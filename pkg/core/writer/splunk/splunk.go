package splunk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/signalfx/golib/v3/trace"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/signalfx/golib/v3/datapoint"
	"github.com/signalfx/golib/v3/event"
	"github.com/signalfx/signalfx-agent/pkg/core/common/httpclient"
	log "github.com/sirupsen/logrus"
)

// Writer posts data to Splunk HTTP Event Collector.
type Writer struct {
	httpClient      *http.Client
	url             string
	hostname        string
	token           string
	traceSource     string
	source          string
	traceSourceType string
	sourceType      string
	index           string
	traceIndex      string
	skipTLSVerify   bool
	events          chan interface{}
	sendQueue       chan []interface{}
}

// Build a Splunk Writer.
func Build(url string, token string, source string, sourceType string, index string, skipTLSVerify bool, hostname string, traceIndex string, traceSource string, traceSourceType string) (Writer, error) {
	handler := Writer{
		url:             url,
		token:           token,
		source:          source,
		sourceType:      sourceType,
		traceIndex:      traceIndex,
		traceSource:     traceSource,
		traceSourceType: traceSourceType,
		index:           index,
		skipTLSVerify:   skipTLSVerify,
		hostname:        hostname,
	}
	err := handler.init()
	return handler, err
}

func (h *Writer) init() error {
	httpConfig := httpclient.HTTPConfig{
		SkipVerify: h.skipTLSVerify,
		UseHTTPS:   strings.HasPrefix(h.url, "https"),
	}

	httpClient, err := httpConfig.Build()
	if err != nil {
		return err
	}
	h.httpClient = httpClient

	if h.hostname == "" {
		h.hostname, _ = os.Hostname()
	}
	h.sendQueue = make(chan []interface{})
	h.events = make(chan interface{})
	go func() {
		eventsToSend := make([]interface{}, 0)
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case e := <-h.events:
				eventsToSend = append(eventsToSend, e)
				if len(eventsToSend) > 50 {
					oldEvents := eventsToSend
					eventsToSend = make([]interface{}, 0)
					h.sendQueue <- oldEvents
				}
				break
			case <-ticker.C:
				if len(eventsToSend) > 0 {
					oldEvents := eventsToSend
					eventsToSend = make([]interface{}, 0)
					h.sendQueue <- oldEvents
				}
			}
		}
	}()
	go func() {
		for {
			events := <-h.sendQueue
			h.logEvents(events)
		}
	}()
	return nil
}

type eventdata struct {
	Category   event.Category    `json:"category"`
	EventType  string            `json:"eventType"`
	Meta       map[string]string `json:"meta"`
	Dimensions map[string]string `json:"dimensions"`
	Properties map[string]string `json:"properties"`
}

type splunkSpan struct {
	Time       int64      `json:"time"`                 // epoch time
	Host       string     `json:"host"`                 // hostname
	Source     string     `json:"source,omitempty"`     // optional description of the source of the event; typically the app's name
	SourceType string     `json:"sourcetype,omitempty"` // optional name of a Splunk parsing configuration; this is usually inferred by Splunk
	Index      string     `json:"index,omitempty"`      // optional name of the Splunk index to store the event in; not required if the token has a default index set in Splunk
	Event      trace.Span `json:"event"`                // a trace span
}

type splunkMetric struct {
	Time       int64             `json:"time"`                 // epoch time
	Host       string            `json:"host"`                 // hostname
	Source     string            `json:"source,omitempty"`     // optional description of the source of the event; typically the app's name
	SourceType string            `json:"sourcetype,omitempty"` // optional name of a Splunk parsing configuration; this is usually inferred by Splunk
	Index      string            `json:"index,omitempty"`      // optional name of the Splunk index to store the event in; not required if the token has a default index set in Splunk
	Event      string            `json:"event"`                // type of event: this is a metric.
	Fields     map[string]string `json:"fields"`               // metric data
}

type splunkEvent struct {
	Time       int64     `json:"time"`                 // epoch time
	Host       string    `json:"host"`                 // hostname
	Source     string    `json:"source,omitempty"`     // optional description of the source of the event; typically the app's name
	SourceType string    `json:"sourcetype,omitempty"` // optional name of a Splunk parsing configuration; this is usually inferred by Splunk
	Index      string    `json:"index,omitempty"`      // optional name of the Splunk index to store the event in; not required if the token has a default index set in Splunk
	Event      eventdata `json:"event"`                // event data
}

func toString(obj interface{}) string {
	if stringer, ok := obj.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%v", obj)
}

func computeTime(timestamp time.Time) int64 {
	if timestamp.IsZero() {
		return time.Now().UnixNano() / time.Millisecond.Nanoseconds()
	}
	return timestamp.UnixNano() / time.Millisecond.Nanoseconds()
}

// LogDataPoint logs a data point as a Splunk metric event
func (h *Writer) LogDataPoint(d *datapoint.Datapoint) {
	fields := make(map[string]string)

	for key, v := range d.Meta {
		if v != nil {
			fields[toString(key)] = toString(v)
		}
	}
	for key, v := range d.Dimensions {
		fields[toString(key)] = toString(v)
	}
	fields["metric_type"] = d.MetricType.String()
	fields[fmt.Sprintf("metric_name:%s", d.Metric)] = d.Value.String()

	newEvent := splunkMetric{
		Time:       computeTime(d.Timestamp),
		Host:       h.hostname,
		Source:     h.source,
		SourceType: h.sourceType,
		Index:      h.index,
		Event:      "metric",
		Fields:     fields,
	}
	h.events <- newEvent
}

// LogEvent logs an event as a Splunk metric event
func (h *Writer) LogEvent(e *event.Event) {
	props := make(map[string]string)
	for key, v := range e.Properties {
		if v != nil {
			props[key] = toString(v)
		}
	}

	meta := make(map[string]string)
	for key, v := range e.Meta {
		if v != nil {
			meta[toString(key)] = toString(v)
		}
	}
	newEvent := splunkEvent{
		Time:       computeTime(e.Timestamp),
		Host:       h.hostname,
		Source:     h.source,
		SourceType: h.sourceType,
		Index:      h.index,
		Event:      eventdata{Properties: props, Dimensions: e.Dimensions, Meta: meta, EventType: e.EventType, Category: e.Category},
	}
	h.events <- newEvent
}

func (h *Writer) LogSpan(span *trace.Span) {

	newEvent := splunkSpan{
		Time:       *span.Timestamp,
		Host:       h.hostname,
		Source:     h.traceSource,
		SourceType: h.traceSourceType,
		Index:      h.traceIndex,
		Event:      *span,
	}
	h.events <- newEvent

}

func (h *Writer) logEvents(events []interface{}) {
	buf := new(bytes.Buffer)
	for _, e := range events {
		b, err := json.Marshal(e)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error marshaling data to Splunk")
		} else {
			buf.Write(b)
			buf.WriteString("\r\n\r\n")
		}
	}

	err := h.doRequest(buf)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error sending data to Splunk")
	}
}

func (h *Writer) doRequest(b io.Reader) error {
	url := h.url
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Splunk "+h.token)

	res, err := h.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		_, _ = io.Copy(ioutil.Discard, res.Body)
		return nil
	default:
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(res.Body)
		responseBody := buf.String()
		err = fmt.Errorf("%s", responseBody)
	}
	return err
}

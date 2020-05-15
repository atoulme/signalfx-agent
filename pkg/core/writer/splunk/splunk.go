package splunk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/signalfx/golib/v3/datapoint"
	"github.com/signalfx/golib/v3/event"
	"github.com/signalfx/signalfx-agent/pkg/core/common/httpclient"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

// Handler posts data to Splunk HTTP Event Collector.
type Handler struct {
	HTTPClient    *http.Client
	URL           string
	Hostname      string
	Token         string
	Source        string
	SourceType    string
	Index         string
	SkipTLSVerify bool
	Errors        chan error
	once          sync.Once
	events        chan splunkMessage
	sendQueue     chan []splunkMessage
}

func (h *Handler) init() {
	if h.HTTPClient == nil {
		httpConfig := httpclient.HTTPConfig{
			SkipVerify: h.SkipTLSVerify,
		}

		httpClient, err := httpConfig.Build()
		if err != nil {

		}
		h.HTTPClient = httpClient
	}
	if h.Hostname == "" {
		h.Hostname, _ = os.Hostname()
	}
	h.sendQueue = make(chan []splunkMessage)
	h.events = make(chan splunkMessage)
	go func() {
		eventsToSend := make([]splunkMessage, 0)
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case e := <-h.events:
				eventsToSend = append(eventsToSend, e)
				if len(eventsToSend) > 50 {
					oldEvents := eventsToSend
					eventsToSend = make([]splunkMessage, 0)
					h.sendQueue <- oldEvents
				}
				break
			case <-ticker.C:
				if len(eventsToSend) > 0 {
					oldEvents := eventsToSend
					eventsToSend = make([]splunkMessage, 0)
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
}

type eventdata struct {
	Category   event.Category    `json:"category"`
	EventType  string            `json:"eventType"`
	Meta       map[string]string `json:"meta"`
	Dimensions map[string]string `json:"dimensions"`
	Properties map[string]string `json:"properties"`
}

type splunkMessage interface {
	Marshall() ([]byte, error)
}

type splunkMetric struct {
	Time       int64             `json:"time"`                 // epoch time in nano-seconds
	Host       string            `json:"host"`                 // hostname
	Source     string            `json:"source,omitempty"`     // optional description of the source of the event; typically the app's name
	SourceType string            `json:"sourcetype,omitempty"` // optional name of a Splunk parsing configuration; this is usually inferred by Splunk
	Index      string            `json:"index,omitempty"`      // optional name of the Splunk index to store the event in; not required if the token has a default index set in Splunk
	Event      string            `json:"event"`                // type of event: this is a metric.
	Fields     map[string]string `json:"fields"`               // metric data
}

type splunkEvent struct {
	Time       int64     `json:"time"`                 // epoch time in nano-seconds
	Host       string    `json:"host"`                 // hostname
	Source     string    `json:"source,omitempty"`     // optional description of the source of the event; typically the app's name
	SourceType string    `json:"sourcetype,omitempty"` // optional name of a Splunk parsing configuration; this is usually inferred by Splunk
	Index      string    `json:"index,omitempty"`      // optional name of the Splunk index to store the event in; not required if the token has a default index set in Splunk
	Event      eventdata `json:"event"`                // event data
}

func (s splunkMetric) Marshall() ([]byte, error) {
	b, err := json.Marshal(s)
	return b, err
}

func (s splunkEvent) Marshall() ([]byte, error) {
	return json.Marshal(s)
}

func toString(obj interface{}) string {
	if stringer, ok := obj.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%v", obj)
}

func (h *Handler) LogDataPoint(d *datapoint.Datapoint) {
	h.once.Do(h.init)
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
		Time:       d.Timestamp.UnixNano(),
		Host:       h.Hostname,
		Source:     h.Source,
		SourceType: h.SourceType,
		Index:      h.Index,
		Event:      "metric",
		Fields:     fields,
	}
	h.events <- newEvent
}

func (h *Handler) LogEvent(e *event.Event) {
	h.once.Do(h.init)
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
		Time:       e.Timestamp.UnixNano(),
		Host:       h.Hostname,
		Source:     h.Source,
		SourceType: h.SourceType,
		Index:      h.Index,
		Event:      eventdata{Properties: props, Dimensions: e.Dimensions, Meta: meta, EventType: e.EventType, Category: e.Category},
	}
	h.events <- newEvent
}

func (h *Handler) logEvents(events []splunkMessage) {
	buf := new(bytes.Buffer)
	for _, e := range events {
		b, err := e.Marshall()
		if err != nil {
			h.Errors <- err
		} else {
			buf.Write(b)
			buf.WriteString("\r\n\r\n")
		}
	}

	err := h.doRequest(buf)
	if err != nil {
		h.Errors <- err
	}
}

func (h *Handler) doRequest(b io.Reader) error {
	url := h.URL
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Splunk "+h.Token)

	res, err := h.HTTPClient.Do(req)
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
		err = fmt.Errorf("%s\n%s", responseBody, b)
	}
	return err
}

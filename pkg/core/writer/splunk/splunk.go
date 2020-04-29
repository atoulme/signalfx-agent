package splunk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/signalfx/golib/v3/datapoint"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
	"github.com/signalfx/golib/v3/event"
)

// Handler posts data to Splunk HTTP Event Collector.
type Handler struct {
	HTTPClient      *http.Client
	URL             string
	Hostname        string
	Token           string
	Source          string
	SourceType      string
	Index           string
	SkipTLSVerify   bool
	Errors          chan error
	once            sync.Once
	events          chan splunkEvent
	sendQueue       chan []splunkEvent
}

func (h *Handler) init() {
	if h.HTTPClient == nil {
		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: h.SkipTLSVerify}}
		h.HTTPClient = &http.Client{Timeout: time.Second * 20, Transport: tr}
	}
	if h.Hostname == "" {
		h.Hostname, _ = os.Hostname()
	}
	h.sendQueue = make(chan []splunkEvent)
	h.events = make(chan splunkEvent)
	go func() {
		eventsToSend := make([]splunkEvent, 0)
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case e := <-h.events:
				eventsToSend = append(eventsToSend, e)
				if len(eventsToSend) > 50 {
					oldEvents := eventsToSend
					eventsToSend = make([]splunkEvent, 0)
					h.sendQueue <- oldEvents
				}
				break
			case <-ticker.C:
				if len(eventsToSend) > 0 {
					oldEvents := eventsToSend
					eventsToSend = make([]splunkEvent, 0)
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
	Category   event.Category `json:"category"`
	EventType  string `json:"eventType"`
	Meta       map[string]string `json:"meta"`
	Dimensions map[string]string `json:"dimensions"`
	Properties map[string]string `json:"properties"`
}

type datapointdata struct {
	EventType  string            `json:"eventType"`
	Meta       map[string]string `json:"meta"`
	Dimensions map[string]string `json:"dimensions"`
	Value      string `json:"value"`
	Metric     string `json:"metric"`
	MetricType datapoint.MetricType `json:"metricType"`
}

type splunkEvent struct {
	Time       int64     `json:"time"`                 // epoch time in nano-seconds
	Host       string    `json:"host"`                 // hostname
	Source     string    `json:"source,omitempty"`     // optional description of the source of the event; typically the app's name
	SourceType string    `json:"sourcetype,omitempty"` // optional name of a Splunk parsing configuration; this is usually inferred by Splunk
	Index      string    `json:"index,omitempty"`      // optional name of the Splunk index to store the event in; not required if the token has a default index set in Splunk
	Event      interface{} `json:"event"`      // event data
}

func toString(obj interface{}) string {
	if stringer, ok := obj.(fmt.Stringer); ok {
		return stringer.String()
	} else {
		return fmt.Sprintf("%v", obj)
	}
}

func (h *Handler) LogDataPoint(d *datapoint.Datapoint) {
	h.once.Do(h.init)

	meta := make(map[string]string)
	for key, v := range d.Meta {
		if v != nil {
			meta[toString(key)] = toString(v)
		}
	}
	event := splunkEvent{
		Time:       d.Timestamp.UnixNano(),
		Host:       h.Hostname,
		Source:     h.Source,
		SourceType: h.SourceType,
		Index:      h.Index,
		Event:  datapointdata{Value: toString(d.Value), Dimensions: d.Dimensions, Meta: meta, Metric: d.Metric, MetricType: d.MetricType},
	}
	h.events <- event
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
	event := splunkEvent{
		Time:       e.Timestamp.UnixNano(),
		Host:       h.Hostname,
		Source:     h.Source,
		SourceType: h.SourceType,
		Index:      h.Index,
		Event:      eventdata{Properties: props, Dimensions: e.Dimensions, Meta: meta, EventType: e.EventType, Category: e.Category},
	}
	h.events <- event
}

func (h *Handler) logEvents(events []splunkEvent) {
	buf := new(bytes.Buffer)
	for _, e := range events {

		b, err := json.Marshal(e)
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

func (h *Handler) doRequest(b *bytes.Buffer) error {
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
		io.Copy(ioutil.Discard, res.Body)
		return nil
	default:
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		responseBody := buf.String()
		err = errors.New(fmt.Sprintf("%s\n%s", responseBody, b))
	}
	return err
}

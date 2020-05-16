package main

import (
	"github.com/opentracing/opentracing-go"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)
import "github.com/signalfx/signalfx-go-tracing/tracing"

func main() {
	tracing.Start(tracing.WithEndpointURL("http://signalfx-agent:9080/v1/trace"), tracing.WithServiceName("tracing"))
	defer tracing.Stop()
	counter := 0
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <- ticker.C:
			span := opentracing.GlobalTracer().StartSpan("main")
			span.SetTag("span.kind", "server")
			span.SetTag("counter", strconv.Itoa(counter))
			time.Sleep(1 * time.Second)
			span.Finish()
			counter++
			break
		case <- c:
			ticker.Stop()
			return
		}
	}

}

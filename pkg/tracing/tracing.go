// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package tracing

import (
	"net/http"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/vogo/clog"
	"sourcegraph.com/sourcegraph/appdash"
	appdashtracer "sourcegraph.com/sourcegraph/appdash/opentracing"
)

var collector appdash.Collector
var tracer opentracing.Tracer

// Start tracing
func Start() {
	collector = appdash.NewRemoteCollector("localhost:3001")
	tracer = appdashtracer.NewTracer(collector)
	opentracing.InitGlobalTracer(tracer)
	clog.Info(nil, "start tracing")
}

//Span new span from request
func Span(r *http.Request) opentracing.Span {
	// Start a new root Span and therefore a new trace.
	span := opentracing.StartSpan(r.URL.Path)

	// OpenTracing allows for arbitrary tags to be added to a Span.
	span.SetTag("Request.Host", r.Host)
	span.SetTag("Request.Address", r.RemoteAddr)
	addHeaderTags(span, r.Header)
	err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, r.Header)
	// We have no better place to record an error than the Span itself :-/
	if err != nil {
		clog.Error(r.Context(), "failed to inject span: %v", err)
	}
	return span
}

const headerTagPrefix = "Request.Header."

//SpanHTTPHeader whether span header
func SpanHTTPHeader(h string) bool {
	return strings.HasPrefix(h, "Ot-Tracer-")
}

// addHeaderTags adds header key:value pairs to a span as a tag with the prefix
// "Request.Header.*"
func addHeaderTags(span opentracing.Span, h http.Header) {
	for k, v := range h {
		span.SetTag(headerTagPrefix+k, strings.Join(v, ", "))
	}
}

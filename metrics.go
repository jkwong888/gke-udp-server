package main

import (
	"log"
	"net/http"

	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
)

func enableObservabilityAndExporters(mux *http.ServeMux) {
	// Stats exporter: Prometheus
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "udp-serve",
	})
	if err != nil {
		log.Fatalf("Failed to create the Prometheus stats exporter: %v", err)
	}

	view.RegisterExporter(pe)

	mux.Handle("/metrics", pe)

	/*
		// Trace exporter: Zipkin
		localEndpoint, err := openzipkin.NewEndpoint("ochttp_tutorial", "localhost:5454")
		if err != nil {
			log.Fatalf("Failed to create the local zipkinEndpoint: %v", err)
		}
		reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
		ze := zipkin.NewExporter(reporter, localEndpoint)
		trace.RegisterExporter(ze)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	*/
}
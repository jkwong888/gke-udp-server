package main

import (
	"net/http"

	"github.com/heptiolabs/healthcheck"
)


func NullHealthCheck() healthcheck.Check {
	return func() error {
		return nil
	}
}

func enableHealthCheck(mux *http.ServeMux) {
	// add health check
	health := healthcheck.NewHandler()

	health.AddReadinessCheck(
		"null",
		NullHealthCheck())

	health.AddLivenessCheck(
		"null",
		NullHealthCheck())

	mux.HandleFunc("/health/live", health.LiveEndpoint)
	mux.HandleFunc("/health/ready", health.ReadyEndpoint)

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
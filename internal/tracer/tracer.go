package tracer

import (
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/proto/zipkin_proto3"
	"github.com/openzipkin/zipkin-go/reporter"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	exampleconfig "github.com/smlobo/zipkin-go-example/internal/config"
)

func NewTracer(serviceName string, port uint16) (*zipkin.Tracer, error) {

	var reporter reporter.Reporter

	if exampleconfig.Config["ENCODING"] == "json" {
		// Default JSON V2 reporter
		reporter = httpreporter.NewReporter(exampleconfig.Config["ENDPOINT"])
	} else if exampleconfig.Config["ENCODING"] == "protobuf" {
		// Protobuf reporter
		reporterOption := httpreporter.Serializer(zipkin_proto3.SpanSerializer{})
		reporter = httpreporter.NewReporter(exampleconfig.Config["ENDPOINT"], reporterOption)
	}

	// Local endpoint represent the local service information
	localEndpoint := &model.Endpoint{ServiceName: serviceName, Port: port}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00)
	// of traces.
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, err
	}

	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		return nil, err
	}

	return tracer, err
}

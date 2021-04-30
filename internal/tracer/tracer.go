// Copyright 2021 Sheldon Lobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tracer

import (
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/proto/zipkin_proto3"
	"github.com/openzipkin/zipkin-go/reporter"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	exampleconfig "github.com/smlobo/zipkin-go-example/internal/config"
	"log"
	"net"
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
	localEndpoint := &model.Endpoint{ServiceName: serviceName, IPv4: getOutboundIP(), Port: port}

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

// Get preferred outbound ip of this machine
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}


package main

import (
	"fmt"
	"github.com/gorilla/mux"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	exampleconfig "github.com/smlobo/zipkin-go-example/internal/config"
	examplehandler "github.com/smlobo/zipkin-go-example/internal/handler"
	exampletracer "github.com/smlobo/zipkin-go-example/internal/tracer"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Global uint16 ports from config
var frontendPort uint16
var backendPort uint16

func setupFrontend() {

	tracer, err := exampletracer.NewTracer("go-frontend", frontendPort)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	// create global zipkin http server middleware
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		tracer, zipkinhttp.TagResponseSize(true),
	)

	// create global zipkin traced http client
	client, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}

	router.Methods("GET").Path("/").HandlerFunc(examplehandler.FrontendHandler(client))

	router.Use(serverMiddleware)

	log.Println("Starting frontend at :", frontendPort)
	portString := fmt.Sprintf(":%d", frontendPort)
	log.Fatal(http.ListenAndServe(portString, router))
}


func setupBackend() {

	tracer, err := exampletracer.NewTracer("go-backend", backendPort)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	// create global zipkin http server middleware
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		tracer, zipkinhttp.TagResponseSize(true),
	)

	router.Methods("POST").Path("/").HandlerFunc(examplehandler.BackendHandler())

	router.Use(serverMiddleware)

	log.Println("Starting backend at :", backendPort)
	portString := fmt.Sprintf(":%d", backendPort)
	log.Fatal(http.ListenAndServe(portString, router))
}


func main() {

	// Usage
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./bin/zipkin-go-example <frontend|backend>")
		os.Exit(1)
	}

	// Do config
	exampleconfig.SetupConfig()

	// Ports
	feP, _ := strconv.Atoi(exampleconfig.Config["FRONTEND_PORT"])
	frontendPort = uint16(feP)
	beP, _ := strconv.Atoi(exampleconfig.Config["BACKEND_PORT"])
	backendPort = uint16(beP)

	if os.Args[1] == "frontend" {
		setupFrontend()
	} else if os.Args[1] == "backend" {
		setupBackend()
	}
}

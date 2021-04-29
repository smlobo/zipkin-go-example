# zipkin-go-example
An example app where two Golang services collaborate on an http request. The timing of these requests are recorded into ZIpkin.

# Setup
* `make module`
* `make`

# Zipkin Receiver
* `docker run -d -p 9411:9411 openzipkin/zipkin`

# Run
Execute the `frontend` and `backend` services on separate terminals.
For example:
* `./bin/zipkin-go-example backend`
* `ENCODING="protobuf" ./bin/zipkin-go-example frontend`

# Generate Traces
* `curl http://localhost:8081`

![Example Trace](screenshot.png?raw=true "Example Trace")

# References
* https://medium.com/devthoughts/instrumenting-a-go-application-with-zipkin-b79cc858ac3e
* https://pkg.go.dev/github.com/openzipkin/zipkin-go

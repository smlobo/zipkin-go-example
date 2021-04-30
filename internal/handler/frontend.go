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

package handler

import (
	"fmt"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	exampleconfig "github.com/smlobo/zipkin-go-example/internal/config"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func FrontendHandler(client *zipkinhttp.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received frontend request to:", r.Host, r.URL.Path, "::", r.Method)

		// retrieve span from context (created by server middleware)
		span := zipkin.SpanFromContext(r.Context())
		span.Tag("frontend_key", "frontend value")

		// doing some expensive calculations....
		time.Sleep(25 * time.Millisecond)
		span.Annotate(time.Now(), "frontend expensive_calc_done")

		// Make wrapped call to backend
		backendURL := fmt.Sprintf("http://localhost:%s/", exampleconfig.Config["BACKEND_PORT"])
		newRequest, err := http.NewRequest("POST", backendURL, nil)
		if err != nil {
			log.Printf("unable to create client: %+v\n", err)
			http.Error(w, err.Error(), 500)
			return
		}

		ctx := zipkin.NewContext(newRequest.Context(), span)

		newRequest = newRequest.WithContext(ctx)

		//backendResponse, err := client.DoWithAppSpan(newRequest, "backend-call")
		backendResponse, err := client.Do(newRequest)
		if err != nil {
			log.Fatal("Bad backend request:", backendURL, ";", err)
		}

		backendBody := "Bad backend"
		if backendResponse.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(backendResponse.Body)
			if err != nil {
				log.Fatal(err)
			}
			backendBody = string(bodyBytes)
		}
		backendResponse.Body.Close()

		time.Sleep(25 * time.Millisecond)

		w.WriteHeader(http.StatusOK)
		responseBody := fmt.Sprintf("From frontend: %s [%s]\n",
			time.Now().Local().Format("15:04:05.000"), backendBody)
		w.Write([]byte(responseBody))
	}
}

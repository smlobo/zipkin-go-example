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
	"log"
	"net/http"
	"time"
)

func BackendHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received backend request to:", r.Host, r.URL.Path, "::", r.Method)

		// retrieve span from context (created by server middleware)
		span := zipkin.SpanFromContext(r.Context())
		span.Tag("backend_key", "backend value")

		// doing some expensive calculations....
		time.Sleep(25 * time.Millisecond)
		span.Annotate(time.Now(), "backend expensive_calc_done")
		time.Sleep(25 * time.Millisecond)

		w.WriteHeader(http.StatusOK)
		responseBody := fmt.Sprintf("From backend: %s", time.Now().Local().Format("15:04:05.000"))
		w.Write([]byte(responseBody))
	}
}

/*
Copyright © 2023 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package tests contains REST API tests for following endpoints:
//
// apiPrefix + "content"
package tests

import (
	"encoding/json"
	"net/http"

	"github.com/verdverm/frisby"
)

const contentEndpointPostfix = "content"

// ContentAndGroups represents response from /content endpoint
type ContentAndGroups struct {
	Content []string `json:"content"`
	Groups  []Group  `json:"groups"`
	Status  string   `json:"status"`
}

// checkContentEndpoint check if the 'content' point (usually
// /api/insights-results-aggregator/v2/content) responds correctly to HTTP GET
// command
func checkContentEndpoint() {
	f := frisby.Create("Check the 'content' REST API point using HTTP GET method").Get(apiURL + contentEndpointPostfix)
	f.Send()
	f.ExpectStatus(http.StatusOK)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := ContentAndGroups{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.Status != "ok" {
			f.AddError(statusShouldBeSetToOK)
		}
		if len(response.Groups) == 0 {
			f.AddError("Groups node is empty")
		}
	}
	f.PrintReport()
}

// check whether other HTTP methods are rejected correctly for the REST API
// 'content' point
func checkWrongMethodsForContentEndpoint() {
	checkGetEndpointByOtherMethods(apiURL+contentEndpointPostfix, false)
}

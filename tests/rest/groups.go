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
// apiPrefix
// apiPrefix + "groups"
package tests

import (
	"encoding/json"

	"github.com/verdverm/frisby"
)

const endpointPostfix = "groups"

// Group structure represents one group entry in groups array
type Group struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// GroupsResponse represents response from /organizations endpoint
type GroupsResponse struct {
	Groups []Group `json:"groups"`
	Status string  `json:"status"`
}

// checkGroupsEndpoint check if the 'groups' point (usually /api/insights-results-aggregator/v2/groups) responds correctly to HTTP GET command
func checkGroupsEndpoint() {
	f := frisby.Create("Check the 'groups' REST API point using HTTP GET method").Get(apiURL + endpointPostfix)
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := GroupsResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.Status != "ok" {
			f.AddError("Expecting 'status' to be set to 'ok'")
		}
		if len(response.Groups) == 0 {
			f.AddError("Groups node is empty")
		}
	}
	f.PrintReport()
}

// check whether other HTTP methods are rejected correctly for the REST API 'groups' point
func checkWrongMethodsForGroupsEndpoint() {
	checkGetEndpointByOtherMethods(apiURL+endpointPostfix, false)
}

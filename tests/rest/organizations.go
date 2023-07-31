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
// apiPrefix + "clusters"
// apiPrefix + "groups"
// apiPrefix + "organizations"
package tests

import (
	"encoding/json"
	"fmt"

	"github.com/verdverm/frisby"
)

const organizationsEndpointPostfix = "organizations"

// OrganizationsResponse represents response from /organizations endpoint
type OrganizationsResponse struct {
	Organizations []int  `json:"organizations"`
	Status        string `json:"status"`
}

func testOrganizationIDExistence(f *frisby.Frisby, organizations []int, organizationID int) {
	for _, orgID := range organizations {
		if orgID == organizationID {
			// found it
			return
		}
	}
	// not found
	errorMessage := fmt.Sprintf("Organization %d has not been found", organizationID)
	f.AddError(errorMessage)
}

// checkOrganizationsEndpoint check if the 'organizations' point (usually /api/insights-results-aggregator/v2/organizations) responds correctly to HTTP GET command
func checkOrganizationsEndpoint() {
	f := frisby.Create("Check the 'organizations' REST API point using HTTP GET method").Get(apiURL + organizationsEndpointPostfix)
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := OrganizationsResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.Status != "ok" {
			f.AddError(statusShouldBeSetToOK)
		}
		if len(response.Organizations) == 0 {
			f.AddError("Organizations node is empty")
		}
		testOrganizationIDExistence(f, response.Organizations, organization1)
		testOrganizationIDExistence(f, response.Organizations, organization2)
	}
	f.PrintReport()
}

// check whether other HTTP methods are rejected correctly for the REST API 'organizations' point
func checkWrongMethodsForOrganizationsEndpoint() {
	checkGetEndpointByOtherMethods(apiURL+organizationsEndpointPostfix, false)
}

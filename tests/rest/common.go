/*
Copyright © 2020 Red Hat, Inc.

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

package tests

import (
	"github.com/verdverm/frisby"
)

// common constants used by REST API tests
const (
	apiURL            = "http://localhost:8080/api/insights-results-aggregator/v2/"
	contentTypeHeader = "Content-Type"

	// ContentTypeJSON represents MIME type for JSON format
	ContentTypeJSON = "application/json; charset=utf-8"

	// ContentTypeText represents MIME type for plain text format
	ContentTypeText = "text/plain; charset=utf-8"
)

// messages
const (
	statusShouldBeSetToOK = "Expecting 'status' to be set to 'ok'"
	noPermissionsForOrg   = "You have no permissions to get or change info about this organization"
)

// known IDs
const (
	organization1 = 11789772
	organization2 = 11940171

	cluster1ForOrg1 = "34c3ecc5-624a-49a5-bab8-4fdc5e51a266"
	cluster2ForOrg1 = "ee7d2bf4-8933-4a3a-8634-3328fe806e08"
	cluster3ForOrg1 = "eeeeeeee-eeee-eeee-eeee-000000000001"
	unknownCluster  = "ffffeeee-eeee-eeee-eeee-000000000001"
)

// StatusOnlyResponse represents response containing just a status
type StatusOnlyResponse struct {
	Status string `json:"status"`
}

// sendAndExpectStatus sends the request to the server and checks whether expected HTTP code (status) is returned
func sendAndExpectStatus(f *frisby.Frisby, expectedStatus int) {
	f.Send()
	f.ExpectStatus(expectedStatus)
	f.PrintReport()
}

// checkGetEndpointByOtherMethods checks whether a 'GET' endpoint respond correctly if other HTTP methods are used
func checkGetEndpointByOtherMethods(endpoint string, includingOptions bool) {
	f := frisby.Create("Check the end point " + endpoint + " with wrong method: POST").Post(endpoint)
	sendAndExpectStatus(f, 405)

	f = frisby.Create("Check the entry point " + endpoint + " with wrong method: PUT").Put(endpoint)
	sendAndExpectStatus(f, 405)

	f = frisby.Create("Check the entry point " + endpoint + " with wrong method: DELETE").Delete(endpoint)
	sendAndExpectStatus(f, 405)

	f = frisby.Create("Check the entry point " + endpoint + " with wrong method: PATCH").Patch(endpoint)
	sendAndExpectStatus(f, 405)

	f = frisby.Create("Check the entry point " + endpoint + " with wrong method: HEAD").Head(endpoint)
	sendAndExpectStatus(f, 405)

	// some endpoints accepts OPTIONS method together with GET one, so this check is fully optional
	if includingOptions {
		f = frisby.Create("Check the entry point " + endpoint + " with wrong method: OPTIONS").Options(endpoint)
		sendAndExpectStatus(f, 405)
	}
}

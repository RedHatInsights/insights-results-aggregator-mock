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

package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/verdverm/frisby"
)

// RequestStruct represents one entry in list of requests
type RequestStruct struct {
	RequestID string    `json:"requestID"`
	Valid     bool      `json:"valid"`
	Received  time.Time `json:"received"`
	Processed time.Time `json:"processed"`
}

// RequestResponse represents response containing list of requests
type RequestResponse struct {
	Cluster  string          `json:"cluster"`
	Requests []RequestStruct `json:"requests"`
	Status   string          `json:"status"`
}

// allRequestsIDsEndpointForCluster helper function constructs URL for
// accessing endpoint to retrieve list of request IDs for given cluster
func allRequestsIDsEndpointForCluster(clusterName string) string {
	return fmt.Sprintf("%scluster/%s/requests/", apiURL, clusterName)
}

// checkListAllRequestIDsForKnownCluster checks if expected structure with
// request IDs is returned for known cluster
func checkListAllRequestIDsForKnownCluster() {
	// clusterName represents known cluster with 12 request IDs
	const clusterName = "34c3ecc5-624a-49a5-bab8-4fdc5e51a266"

	url := allRequestsIDsEndpointForCluster(clusterName)
	f := frisby.Create("Check the 'requests' REST API point using HTTP GET method with known cluster").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusOK)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := RequestResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.Status != "ok" {
			f.AddError(statusShouldBeSetToOK)
		}
		if response.Cluster != clusterName {
			f.AddError("Improper cluster name returned")
		}
		if len(response.Requests) != 12 {
			f.AddError("Improper number of request IDs returned")
		}
		// just check the first one
		if response.Requests[0].RequestID != "3nl2vda87ld6e3s25jlk7n2dna" {
			f.AddError("Improper request ID detected")
		}
	}
	f.PrintReport()
}

// checkListAllRequestIDsEmptyList checks if empty list of request IDs is
// returned for known cluster
func checkListAllRequestIDsEmptyList() {
	// clusterName represents known cluster with no request IDs
	const clusterName = "eeeeeeee-eeee-eeee-eeee-000000000001"

	url := allRequestsIDsEndpointForCluster(clusterName)
	f := frisby.Create("Check the 'requests' REST API point using HTTP GET method with known cluster").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusOK)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := RequestResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.Status != "ok" {
			f.AddError(statusShouldBeSetToOK)
		}
		if response.Cluster != clusterName {
			f.AddError("Improper cluster name returned")
		}
		// no requests IDs should be returned
		if len(response.Requests) != 0 {
			f.AddError("Improper number of request IDs returned")
		}
	}
	f.PrintReport()
}

// checkListAllRequestIDsForUnknownCluster checks how unknown cluster is
// handled by the mock service
func checkListAllRequestIDsForUnknownCluster() {
	// clusterName represents unknown cluster
	const clusterName = "ffffffff-ffff-ffff-ffff-000000000001"

	url := allRequestsIDsEndpointForCluster(clusterName)
	f := frisby.Create("Check the 'requests' REST API point using HTTP GET method with known cluster").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusNotFound)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)
	f.PrintReport()
}

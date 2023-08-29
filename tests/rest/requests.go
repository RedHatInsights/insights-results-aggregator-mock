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
	"fmt"
	"time"
)

// RequestStruct represents one entry in list of requests
type RequestStruct struct {
	RequestID string    `json:"requestID"`
	Valid     bool      `json:"valid"`
	Received  time.Time `json:"received"`
	Processed time.Time `json:processed"`
}

// RequestsResponse represents response containing list of requests
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

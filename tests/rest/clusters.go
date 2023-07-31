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

// ClustersResponse represents response containing list of clusters for given organization
type ClustersResponse struct {
	Clusters []string `json:"clusters"`
	Status   string   `json:"status"`
}

func clustersEndpointForOrganization(orgID int) string {
	return fmt.Sprintf("%sorganizations/%d/clusters", apiURL, orgID)
}

func testClusterIDExistence(f *frisby.Frisby, clusters []string, searchedClusterID string) {
	for _, clusterID := range clusters {
		if clusterID == searchedClusterID {
			// found it
			return
		}
	}
	// not found
	errorMessage := fmt.Sprintf("Cluster %s has not been found", searchedClusterID)
	f.AddError(errorMessage)
}

// checkClustersEndpoint check if the 'clusters' point (usually
// /api/insights-results-aggregator/v2/organizations/{org_id}/clusters)
// responds correctly to HTTP GET command
func checkClustersEndpointForOrganization1() {
	f := frisby.Create("Check the 'clusters' REST API point using HTTP GET method").Get(clustersEndpointForOrganization(organization1))
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := ClustersResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.Status != "ok" {
			f.AddError(statusShouldBeSetToOK)
		}
		if len(response.Clusters) == 0 {
			f.AddError("Clusters node is empty")
		}
		testClusterIDExistence(f, response.Clusters, cluster1ForOrg1)
		testClusterIDExistence(f, response.Clusters, cluster2ForOrg1)
		testClusterIDExistence(f, response.Clusters, cluster3ForOrg1)
	}
	f.PrintReport()
}

// checkClustersEndpoint check if the 'clusters' point (usually
// /api/insights-results-aggregator/v2/organizations/{org_id}/clusters)
// responds for organization w/o the right permission to retrieve info
func checkClustersEndpointForOrganization2() {
	f := frisby.Create("Check the 'clusters' REST API point using HTTP GET method").Get(clustersEndpointForOrganization(organization2))
	f.Send()
	f.ExpectStatus(403)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := StatusOnlyResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.Status != "You have no permissions to get or change info about this organization" {
			f.AddError("Wrong status message received: " + response.Status)
		}
	}
	f.PrintReport()
}

// check whether other HTTP methods are rejected correctly for the REST API 'clusters' point
func checkWrongMethodsForClustersEndpoint() {
	// known organizations
	checkGetEndpointByOtherMethods(clustersEndpointForOrganization(organization1), false)
	checkGetEndpointByOtherMethods(clustersEndpointForOrganization(organization2), false)

	// unknown organizations
	checkGetEndpointByOtherMethods(clustersEndpointForOrganization(1), false)
	checkGetEndpointByOtherMethods(clustersEndpointForOrganization(2), false)
}

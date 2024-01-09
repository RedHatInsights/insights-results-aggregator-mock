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

// ClustersDetails represents response containing list of clusters for which
// selected rule was hit
type ClustersDetails struct {
	MetaData ClustersDetailsMetadata `json:"meta"`
	Clusters []string                `json:"data"`
}

// ClustersDetailsMetadata represents metadata about number of cluster hits
// etc.
type ClustersDetailsMetadata struct {
	Count       int       `json:"count"`
	Component   string    `json:"component"`
	ErrorKey    string    `json:"error_key"`
	GeneratedAt time.Time `json:"generated_at"`
}

// clustersDetailsEndpointForRule constructs an URL for clusters detail
// endpoint
func clustersDetailsEndpointForRule(component, errorKey string) string {
	return fmt.Sprintf("%srule/%s|%s/clusters_detail", apiURL, component, errorKey)
}

// checkRetrieveClusterDetailsForKnownRule checks if the
// 'rule/{rule}/clusters_detail' point responds correctly to HTTP GET command
// (for known rule)
func checkRetrieveClusterDetailsForKnownRule() {
	const component = "ccx_rules_ocp.external.rules.nodes_requirements_check.report"
	const errorKey = "NODES_MINIMUM_REQUIREMENTS_NOT_MET"

	url := clustersDetailsEndpointForRule(component, errorKey)
	f := frisby.Create("Check the 'rule/{rule}/clusters_detail' REST API point using HTTP GET method (known rule)").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusOK)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := ClustersDetails{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.MetaData.Count != 24 {
			f.AddError("Improper metadata about number of clusters returned")
		}
		if response.MetaData.Component != component {
			f.AddError("Invalid component")
		}
		if response.MetaData.ErrorKey != errorKey {
			f.AddError("Invalid error key")
		}
		if len(response.Clusters) != 24 {
			f.AddError("Improper number of clusters returned")
		}
		// try just few clusters
		testClusterIDExistence(f, response.Clusters, "00000001-624a-49a5-bab8-4fdc5e51a266")
		testClusterIDExistence(f, response.Clusters, "00000001-6577-4e80-85e7-697cb646ff37")
		testClusterIDExistence(f, response.Clusters, "00000001-8933-4a3a-8634-3328fe806e08")
	}
	f.PrintReport()
}

// checkRetrieveClusterDetailsForUnknownRule checks if the
// 'rule/clusters_detail' point responds correctly to HTTP GET command for
// unknown rule
func checkRetrieveClusterDetailsForUnknownRule() {
	const component = "this.is.unknown.component"
	const errorKey = "THIS_IS_NOT_KNOW_ERROR_KEY"

	url := clustersDetailsEndpointForRule(component, errorKey)
	f := frisby.Create("Check the 'rule/{rule}/clusters_detail' REST API point using HTTP GET method (unknown rule)").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusOK)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := ClustersDetails{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.MetaData.Count != 0 {
			f.AddError("Improper metadata about number of clusters returned")
		}
		if response.MetaData.Component != component {
			f.AddError("Invalid component")
		}
		if response.MetaData.ErrorKey != errorKey {
			f.AddError("Invalid error key")
		}
		if len(response.Clusters) != 0 {
			f.AddError("Improper number of clusters returned")
		}
	}
	f.PrintReport()
}

// checkRetrieveClusterDetailsForImproperRule checks if the
// 'rule/clusters_detail' point responds correctly to HTTP GET command for
// improper rule
func checkRetrieveClusterDetailsForImproperRule() {
	const component = ""
	const errorKey = ""

	url := clustersDetailsEndpointForRule(component, errorKey)
	f := frisby.Create("Check the 'rule/clusters_detail' REST API point using HTTP GET method (improper rule)").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusBadRequest)

	f.PrintReport()
}

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
	"net/http"

	"github.com/verdverm/frisby"

	"github.com/RedHatInsights/insights-results-aggregator-mock/types"
)

// FullReportResponse represents response containing report for given cluster
type FullReportResponse struct {
	Report types.ReportResponse `json:"report"`
	Status string               `json:"status"`
}

// reportEndpoint helper function constructs URL for accessing endpoint to
// retrieve report for given organization and cluster
func reportEndpoint(orgID int, clusterName string) string {
	return fmt.Sprintf("%sreport/%d/%s", apiURL, orgID, clusterName)
}

// checkReportForKnownOrganization checks if proper report is returned for
// known organization ID and known cluster name
func checkReportForKnownOrganizationKnownCluster() {
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method").Get(reportEndpoint(organization1, cluster1ForOrg1))
	f.Send()
	f.ExpectStatus(http.StatusOK)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := FullReportResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.Status != "ok" {
			f.AddError(statusShouldBeSetToOK)
		}
	}
	f.PrintReport()
}

// checkReportForUknownOrganization checks how uknown organization ID is
// checked by REST API handler
func checkReportForUnknownOrganization() {
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method").Get(reportEndpoint(1234, unknownCluster))
	f.Send()
	f.ExpectStatus(http.StatusNotFound)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	f.PrintReport()
}

// checkReportForImproperOrganization checks how improperly entered
// organization ID is checked by REST API handler
func checkReportForImproperOrganization() {
	url := fmt.Sprintf("%sreport/foobar/%s", apiURL, cluster1ForOrg1)
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusBadRequest)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	f.PrintReport()
}

// checkReportForKnownOrganizationUnknownCluster checks how unknown cluster
// name is checked by REST API handler
func checkReportForKnownOrganizationUnknownCluster() {
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method").Get(reportEndpoint(organization1, unknownCluster))
	f.Send()
	f.ExpectStatus(http.StatusNotFound)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	f.PrintReport()
}

// checkReportForKnownOrganizationWrongCluster checks how improper cluster name
// is checked by REST API handler
func checkReportForKnownOrganizationWrongCluster() {
	clusterName := "abcdefghijklmnopqrstuvwyz"
	f := frisby.Create("Check the 'report' REST API point using HTTP GET method").Get(reportEndpoint(organization1, clusterName))
	f.Send()
	f.ExpectStatus(http.StatusBadRequest)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	f.PrintReport()
}

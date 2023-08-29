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

	"github.com/verdverm/frisby"
)

// ListOfDVONamespaces structure represents response for namespaces/dvo REST
// API endpoint
type ListOfDVONamespaces struct {
	Status    string            `json:"status"`
	Workloads []DVOWorkloadItem `json:"workloads"`
}

// DVOWorkloadItem structure represents one entry in list of workloads
type DVOWorkloadItem struct {
	Cluster   ClusterEntry   `json:"cluster"`
	Namespace NamespaceEntry `json:"namespace"`
	Reports   []Report       `json:"reports"`
}

// ClusterEntry structure represents cluster info in namespaces/dvo payload
type ClusterEntry struct {
	UUID        string `json":"uuid"`
	DisplayName string `json:"display_name"`
}

// NamespaceEntry structure represents namespace info in namespaces/dvo payload
type NamespaceEntry struct {
	UUID string `json":"uuid"`
	Name string `json:"name"`
}

// Report structure represents one report in namespaces/dvo list of reports
type Report struct {
	Check       string `json:"check"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	Remediation string `json:"remediation"`
}

// dvoNamespacesEndpoint constructs an URL for list of all DVO namespaces
func dvoNamespacesEndpoint() string {
	return fmt.Sprintf("%snamespaces/dvo", apiURL)
}

func checkListOfDVONamespaces() {
	url := dvoNamespacesEndpoint()
	f := frisby.Create("Check the 'namespaces/dvo' REST API point using HTTP GET method").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusOK)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := ListOfDVONamespaces{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}

		// check elementary metadata
		if response.Status != "ok" {
			f.AddError("Status is not set to ok")
		}
		if len(response.Workloads) != 1 {
			f.AddError("Just one workload is expected")
		}

		// check cluster entry
		uuid := response.Workloads[0].Cluster.UUID
		if uuid != "00000001-0001-0001-0001-000000000001" {
			f.AddError("Improper cluster UUID: " + uuid)
		}
		displayName := response.Workloads[0].Cluster.DisplayName
		if displayName != "Cluster #1" {
			f.AddError("Improper cluster display name: " + displayName)
		}

		// check namespace entry
		uuid = response.Workloads[0].Namespace.UUID
		if uuid != "00000001-0001-0001-0001-000000000001" {
			f.AddError("Improper namespace UUID: " + uuid)
		}
		name := response.Workloads[0].Namespace.Name
		if name != "Namespace #2" {
			f.AddError("Improper namespace name: " + name)
		}

		// check reports
		if len(response.Workloads[0].Reports) != 2 {
			f.AddError("Two reports are expected")
		}

		// test first report
		report := response.Workloads[0].Reports[0]
		expected := Report{
			Check:       "no_anti_affinity",
			Kind:        "Deployment",
			Description: "Indicates when... ... ...",
			Remediation: "Specify anti-affinity in your pod specification ... ... ...",
		}
		if report != expected {
			f.AddError("First report is different from expected one")
		}

		// test second report
		report = response.Workloads[0].Reports[1]
		expected = Report{
			Check:       "run_as_non_root",
			Kind:        "Runtime",
			Description: "Indicates when... ... ...",
			Remediation: "Select different user to run this deployment... ... ...",
		}
		if report != expected {
			f.AddError("First report is different from expected one")
		}

	}
	f.PrintReport()
}

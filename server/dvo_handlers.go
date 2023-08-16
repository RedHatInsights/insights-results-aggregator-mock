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

package server

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Constants used to produce HTTP responses
const (
	cluster1UUID        = "00000001-0001-0001-0001-000000000001"
	cluster1DisplayName = "Cluster #1"

	namespace1UUID     = "00000001-0001-0001-0001-000000000001"
	namespace1FullName = "Namespace #1"
	namespace2UUID     = "00000001-0001-0001-0001-000000000001"
	namespace2FullName = "Namespace #2"

	dvoReport1Check       = "no_anti_affinity"
	dvoReport1Kind        = "Deployment"
	dvoReport1Description = "Indicates when... ... ..."
	dvoReport1Remediation = "Specify anti-affinity in your pod specification ... ... ..."

	dvoReport2Check       = "run_as_non_root"
	dvoReport2Kind        = "Runtime"
	dvoReport2Description = "Indicates when... ... ..."
	dvoReport2Remediation = "Select different user to run this deployment... ... ..."
)

// AllDVONamespacesResponse is a data structure that represents list of namespace
type AllDVONamespacesResponse struct {
	Status    string     `json:"status"`
	Workloads []Workload `json:"workloads"`
}

// Workload structure represents one workload entry in list of workloads
type Workload struct {
	ClusterEntry ClusterEntry   `json:"cluster"`
	Namespace    NamespaceEntry `json:"namespace"`
	Reports      []DVOReport    `json:"reports"`
}

// ClusterEntry structure contains cluster UUID and cluster name
type ClusterEntry struct {
	UUID        string `json:"uuid"`
	DisplayName string `json:"display_name"`
}

// NamespaceEntry structure contains basic information about namespace
type NamespaceEntry struct {
	UUID     string `json:"uuid"`
	FullName string `json:"name"`
}

// DVOReport structure represents one DVO-related report
type DVOReport struct {
	Check       string `json:"check"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	Remediation string `json:"remediation"`
}

// allDVONamespaces handler returns list of all DVO namespaces. Currently it
// does not depend on Organization ID as this information is passed through
// Bearer token in real Smart Proxy service. The format of output should be:
//
//	  {
//	    "status": "ok",
//	    "workloads": [
//	        {
//	            "cluster": {
//	                "uuid": "{cluster UUID}",
//	                "display_name": "{cluster UUID or displayable name}",
//	            },
//	            "namespace": {
//	                "uuid": "{namespace UUID}",
//	                "name": "{namespace real name}", // optional, might be null
//	            },
//	            "reports": [
//	                {
//	                    "check": "{for example no_anti_affinity}", // taken from the original full name deploment_validation_operator_no_anti_affinity
//	                    "kind": "{kind attribute}",
//	                    "description": {description}",
//	                    "remediation": {remediation}",
//	                },
//	            ]
//	    ]
//	}
func (server *HTTPServer) allDVONamespaces(writer http.ResponseWriter, request *http.Request) {
	log.Info().Msg("All DVO namespaces handler")

	// prepare response structure
	var responseData AllDVONamespacesResponse
	responseData.Status = "ok"
	responseData.Workloads = []Workload{
		Workload{
			ClusterEntry{
				UUID:        cluster1UUID,
				DisplayName: cluster1DisplayName,
			},
			NamespaceEntry{
				UUID:     namespace2UUID,
				FullName: namespace2FullName,
			},
			[]DVOReport{
				DVOReport{
					Check:       dvoReport1Check,
					Kind:        dvoReport1Kind,
					Description: dvoReport1Description,
					Remediation: dvoReport1Remediation,
				},
				DVOReport{
					Check:       dvoReport2Check,
					Kind:        dvoReport2Kind,
					Description: dvoReport2Description,
					Remediation: dvoReport2Remediation,
				},
			},
		},
	}

	// transform response structure into proper JSON payload
	bytes, err := json.MarshalIndent(responseData, "", "\t")
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
		return
	}

	// and send the response to client
	_, err = writer.Write(bytes)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

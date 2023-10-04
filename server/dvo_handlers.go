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
	"time"

	"github.com/rs/zerolog/log"
)

// AllDVONamespacesResponse is a data structure that represents list of namespace
type AllDVONamespacesResponse struct {
	Status    string     `json:"status"`
	Workloads []Workload `json:"workloads"`
}

// Workload structure represents one workload entry in list of workloads
type Workload struct {
	ClusterEntry  ClusterEntry   `json:"cluster"`
	Namespace     NamespaceEntry `json:"namespace"`
	MetadataEntry MetadataEntry  `json:"metadata"`
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

// MetadataEntry structure contains basic information about workload metadata
type MetadataEntry struct {
	Recommendations int    `json:"recommendations"`
	Objects         int    `json:"objects"`
	ReportedAt      string `json:"reported_at"`
	LastCheckedAt   string `json:"last_checked_at"`
	HighestSeverity int    `json:"highest_severity"`
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
//		  {
//		    "status": "ok",
//		    "workloads": [
//		        {
//		            "cluster": {
//		                "uuid": "{cluster UUID}",
//		                "display_name": "{cluster UUID or displayable name}",
//		            },
//		            "namespace": {
//		                "uuid": "{namespace UUID}",
//		                "name": "{namespace real name}", // optional, might be null
//		            },
//	                 metadata": {
//	                     "recommendations": "{number of recommendations"}, // stored in DVO_REPORT table, computed as SELECT count(distinct(recommendation)) WHERE cluster="{cluster UUID}" and namespace="{namespace UUID}"
//	                     "objects": "{number of objects}",                 // stored in DVO_REPORT table, computed as SELECT count(distinct(object)) WHERE cluster="{cluster UUID}" and namespace="{namespace UUID}"
//	                     "reported_at": "{reported_at}",                   // stored in DVO_REPORT table
//	                     "last_checked_at": "{last_checked_at}",           // stored in DVO_REPORT table
//	                     "highest_severity": "{highest_severity}",         // computed with the help of Content Service
//	                 },
//	             },
//		    ]
//		}
func (server *HTTPServer) allDVONamespaces(writer http.ResponseWriter, _ *http.Request) {
	log.Info().Msg("All DVO namespaces handler")
	// set the response header
	writer.Header().Set(contentType, appJSON)

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
			MetadataEntry{
				Recommendations: 100,
				Objects:         1000,
				ReportedAt:      time.Now().Format(time.RFC3339),
				LastCheckedAt:   time.Now().Format(time.RFC3339),
				HighestSeverity: 5,
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

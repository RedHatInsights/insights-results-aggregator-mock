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

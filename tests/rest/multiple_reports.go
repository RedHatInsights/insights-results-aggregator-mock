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

// MultipleReportsResponse represents response from the server that contains
// results for multiple clusters together with overall status
type MultipleReportsResponse struct {
	Clusters    []string               `json:"clusters"`
	Errors      []string               `json:"errors"`
	Reports     map[string]interface{} `json:"reports"`
	GeneratedAt string                 `json:"generated_at"`
	Status      string                 `json:"status"`
}

// ClusterListInRequest represents request body containing list of clusters
type ClusterListInRequest struct {
	Clusters []string `json:"clusters"`
}

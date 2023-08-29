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

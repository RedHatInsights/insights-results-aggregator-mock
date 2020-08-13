/*
Copyright © 2020 Red Hat, Inc.

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

// Generated documentation is available at:
// https://godoc.org/github.com/RedHatInsights/insights-results-aggregator-mock/groups
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-results-aggregator-mock/packages/groups/groups.html
package groups

import (
	"io/ioutil"
	"path/filepath"

	"github.com/go-yaml/yaml"
	"github.com/rs/zerolog/log"
)

// Group represent the relative information about a group
type Group struct {
	Name        string   `yaml:"name" json:"title"`
	Description string   `yaml:"description" json:"description"`
	Tags        []string `yaml:"tags" json:"tags"`
}

// ParseGroupConfigFile parses the groups configuration file and return the read groups
func ParseGroupConfigFile(groupConfigPath string) (map[string]Group, error) {
	configBytes, err := ioutil.ReadFile(filepath.Clean(groupConfigPath))
	if err != nil {
		log.Error().Err(err).Msg("Error reading groups configuration file")
		return nil, err
	}

	var groups map[string]Group

	err = yaml.Unmarshal(configBytes, &groups)

	if err != nil {
		log.Error().Err(err).Msg("Error parsing groups configuration file")
		return nil, err
	}

	return groups, nil
}

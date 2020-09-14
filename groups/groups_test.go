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

package groups_test

import (
	"testing"

	"github.com/RedHatInsights/insights-results-aggregator-mock/groups"
)

// TestParseGroupConfigFileNonExistingFile check whether non existing file is detected properly
func TestParseGroupConfigFileNonExistingFile(t *testing.T) {
	// the following file does not exist
	_, err := groups.ParseGroupConfigFile("this does not exist")
	if err == nil {
		t.Log(err)
		t.Fatal("Error should be returned for non existing file")
	}
}

// TestParseGroupConfigFileNonExistingFile check whether non YAML file is detected properly
func TestParseGroupConfigFileNotYamlFile(t *testing.T) {
	// the following file does exist, but it is not proper YAML file
	_, err := groups.ParseGroupConfigFile("../LICENSE")
	if err == nil {
		t.Log(err)
		t.Fatal("Error should be returned for improper file format")
	}
}

// TestParseGroupConfigFileProperYamlFile is basic test for checking whether group configuration file can be read properly
func TestParseGroupConfigFileProperYamlFile(t *testing.T) {
	// the following file does exist, but it is not proper YAML file
	_, err := groups.ParseGroupConfigFile("../groups_config.yaml")
	if err != nil {
		t.Log(err)
		t.Fatal("Error should not be returned for existing and proper file")
	}
	// TODO: more checks will need test_config.yaml
}

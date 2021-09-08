// Copyright 2021 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package content

import (
	"io/ioutil"
	"os"

	"encoding/json"
)

func ParseContent(filePath string) ([]RuleContent, error) {
	var ruleContent []RuleContent

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return ruleContent, err
	}

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return ruleContent, err
	}

	err = jsonFile.Close()
	if err != nil {
		return ruleContent, err
	}

	err = json.Unmarshal(bytes, &ruleContent)
	if err != nil {
		return ruleContent, err
	}

	return ruleContent, nil
}

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
	"time"

	"github.com/verdverm/frisby"
)

// AckListResponse represents response containing list of acks
type AckListResponse struct {
	AckListMetaData AckListMetadata `json:"meta"`
	Acks            []Ack           `json:"data"`
}

// AckListMetadata represents metadata about rule acks
type AckListMetadata struct {
	Count int `json:"count"`
}

// Ack represents one rule Ack
type Ack struct {
	Rule          string    `json:"rule"`
	Justification string    `json:"justification"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// RuleAckRequest represents structure used to request rule to be acked via
// REST API
type RuleAckRequest struct {
	Rule          string `json:"rule_id"`
	Justification string `json:"justification"`
}

// Justification represents structure with justification for acking a rule
type Justification struct {
	Justification string `json:"justification"`
}

// ackListEndpoint constructs an URL for list of acks
func ackListEndpoint() string {
	return fmt.Sprintf("%sack", apiURL)
}

// ackRuleEndpoint constructs an URL for getting/changing/removing rule ack
func ackRuleEndpoint(ruleSelector string) string {
	return fmt.Sprintf("%sack/%s", apiURL, ruleSelector)
}

func testRuleExistence(f *frisby.Frisby, acks []Ack, searchedRule string) {
	for _, ack := range acks {
		if ack.Rule == searchedRule {
			// found it
			return
		}
	}
	// not found
	errorMessage := fmt.Sprintf("Rule %s has not been found in list of acks", searchedRule)
	f.AddError(errorMessage)
}

// ListOfAcks checks if the 'ack' point responds correctly to HTTP GET command
func checkRetrieveListOfAcks() {
	url := ackListEndpoint()
	f := frisby.Create("Check the 'ack' REST API point using HTTP GET method").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusOK)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := AckListResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.AckListMetaData.Count != 5 {
			f.AddError("Improper metadata about number of acks")
		}
		if len(response.Acks) != 5 {
			f.AddError("Improper number of acks returned")
		}
		// try few rule acks
		testRuleExistence(f, response.Acks, "ccx_rules_ocp.external.rules.nodes_requirements_check.report|NODES_MINIMUM_REQUIREMENTS_NOT_MET")
		testRuleExistence(f, response.Acks, "ccx_rules_ocp.external.bug_rules.bug_1766907.report|BUGZILLA_BUG_1766907")
		testRuleExistence(f, response.Acks, "ccx_rules_ocp.external.rules.nodes_kubelet_version_check.report|NODE_KUBELET_VERSION")
		testRuleExistence(f, response.Acks, "ccx_rules_ocp.external.rules.samples_op_failed_image_import_check.report|SAMPLES_FAILED_IMAGE_IMPORT_ERR")
		testRuleExistence(f, response.Acks, "ccx_rules_ocp.external.rules.cluster_wide_proxy_auth_check.report|AUTH_OPERATOR_PROXY_ERROR")
	}
	f.PrintReport()
}

// checkAckRule checks if it is possible to ack one selected rule
func checkAckRule() {
	url := ackListEndpoint()
	f := frisby.Create("Check the 'ack' REST API point using HTTP POST method").Post(url)

	// set the payload to be sent
	f.SetJson(RuleAckRequest{
		Rule:          ackedRule1,
		Justification: "justification"})

	f.Send()
	f.ExpectStatus(http.StatusCreated)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := Ack{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
	}
	f.PrintReport()
}

// checkAckRuleWithIncorrectName checks how improper rule name is handled by ack endpoint
func checkAckRuleWithIncorrectName() {
	url := ackListEndpoint()
	f := frisby.Create("Check the 'ack' REST API point using HTTP POST method (incorrect variant)").Post(url)

	// set the payload to be sent
	f.SetJson(RuleAckRequest{
		Rule:          "incorrect rule name",
		Justification: "justification"})

	f.Send()
	f.ExpectStatus(http.StatusBadRequest)
	f.PrintReport()
}

// checkAckRuleViaGetEndpoint checks if it is possible to ack one selected rule
func checkAckRuleViaGetEndpoint() {
	url := ackRuleEndpoint(ackedRule2)
	f := frisby.Create("Check the 'ack/{rule_selector}' REST API point using HTTP GET method").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusOK)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := Ack{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
	}
	f.PrintReport()
}

// checkAckIncorrectRule checks if/how incorrect rule ack can be Ackd
func checkAckIncorrectRule() {
	url := ackRuleEndpoint(incorrectRule)
	f := frisby.Create("Check the 'ack/{rule_selector}' REST API point using HTTP GET method for incorrect rule").Get(url)
	f.Send()
	f.ExpectStatus(http.StatusBadRequest)
	f.PrintReport()
}

// checkUpdateExistingRule checks if it is possible to update ack for one selected rule
func checkUpdateExistingRule() {
	url := ackRuleEndpoint(ackedRule2)
	f := frisby.Create("Check the 'ack/{rule_selector}' REST API point using HTTP PUT method").Put(url)
	// set the payload to be sent
	f.SetJson(Justification{
		Justification: "justification"})

	f.Send()
	f.ExpectStatus(http.StatusOK)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := Ack{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
	}
	f.PrintReport()
}

// checkUpdateNonExistingRule checks if/how non-existing rule ack can be updated
func checkUpdateNonExistingRule() {
	url := ackRuleEndpoint(nonExistingRule)
	f := frisby.Create("Check the 'ack/{rule_selector}' REST API point using HTTP PUT method for non-existing rule").Put(url)
	// set the payload to be sent
	f.SetJson(Justification{
		Justification: "justification"})

	f.Send()
	f.ExpectStatus(http.StatusNotFound)
	f.PrintReport()
}

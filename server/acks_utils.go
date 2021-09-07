/*
Copyright © 2021 Red Hat, Inc.

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

// Helper functions to be called from request handlers defined in the source
// file acks_handlers.go.

import (
	"errors"
	"time"

	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/RedHatInsights/insights-results-aggregator-mock/types"
)

// Log error messages
const (
	unableToReadRuleSelector = "Unable to read rule selector"
)

// formattedNow function returns current time formatted according to RFC3339
func formattedNow() string {
	return time.Now().Format(time.RFC3339)
}

// handleImproperSelector function handles situation when rule selector can not
// be read from client request. HTTP code 400 Bad Request is returned in this
// situation.
func handleImproperSelector(writer http.ResponseWriter, err error) {
	log.Error().Err(err).Msg(unableToReadRuleSelector)
	// return HTTP code 400 Bad Request
	http.Error(writer, err.Error(), http.StatusBadRequest)
}

// handleMissingRule function handles situation when rule acknowledge can not
// be found in internal data structure.  HTTP code 404 Not Found is returned in
// this situation.
func handleMissingRule(writer http.ResponseWriter, ruleSelector string) {
	err := errors.New("rule not found -> justification can not be changed")
	log.Error().Err(err).Msg("")
	// return 404
	http.Error(writer, err.Error(), http.StatusNotFound)
}

// addNewRule function add a new rule to existing map of acknowledges.
func addNewRule(ruleSelector types.RuleSelector, justification string, createdBy string) {
	// add new rule
	acks[ruleSelector] = types.Acknowledge{
		Acknowledged:  true,
		Rule:          string(ruleSelector),
		Justification: justification,
		CreatedBy:     createdBy,
		CreatedAt:     formattedNow(),
		UpdatedAt:     formattedNow(),
	}
}

// updateRuleJustification function updates justification of given rule. It
// also changes UpdatedAt attribute.
func updateRuleJustification(ruleSelector types.RuleSelector, justification types.AcknowledgementJustification) {
	// (it is impossible to change the struct in a map directly!)
	ack := acks[ruleSelector]

	// TODO: ask if that attribute needs to be updated as well
	// ack.CreatedBy = defaultUserName
	ack.UpdatedAt = formattedNow()
	ack.Justification = justification.Value

	// update map
	acks[ruleSelector] = ack
}

// updateRuleUpdatedAt function just UpdatedAt attribute.
func updateRuleUpdatedAt(ruleSelector types.RuleSelector) {
	// (it is impossible to change the struct in a map directly!)
	ack := acks[ruleSelector]

	// TODO: ask if that attribute needs to be updated as well
	// ack.CreatedBy = defaultUserName
	ack.UpdatedAt = formattedNow()

	// update map
	acks[ruleSelector] = ack
}

// returnRuleAckToClient returns information about selected rule ack to client.
// This function also tries to process all errors.
func returnRuleAckToClient(writer http.ResponseWriter, ack types.Acknowledge) {
	// serialize the above data structure into JSON format
	bytes, err := json.MarshalIndent(ack, "", "\t")
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
		return
	}

	// and send the serialized structure to client
	_, err = writer.Write(bytes)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

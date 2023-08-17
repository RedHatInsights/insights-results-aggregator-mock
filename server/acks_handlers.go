/*
Copyright © 2021, 2022 Red Hat, Inc.

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

	"github.com/RedHatInsights/insights-results-aggregator-mock/types"
)

// HTTP response-related constants
const (
	contentType = "Content-Type"
	appJSON     = "application/json; charset=utf-8"
)

// other constants
const (
	defaultUserName      = "onlineTester"
	defaultJustification = "?"
)

// method readAckList list acks from this account where the rule is active.
// Will return an empty list if this account has no acks.
//
// Response format should look like:
//
//	{
//	  "meta": {
//	    "count": 0
//	  },
//	  "links": {
//	    "first": "string",
//	    "previous": "string",
//	    "next": "string",
//	    "last": "string"
//	  },
//	  "data": [
//	    {
//	      "rule": "string",
//	      "justification": "string",
//	      "created_by": "string",
//	      "created_at": "2021-09-04T17:11:35.130Z",
//	      "updated_at": "2021-09-04T17:11:35.130Z"
//	    }
//	  ]
//	}
//
// Please note that for the sake of simplicity we don't use links section as
// pagination is not supported ATM.
func (server *HTTPServer) readAckList(writer http.ResponseWriter, _ *http.Request) {
	// set the response header
	writer.Header().Set(contentType, appJSON)

	var responseBody types.AcknowledgementsResponse

	// fill-in metadata part of response body
	responseBody.Metadata.Count = len(acks)

	// fill-in data part of response body
	responseBody.Data = make([]types.Acknowledge, len(acks))

	i := 0
	for _, ack := range acks {
		responseBody.Data[i] = ack
		i++
	}

	// serialize the above data structure into JSON format
	bytes, err := json.MarshalIndent(responseBody, "", "\t")
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

// method acknowledgePost acknowledges (and therefore hides) a rule from view
// in an account. If there's already an acknowledgement of this rule by this
// account, then return that. Otherwise, a new ack is created.
//
// An example request:
//
//	{
//	  "rule_id": "string",
//	  "justification": "string"
//	}
//
// An example response:
//
//	{
//	  "rule": "string",
//	  "justification": "string",  <- can not be set by this call!!!
//	  "created_by": "string",
//	  "created_at": "2021-09-04T17:52:48.976Z",
//	  "updated_at": "2021-09-04T17:52:48.976Z"
//	}
func (server *HTTPServer) acknowledgePost(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(contentType, appJSON)

	// try to read request body
	var parameters types.AcknowledgementRuleSelectorJustification
	err := json.NewDecoder(request.Body).Decode(&parameters)

	if err != nil {
		log.Error().Err(err).Msg("wrong payload provided by client")
		// return HTTP code 400 to client
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// we seem to have proper data -> let's display them
	log.Info().
		Str("rule", string(parameters.RuleSelector)).
		Str("value", parameters.Value).
		Msg("Proper payload provided")

	// check if rule selector has the proper format
	_, _, err = parseRuleSelector(parameters.RuleSelector)
	if err != nil {
		log.Error().Err(err).Msg("improper rule selector format")
		// return 400
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// try to find the rule in map of active rules
	_, found := acks[parameters.RuleSelector]
	if !found {
		// rule not found -> add a new one
		addNewRule(parameters.RuleSelector, parameters.Value, defaultUserName)
		// update HTTP status code accordingly
		writer.WriteHeader(http.StatusCreated)
	}

	// return existing rule or the new one (if created)
	ack := acks[parameters.RuleSelector]
	returnRuleAckToClient(writer, &ack)
}

// method getAcknowledge acknowledges (and therefore hides) a rule from view in
// an account. This view handles listing, retrieving, creating and deleting
// acks. Acks are created and deleted by Insights rule ID, not by their own ack
// ID.
//
// An example response:
//
//	{
//	  "rule": "string",
//	  "justification": "string",  <- can not be set by this call!!!
//	  "created_by": "string",
//	  "created_at": "2021-09-04T17:52:48.976Z",
//	  "updated_at": "2021-09-04T17:52:48.976Z"
//	}
//
// Please note that it is impossible to set "justification" field into any
// value that makes sense!
func (server *HTTPServer) getAcknowledge(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(contentType, appJSON)

	// read the selector
	ruleSelector, err := readRuleSelector(writer, request)

	// check for missing/improper selector
	if err != nil {
		handleImproperSelector(writer, err)
		// everything has been handled already
		return
	}

	// try to find the rule in map of active rules
	_, found := acks[ruleSelector]
	if !found {
		// rule not found -> add a new one
		addNewRule(ruleSelector, defaultJustification, defaultUserName)
	} else {
		// rule has been found -> just update it
		updateRuleUpdatedAt(ruleSelector)
	}

	ack := acks[ruleSelector]
	returnRuleAckToClient(writer, &ack)
}

// method updateAcknowledge updates an acknowledgement for a rule, by rule ID.
// A new justification can be supplied. The username is taken from the
// authenticated request. The updated ack is returned.
//
// An example of request:
//
//	{
//	   "justification": "string"
//	}
//
// An example response:
//
//	{
//	  "rule": "string",
//	  "justification": "string",
//	  "created_by": "string",
//	  "created_at": "2021-09-04T17:52:48.976Z",
//	  "updated_at": "2021-09-04T17:52:48.976Z"
//	}
//
// Additionally, if rule is not found, 404 is returned (not mentioned in
// original REST API specification).
func (server *HTTPServer) updateAcknowledge(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(contentType, appJSON)

	// read the selector
	ruleSelector, err := readRuleSelector(writer, request)

	// check for missing/improper selector
	if err != nil {
		handleImproperSelector(writer, err)
		// everything has been handled already
		return
	}

	// try to find the rule in map of active rules
	_, found := acks[ruleSelector]
	if !found {
		handleMissingRule(writer, string(ruleSelector))
		// everything has been handled already
		return
	}

	// try to read request body
	var justification types.AcknowledgementJustification
	err = json.NewDecoder(request.Body).Decode(&justification)

	if err != nil {
		log.Error().Err(err).Msg("justification provided by client")
		// return 400
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info().
		Str("value", justification.Value).
		Msg("Justification provided")

	// update existing rule
	updateRuleJustification(ruleSelector, justification)

	ack := acks[ruleSelector]
	returnRuleAckToClient(writer, &ack)
}

// method deleteAcknowledge deletes an acknowledgement for a rule, by its rule
// ID. If the ack existed, it is deleted and a 204 is returned. Otherwise, a
// 404 is returned.
func (server *HTTPServer) deleteAcknowledge(writer http.ResponseWriter, request *http.Request) {
	// read the selector
	ruleSelector, err := readRuleSelector(writer, request)

	// check for missing/improper selector
	if err != nil {
		handleImproperSelector(writer, err)
		// everything has been handled already
		return
	}

	// try to find the rule in map of active rules
	_, found := acks[ruleSelector]
	if !found {
		// return 404
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// rule has been acknowledget -> we can delete it
	delete(acks, ruleSelector)

	// return 204 -> rule ack has been deleted
	writer.WriteHeader(http.StatusNoContent)
}

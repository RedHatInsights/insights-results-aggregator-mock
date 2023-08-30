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

const (
	clusterWithPositiveRisksPrediction = "00000001-624a-49a5-bab8-4fdc5e51a266"
	clusterWithNegativeRisksPrediction = "00000003-eeee-eeee-eeee-000000000001"
	clusterWithNoContent               = "6cab9726-c2be-438e-af11-db846a678abb"
	clusterReturningUnavaliableService = "c60ba611-6af4-4d62-9b9e-36344da5e7bc"
	clusterReturningNotFound           = "234ec1a1-4679-4122-aacb-f0ae9f9e1a56"
	clusterWithImproperName            = "hey this is not proper"
)

// URPResponse structure represents response payload returned from URP endpoint
type URPResponse struct {
	Status             string             `json:"status"`
	Meta               URPMetadata        `json:"meta"`
	URPRecommendations URPRecommendations `json:"upgrade_recommendation"`
}

// URPMetadata structure represents sub-node in response payload returned from
// URP endpoint
type URPMetadata struct {
	LastCheckedAt time.Time `json:"last_checked_at"`
}

// URPRecommendations structure represents sub-node in response payload
// returned from URP endpoint
type URPRecommendations struct {
	UpgradeRecommended     bool           `json:"upgrade_recommended"`
	UpgradeRiskPredictions URPPredictions `json:"upgrade_risks_predictors"`
}

// URPPredictions structure represends sub-node in URP recommendation
type URPPredictions struct {
	Alerts             []Alert             `json:"alerts"`
	OperatorConditions []OperatorCondition `json:"operator_conditions"`
}

// Alert represents one entry in URP predictions
type Alert struct {
	Name      string `json:"name"`
	NameSpace string `json:"namespace"`
	Severity  string `json:"severity"`
	URL       string `json:"url"`
}

// OperatorCondition represents one entry in URP predictions
type OperatorCondition struct {
	Name      string `json:"name"`
	Condition string `json:"condition"`
	Reason    string `json:"reason"`
	URL       string `json:"url"`
}

// access the endpoint to retrieve results for given list of clusters using
// POST method
func constructURLUpgradeRiskEndpoint(cluster string) string {
	return fmt.Sprintf("%scluster/%s/upgrade-risks-prediction", apiURL, cluster)
}

// checkUpgradeRiskEndpointWithClusterWithPositiveRiskPrediction check how/if
// URP endpoint returns positive risk prediction for given cluster
func checkUpgradeRiskEndpointWithClusterWithPositiveRiskPrediction() {
	url := constructURLUpgradeRiskEndpoint(clusterWithPositiveRisksPrediction)

	// send request to the endpoint
	f := frisby.Create("Check the endpoint to return upgrade risk predictions for cluster with positive risk prediction").Get(url)
	f.Send()

	// check the response from server
	f.ExpectStatus(http.StatusOK)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response payload
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := URPResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.Status != "ok" {
			f.AddError(statusShouldBeSetToOK)
		}
		if !response.URPRecommendations.UpgradeRecommended {
			f.AddError("Upgrade should be recommended")
		}
		if len(response.URPRecommendations.UpgradeRiskPredictions.Alerts) != 0 {
			f.AddError("Zero alerts are expected")
		}
		if len(response.URPRecommendations.UpgradeRiskPredictions.OperatorConditions) != 0 {
			f.AddError("Zero operator conditions are expected")
		}
	}

	f.PrintReport()
}

// checkUpgradeRiskEndpointWithClusterWithPositiveRiskPrediction check how/if
// URP endpoint returns negative risk prediction for given cluster
func checkUpgradeRiskEndpointWithClusterWithNegativeRiskPrediction() {
	url := constructURLUpgradeRiskEndpoint(clusterWithNegativeRisksPrediction)

	// send request to the endpoint
	f := frisby.Create("Check the endpoint to return upgrade risk predictions for cluster with negative risk prediction").Get(url)
	f.Send()

	// check the response from server
	f.ExpectStatus(http.StatusOK)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response payload
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := URPResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.Status != "ok" {
			f.AddError(statusShouldBeSetToOK)
		}
		if response.URPRecommendations.UpgradeRecommended {
			f.AddError("Upgrade should NOT be recommended")
		}
		alerts := response.URPRecommendations.UpgradeRiskPredictions.Alerts
		if len(alerts) != 3 {
			f.AddError("Exactly 3 alerts are expected")
		}
		conditions := response.URPRecommendations.UpgradeRiskPredictions.OperatorConditions
		if len(conditions) != 4 {
			f.AddError("Exactly 4 operator conditions are expected")
		}
		// TODO: exact check of alerts
		// TODO: exact check of conditions
	}

	f.PrintReport()
}

// checkUpgradeRiskEndpointWithClusterWithNoContent check how/if URP endpoint
// returns NoContent status for selected cluster
func checkUpgradeRiskEndpointWithClusterWithNoContent() {
	url := constructURLUpgradeRiskEndpoint(clusterWithNoContent)

	// send request to the endpoint
	f := frisby.Create("Check the endpoint to return upgrade risk predictions for cluster with no content").Get(url)
	f.Send()

	// check the response from server
	f.ExpectStatus(http.StatusNoContent)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// nothing more should be checked there

	f.PrintReport()
}

// checkUpgradeRiskEndpointUnavailableServiceCase check how/if URP endpoint
// returns ServiceUnavailable status for selected cluster
func checkUpgradeRiskEndpointUnavailableServiceCase() {
	url := constructURLUpgradeRiskEndpoint(clusterReturningUnavaliableService)

	// send request to the endpoint
	f := frisby.Create("Check the endpoint to return upgrade risk predictions for cluster returning unavailable service").Get(url)
	f.Send()

	// check the response from server
	f.ExpectStatus(http.StatusServiceUnavailable)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response payload
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := URPResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.Status != "AMS service unavailable" {
			f.AddError("Unexpected status: " + response.Status)
		}
	}

	f.PrintReport()
}

// checkUpgradeRiskEndpointNotFoundCase check how/if URP endpoint returns
// NotFound status for selected cluster
func checkUpgradeRiskEndpointNotFoundCase() {
	url := constructURLUpgradeRiskEndpoint(clusterReturningNotFound)

	// send request to the endpoint
	f := frisby.Create("Check the endpoint to return upgrade risk predictions for cluster returning NotFound").Get(url)
	f.Send()

	// check the response from server
	f.ExpectStatus(http.StatusNotFound)
	f.ExpectHeader(contentTypeHeader, ContentTypeJSON)

	// check the response payload
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		response := URPResponse{}
		err := json.Unmarshal(text, &response)
		if err != nil {
			f.AddError(err.Error())
		}
		if response.Status != "No data for the cluster" {
			f.AddError("Unexpected status: " + response.Status)
		}
	}

	f.PrintReport()
}

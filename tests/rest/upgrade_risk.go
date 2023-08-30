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

import "time"

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

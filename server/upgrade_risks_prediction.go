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

package server

import (
	"net/http"

	"github.com/RedHatInsights/insights-operator-utils/responses"
	"github.com/rs/zerolog/log"
)

const (
	// ClusterOk is the cluster name for a OK response with no upgrade risks detected
	ClusterOk = "00000001-624a-49a5-bab8-4fdc5e51a266"
	// ClusterOkFailUpgrade is the cluster name for a OK response with upgrade risks detected
	ClusterOkFailUpgrade = "00000003-eeee-eeee-eeee-000000000001"
	// ClusterManaged is the cluster name for the response when a cluster in "managed"
	ClusterManaged = "6cab9726-c2be-438e-af11-db846a678abb"
	// ClusterNoAMS is the cluster name for the response when the AMS API is not available
	ClusterNoAMS = "c60ba611-6af4-4d62-9b9e-36344da5e7bc"
	// ClusterUnavailable is the cluster name for the response when the Upgrade risks prediction service is unavailable
	ClusterUnavailable = "897ec1a1-4679-4122-aacb-f0ae9f9e1a5f"
)

// method upgradeRisksPrediction return a recommendation to upgrade or not a cluster
// and a list of the alerts/operator conditions that were taken into account if the
// upgrade is not recommended.
//
// Response format should look like:
//
//	{
//		"upgrade_risks_predictors": {
//			"alerts": ["alert1", "alert2"],
//			"operator_conditions": ["foc1", "foc2"]
//		}
//	}
func (server *HTTPServer) upgradeRisksPrediction(writer http.ResponseWriter, request *http.Request) {
	clusterName, err := readClusterName(writer, request)
	if err != nil {
		// everything has been handled already
		return
	}

	switch clusterName {
	case ClusterManaged:
		log.Info().Msg("managed cluster case")
		err = responses.SendNoContent(writer)
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}

	case ClusterNoAMS:
		log.Info().Msg("No AMS available case")
		err = responses.SendServiceUnavailable(writer, "AMS service unavailable")
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}

	case ClusterUnavailable:
		log.Info().Msg("No Upgrade Risks Prediction service available case")
		err = responses.SendServiceUnavailable(writer, "Upgrade Risks Prediction service unavailable")
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}

	default:
		prediction, err := server.Storage.GetPredictionForCluster(clusterName)
		if err != nil {
			log.Error().Err(err).Msg("error retrieving upgrade prediction from storage")
			handleServerError(err)
			err = responses.SendNotFound(writer, err.Error())
			if err != nil {
				log.Error().Err(err).Msg(responseDataError)
			}
			return
		}

		if clusterName == ClusterOkFailUpgrade {
			prediction.Predictors.Alerts = append(prediction.Predictors.Alerts, "alert1", "alert2")
			prediction.Predictors.OperatorConditions = append(prediction.Predictors.OperatorConditions, "foc1", "foc2")
		}

		writer.Header().Set(contentType, appJSON)
		err = responses.SendOK(writer, responses.BuildOkResponseWithData("upgrade_recommendation", prediction))
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}
	}
}

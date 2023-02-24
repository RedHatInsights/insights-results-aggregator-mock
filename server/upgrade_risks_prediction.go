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
	"time"

	"github.com/RedHatInsights/insights-operator-utils/responses"
	"github.com/rs/zerolog/log"
)

// method upgradeRisksPrediction return a recommendation to upgrade or not a cluster
// and a list of the alerts/operator conditions that were taken into account if the
// upgrade is not recommended.
//
// Response format should look like:
//
//	{
//		"upgrade_recommended": false,
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

	meta := make(map[string]interface{})
	meta["last_checked_at"] = time.Now().UTC().Format(time.RFC3339)

	response := make(map[string]interface{})
	response["upgrade_recommendation"] = prediction
	response["meta"] = meta
	response["status"] = "ok"

	writer.Header().Set(contentType, appJSON)

	err = responses.SendOK(writer, response)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

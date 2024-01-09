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

	httputils "github.com/RedHatInsights/insights-operator-utils/http"
	"github.com/RedHatInsights/insights-operator-utils/responses"
	"github.com/rs/zerolog/log"

	"github.com/RedHatInsights/insights-results-aggregator-mock/types"
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
	// ClusterNoData is the cluster name for the response when the Upgrade risks prediction service returns a 404
	ClusterNoData = "234ec1a1-4679-4122-aacb-f0ae9f9e1a56"
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
//			"alerts": [
//				{
//					"name": "APIRemovedInNextEUSReleaseInUse",
//					"namespace": "openshift-kube-apiserver",
//					"severity": "info",
//                  "url": "${CONSOLE_URL}/monitoring/alerts?orderBy=asc&sortBy=Severity&alert-name=${ALERT_NAME}"
//				}
//			],
//			"operator_conditions": [
//				{
//					"name": "authentication",
//					"condition": "Failing",
//					"reason": "AsExpected",
//                  "url": "${CONSOLE_URL}/k8s/cluster/config.openshift.io~v1~ClusterOperator/${OPERATOR_NAME}"
//				}
//			]
//		}
//	}

func (server *HTTPServer) upgradeRisksPrediction(writer http.ResponseWriter, request *http.Request) {
	clusterName, err := readClusterName(writer, request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
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

	case ClusterNoData:
		log.Info().Msg("No data for the cluster")
		err = responses.SendNotFound(writer, "No data for the cluster")
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}

	default:
		server.sendOkResponse(clusterName, writer)
	}
}

func (server *HTTPServer) upgradeRisksPredictionMultiCluster(writer http.ResponseWriter, request *http.Request) {
	// try to read list of cluster IDs
	clusterList, successful := httputils.ReadClusterListFromBody(writer, request)
	if !successful {
		// wrong state has been handled already
		log.Error().Msg("No cluster list provided")
		return
	}

	responseArray := []types.ClusterUpgradeRiskPrediction{}
	for _, cluster := range clusterList {
		clusterName := types.ClusterName(cluster)

		switch clusterName {
		case ClusterManaged:
			log.Info().Msg("managed cluster case")
			responseArray = append(responseArray, types.ClusterUpgradeRiskPrediction{
				Cluster: cluster,
				Status:  "ok",
			})

		case ClusterNoAMS:
			log.Info().Msg("No AMS available case")
			responseArray = append(responseArray, types.ClusterUpgradeRiskPrediction{
				Cluster: cluster,
				Status:  "AMS service not available",
			})

		case ClusterUnavailable:
			log.Info().Msg("No Upgrade Risks Prediction service available case")
			responseArray = append(responseArray, types.ClusterUpgradeRiskPrediction{
				Cluster: cluster,
				Status:  "Upgrade Risks Prediction service unavailable",
			})

		case ClusterOkFailUpgrade:
			log.Info().Msg("Cluster is not recommended to upgrade")
			responseArray = append(responseArray, types.ClusterUpgradeRiskPrediction{
				Cluster:     cluster,
				Status:      "ok",
				Recommended: false,
				Predictors:  buildNotEmptyPredictors(),
			})

		case ClusterOk:
			responseArray = append(responseArray, types.ClusterUpgradeRiskPrediction{
				Cluster:     cluster,
				Status:      "ok",
				Recommended: true,
				Predictors: types.UpgradeRisksPredictors{
					Alerts:             []types.Alert{},
					OperatorConditions: []types.OperatorCondition{},
				},
			})

		case ClusterNoData:
		default:
			log.Info().Msg("No data for the cluster")
			responseArray = append(responseArray, types.ClusterUpgradeRiskPrediction{
				Cluster: cluster,
				Status:  "No data for the cluster",
			})
		}
	}

	err := responses.SendOK(
		writer,
		map[string]interface{}{
			"predictions": responseArray,
			"status":      "ok",
		},
	)

	if err != nil {
		log.Error().Err(err).Msg("Error sending response")
		handleServerError(err)
	}
}

func (server *HTTPServer) sendOkResponse(clusterName types.ClusterName, writer http.ResponseWriter) {
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
		buildNotRecommendedPrediction(prediction)
	}

	writer.Header().Set(contentType, appJSON)
	resp := responses.BuildOkResponseWithData("upgrade_recommendation", prediction)
	resp["meta"] = map[string]string{
		"last_checked_at": time.Now().UTC().Format(time.RFC3339),
	}
	err = responses.SendOK(writer, resp)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

func buildNotRecommendedPrediction(prediction *types.UpgradeRiskPrediction) {
	prediction.Recommended = false
	prediction.Predictors = buildNotEmptyPredictors()
}

func buildNotEmptyPredictors() types.UpgradeRisksPredictors {
	return types.UpgradeRisksPredictors{
		Alerts: []types.Alert{
			types.Alert{
				Name:      "alert1",
				Namespace: "namespace1",
				Severity:  "info",
				URL:       "https://my-cluster.com/monitoring/alerts?orderBy=asc&sortBy=Severity&alert-name=alert1",
			},
			types.Alert{
				Name:      "alert2",
				Namespace: "namespace2",
				Severity:  "warning",
				URL:       "https://my-cluster.com/monitoring/alerts?orderBy=asc&sortBy=Severity&alert-name=alert2",
			},
			types.Alert{
				Name:      "alert3",
				Namespace: "namespace3",
				Severity:  "critical",
				URL:       "https://my-cluster.com/monitoring/alerts?orderBy=asc&sortBy=Severity&alert-name=alert3",
			},
		},

		OperatorConditions: []types.OperatorCondition{
			types.OperatorCondition{
				Name:      "foc1",
				Condition: "Degraded",
				Reason:    "NotExpected",
				URL:       "https://my-cluster.com/k8s/cluster/config.openshift.io~v1~ClusterOperator/foc1",
			},
			types.OperatorCondition{
				Name:      "foc2",
				Condition: "Failing",
				Reason:    "NotExpected",
				URL:       "https://my-cluster.com/k8s/cluster/config.openshift.io~v1~ClusterOperator/foc2",
			},
			types.OperatorCondition{
				Name:      "foc3",
				Condition: "Not Available",
				Reason:    "NotExpected",
				URL:       "https://my-cluster.com/k8s/cluster/config.openshift.io~v1~ClusterOperator/foc3",
			},
			types.OperatorCondition{
				Name:      "foc4",
				Condition: "Not Upgradeable",
				Reason:    "NotExpected",
				URL:       "https://my-cluster.com/k8s/cluster/config.openshift.io~v1~ClusterOperator/foc4",
			},
		},
	}
}

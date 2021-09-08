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

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/RedHatInsights/insights-operator-utils/responses"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/RedHatInsights/insights-results-aggregator-mock/data"
	"github.com/RedHatInsights/insights-results-aggregator-mock/groups"
	"github.com/RedHatInsights/insights-results-aggregator-mock/types"
)

const failureClusterIDPrefix = "ffffffff-ffff-ffff-ffff-"

const unableToReadReportErrorMessage = "Unable to read report for cluster"

// readOrganizationID retrieves organization id from request
// if it's not possible, it writes http error to the writer and returns error
func readOrganizationID(writer http.ResponseWriter, request *http.Request) (types.OrgID, error) {
	organizationID, err := getRouterPositiveIntParam(request, "organization")
	if err != nil {
		return 0, err
	}
	return types.OrgID(organizationID), nil
}

// readRuleSelector retrieves rule selector from request
func readRuleSelector(writer http.ResponseWriter, request *http.Request) (types.RuleSelector, error) {
	ruleSelector, err := getRouterParam(request, "rule_selector")
	if err != nil {
		return "", err
	}

	// check if the rule selector seems to be correct
	_, _, err = parseRuleSelector(types.RuleSelector(ruleSelector))
	if err != nil {
		return "", err
	}
	return types.RuleSelector(ruleSelector), nil
}

// readClusterName retrieves cluster name from request
// if it's not possible, it writes http error to the writer and returns error
func readClusterName(writer http.ResponseWriter, request *http.Request) (types.ClusterName, error) {
	clusterName, err := getRouterParam(request, "cluster")
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}
	return types.ClusterName(clusterName), nil
}

// getRouterParam retrieves parameter from URL like `/organization/{org_id}`
func getRouterParam(request *http.Request, paramName string) (string, error) {
	value, found := mux.Vars(request)[paramName]
	if !found {
		return "", errors.New("Missing param")
	}

	return value, nil
}

// getRouterPositiveIntParam retrieves parameter from URL like `/organization/{org_id}`
// and check it for being valid and positive integer, otherwise returns error
func getRouterPositiveIntParam(request *http.Request, paramName string) (uint64, error) {
	value, err := getRouterParam(request, paramName)
	if err != nil {
		return 0, err
	}

	uintValue, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}

	if uintValue == 0 {
		return 0, err
	}

	return uintValue, nil
}

// mainEndpoint will handle the requests for / endpoint
func (server *HTTPServer) mainEndpoint(writer http.ResponseWriter, _ *http.Request) {
	err := responses.SendOK(writer, responses.BuildOkResponse())
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

// serveAPISpecFile serves an OpenAPI specifications file specified in config file
func (server *HTTPServer) serveAPISpecFile(writer http.ResponseWriter, request *http.Request) {
	absPath, err := filepath.Abs(server.Config.APISpecFile)
	if err != nil {
		const message = "Error creating absolute path of OpenAPI spec file"
		log.Error().Err(err).Msg(message)
		handleServerError(err)
		return
	}

	http.ServeFile(writer, request, absPath)
}

// listOfGroups returns the list of defined groups
func (server *HTTPServer) listOfGroups(writer http.ResponseWriter, request *http.Request) {
	log.Info().Msg("List of groups handler")

	// let's mimick Content Service behaviour preciselly
	if server.groupsList == nil {
		log.Info().Msg("Initializing group list for the first time")
		server.groupsList = make([]groups.Group, 0, len(server.Groups))

		for _, group := range server.Groups {
			server.groupsList = append(server.groupsList, group)
		}
	}

	log.Info().Int("items", len(server.groupsList)).Msg("Group list")
	err := responses.SendOK(writer, responses.BuildOkResponseWithData("groups", server.groupsList))
	if err != nil {
		log.Error().Err(err)
		handleServerError(err)
		return
	}
}

func (server *HTTPServer) listOfOrganizations(writer http.ResponseWriter, _ *http.Request) {
	organizations, err := server.Storage.ListOfOrgs()
	if err != nil {
		log.Error().Err(err).Msg("Unable to get list of organizations")
		return
	}
	err = responses.SendOK(writer, responses.BuildOkResponseWithData("organizations", organizations))
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

func (server *HTTPServer) listOfClustersForOrganization(writer http.ResponseWriter, request *http.Request) {
	organizationID, err := readOrganizationID(writer, request)

	if err != nil {
		// everything has been handled already
		return
	}

	clusters, err := server.Storage.ListOfClustersForOrg(organizationID)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get list of clusters")
		handleServerError(err)
		err := responses.SendForbidden(writer, err.Error())
		if err != nil {
			log.Error().Err(err).Msg("Unable send forbidden response")
		}
		return
	}
	err = responses.SendOK(writer, responses.BuildOkResponseWithData("clusters", clusters))
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

func (server *HTTPServer) readReportForCluster(writer http.ResponseWriter, request *http.Request) {
	clusterName, err := readClusterName(writer, request)
	if err != nil {
		// everything has been handled already
		return
	}

	if strings.HasPrefix(string(clusterName), failureClusterIDPrefix) {
		s := string(clusterName)
		log.Info().Str("Cluster name", s).Msg("Failed clusters")
		suffix := s[len(s)-3:]
		code, err := strconv.Atoi(suffix)
		if err != nil {
			handleServerError(err)
			return
		}
		log.Info().Int("Code", int(code)).Msg("Failed clusters")
		writer.WriteHeader(code)
		return
	}
	report, err := server.Storage.ReadReportForCluster(clusterName)
	if err != nil {
		log.Error().Err(err).Msg(unableToReadReportErrorMessage)
		handleServerError(err)
		return
	}

	r := []byte(report)
	_, err = writer.Write(r)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

// ClusterList is a data structure that store list of cluster IDs (names).
type ClusterList struct {
	Clusters []string `json:"clusters"`
}

// ClusterReports is a data structure containing list of clusters, list of
// errors and dictionary with results per cluster.
type ClusterReports struct {
	ClusterList []types.ClusterName               `json:"clusters"`
	Errors      []types.ClusterName               `json:"errors"`
	Reports     map[types.ClusterName]interface{} `json:"reports"`
	GeneratedAt string                            `json:"generated_at"`
}

func (server *HTTPServer) readReportForAllClustersInOrg(writer http.ResponseWriter, request *http.Request) {
	organizationID, err := readOrganizationID(writer, request)

	if err != nil {
		// everything has been handled already
		return
	}
	log.Info().Int("OrgID", int(organizationID)).Msg("Organization ID to get list of results")

	var generatedReports ClusterReports
	generatedReports.GeneratedAt = time.Now().UTC().Format(time.RFC3339)

	generatedReports.Reports = make(map[types.ClusterName]interface{})

	bytes, err := json.MarshalIndent(generatedReports, "", "\t")
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
		return
	}
	_, err = writer.Write(bytes)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

func (server *HTTPServer) readReportForClusters(writer http.ResponseWriter, request *http.Request) {
	var clusterList ClusterList
	var generatedReports ClusterReports
	generatedReports.GeneratedAt = time.Now().UTC().Format(time.RFC3339)

	generatedReports.Reports = make(map[types.ClusterName]interface{})

	err := json.NewDecoder(request.Body).Decode(&clusterList)

	if err != nil {
		log.Error().Err(err).Msg("getting list of clusters")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	for _, clusterName := range clusterList.Clusters {
		log.Info().Str("cluster name", clusterName).Msg("result for cluster")
		clusterName := types.ClusterName(clusterName)
		reportStr, err := server.Storage.ReadReportForCluster(clusterName)
		if err != nil {
			log.Error().Err(err).Msg(unableToReadReportErrorMessage)
			generatedReports.Errors = append(generatedReports.Errors, clusterName)
			// if error happen, simply go to the next cluster
			continue
		}
		var report interface{}
		err = json.Unmarshal([]byte(reportStr), &report)
		if err != nil {
			log.Error().Err(err).Msg("Unable to unmarshal report for cluster")
			generatedReports.Errors = append(generatedReports.Errors, clusterName)
			// if error happen, simply go to the next cluster
			continue
		}
		generatedReports.ClusterList = append(generatedReports.ClusterList, clusterName)
		generatedReports.Reports[clusterName] = report
	}
	bytes, err := json.MarshalIndent(generatedReports, "", "\t")
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
		return
	}
	_, err = writer.Write(bytes)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

func (server *HTTPServer) readReportForOrganizationAndCluster(writer http.ResponseWriter, request *http.Request) {
	organizationID, err := readOrganizationID(writer, request)
	if err != nil {
		// everything has been handled already
		return
	}

	clusterName, err := readClusterName(writer, request)
	if err != nil {
		// everything has been handled already
		return
	}

	report, err := server.Storage.ReadReportForOrganizationAndCluster(organizationID, clusterName)
	if err != nil {
		log.Error().Err(err).Msg(unableToReadReportErrorMessage)
		handleServerError(err)
		return
	}

	r := []byte(report)
	_, err = writer.Write(r)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

func parseRuleSelector(ruleSelector types.RuleSelector) (types.Component, types.ErrorKey, error) {
	splitedRuleID := strings.Split(string(ruleSelector), "|")

	if len(splitedRuleID) != 2 {
		err := fmt.Errorf("invalid rule ID, it must contain only rule ID and error key separated by |")
		log.Error().Err(err)
		return types.Component(""), types.ErrorKey(""), err
	}

	IDValidator := regexp.MustCompile(`^[a-zA-Z_0-9.]+$`)

	isRuleIDValid := IDValidator.MatchString(splitedRuleID[0])
	isErrorKeyValid := IDValidator.MatchString(splitedRuleID[1])

	if !isRuleIDValid || !isErrorKeyValid {
		err := fmt.Errorf("invalid rule ID, each part of ID must contain only latin characters, number, underscores or dots")
		log.Error().Err(err)
		return types.Component(""), types.ErrorKey(""), err
	}

	return types.Component(splitedRuleID[0]), types.ErrorKey(splitedRuleID[1]), nil
}

func readClustersHittingRule(component types.Component, errorKey types.ErrorKey) []types.ClusterName {
	var clusterList []types.ClusterName

	// TODO: quick and dirty linear search should be imroved later if required
	for _, ruleHit := range data.RuleHits {
		if ruleHit.Component == component && ruleHit.ErrorKey == errorKey {
			clusterList = append(clusterList, ruleHit.Cluster)
		}
	}

	return clusterList
}

// HittingClustersMetadata used to store metadata of hitting clusters
type HittingClustersMetadata struct {
	Count       int             `json:"count"`
	Component   types.Component `json:"component"`
	ErrorKey    types.ErrorKey  `json:"error_key"`
	GeneratedAt string          `json:"generated_at"`
}

// HittingClusters is a data structure containing list of clusters
// hitting the given rule.
type HittingClusters struct {
	Metadata    HittingClustersMetadata `json:"meta"`
	ClusterList []types.ClusterName     `json:"data"`
}

// ruleClusterDetailEndpoint methods implements endpoint that should return a list of all the clusters IDs affected by this rule
func (server *HTTPServer) ruleClusterDetailEndpoint(writer http.ResponseWriter, request *http.Request) {
	// read the selector
	ruleSelector, err := readRuleSelector(writer, request)

	// check for missing/improper selector
	if err != nil {
		log.Error().Err(err).Msg("unable to read rule selector")
		// everything has been handled already
		return
	}

	// parse the selector
	component, errorKey, err := parseRuleSelector(ruleSelector)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse rule selector")
		// everything has been handled already
		return
	}

	// read all clusters hitting given rule
	log.Info().
		Str("component", string(component)).
		Str("error key", string(errorKey)).
		Msg("Reading clusters hitting given rule")
	clusters := readClustersHittingRule(component, errorKey)
	log.Info().Int("cluster count", len(clusters)).Msg("Clusters hitting the rule")

	// prepare response
	var hittingClusters HittingClusters

	// first fill-in metadata
	hittingClusters.Metadata.GeneratedAt = time.Now().UTC().Format(time.RFC3339)
	hittingClusters.Metadata.Count = len(clusters)
	hittingClusters.Metadata.Component = component
	hittingClusters.Metadata.ErrorKey = errorKey

	// second fill-in list of clusters
	hittingClusters.ClusterList = clusters

	bytes, err := json.MarshalIndent(hittingClusters, "", "\t")
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
		return
	}
	_, err = writer.Write(bytes)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

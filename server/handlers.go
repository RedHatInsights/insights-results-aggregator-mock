/*
Copyright © 2020, 2023 Red Hat, Inc.

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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/RedHatInsights/insights-operator-utils/responses"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/RedHatInsights/insights-results-aggregator-mock/data"
	"github.com/RedHatInsights/insights-results-aggregator-mock/groups"
	"github.com/RedHatInsights/insights-results-aggregator-mock/types"
)

const failureClusterIDPrefix = "ffffffff-ffff-ffff-ffff-"

const requestParameter = "Request parameter"

const unableToReadReportErrorMessage = "Unable to read report for cluster"
const requestsForClusterNotFound = "Requests for cluster not found"

// StatusProcessed is message returned for already processed reports (rule hits)
const StatusProcessed = "processed"

// readOrganizationID retrieves organization id from request
// if it's not possible, it writes http error to the writer and returns error
func readOrganizationID(_ http.ResponseWriter, request *http.Request) (types.OrgID, error) {
	organizationID, err := getRouterPositiveIntParam(request, "organization")
	if err != nil {
		return 0, err
	}
	return types.OrgID(organizationID), nil
}

// readRuleSelector retrieves rule selector from request
func readRuleSelector(_ http.ResponseWriter, request *http.Request) (types.RuleSelector, error) {
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

// ValidateClusterName checks that the cluster name is a valid UUID.
// Converted cluster name is returned if everything is okay, otherwise an error is returned.
func ValidateClusterName(clusterName string) (types.ClusterName, error) {
	if _, err := uuid.Parse(clusterName); err != nil {
		message := fmt.Sprintf("invalid cluster name: '%s'. Error: %s", clusterName, err.Error())

		log.Error().Err(err).Msg(message)
		return "", err
	}

	return types.ClusterName(clusterName), nil
}

// ValidateRequestID checks that the request ID has proper format.
// Converted request ID is returned if everything is okay, otherwise an error is returned.
func ValidateRequestID(requestID string) (types.RequestID, error) {
	IDValidator := regexp.MustCompile(`^[a-zA-Z0-9]+$`)

	if !IDValidator.MatchString(requestID) {
		message := fmt.Sprintf("invalid request ID: '%s'", requestID)
		err := errors.New(message)
		log.Error().Err(err).Msg(message)
		return "", err
	}

	return types.RequestID(requestID), nil
}

// readClusterName retrieves cluster name from request
// if it's not possible, it writes http error to the writer and returns error
func readClusterName(_ http.ResponseWriter, request *http.Request) (types.ClusterName, error) {
	clusterName, err := getRouterParam(request, "cluster")
	if err != nil {
		return "", err
	}

	validatedClusterName, err := ValidateClusterName(clusterName)
	if err != nil {
		return "", err
	}

	return validatedClusterName, nil
}

// readRequestID retrieves request ID from request
// if it's not possible, it writes http error to the writer and returns error
func readRequestID(_ http.ResponseWriter, request *http.Request) (types.RequestID, error) {
	requestID, err := getRouterParam(request, "request_id")
	if err != nil {
		return "", err
	}

	validatedRequestID, err := ValidateRequestID(requestID)
	if err != nil {
		return "", err
	}

	return validatedRequestID, nil
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

// serveContentWithGroups method implements the /content endpoint that also
// returns group info
func (server *HTTPServer) serveContentWithGroups(writer http.ResponseWriter, _ *http.Request) {
	log.Info().Msg("Content with groups handler")

	server.initGroupList()

	// prepare data structure
	responseData := map[string]interface{}{"status": "ok"}
	responseData["content"] = server.Content
	responseData["groups"] = server.groupsList

	err := responses.SendOK(writer, responseData)
	if err != nil {
		handleServerError(err)
		return
	}
}

/*
// serveContent method implements the /content endpoint
func (server *HTTPServer) serveContent(writer http.ResponseWriter, request *http.Request) {
	err := responses.SendOK(writer, responses.BuildOkResponseWithData("content", server.Content))
	if err != nil {
		handleServerError(err)
		return
	}

}
*/

func (server *HTTPServer) initGroupList() {
	// let's mimick Content Service behaviour preciselly
	if server.groupsList == nil {
		log.Info().Msg("Initializing group list for the first time")
		server.groupsList = make([]groups.Group, 0, len(server.Groups))

		for _, group := range server.Groups {
			server.groupsList = append(server.groupsList, group)
		}
	}

	log.Info().Int("items", len(server.groupsList)).Msg("Group list")
}

// listOfGroups returns the list of defined groups
func (server *HTTPServer) listOfGroups(writer http.ResponseWriter, _ *http.Request) {
	log.Info().Msg("List of groups handler")

	server.initGroupList()

	err := responses.SendOK(writer, responses.BuildOkResponseWithData("groups", server.groupsList))
	if err != nil {
		log.Error().Err(err).Msg("List of groups handler")
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
		writer.WriteHeader(http.StatusBadRequest)
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
		log.Info().Int("Code", code).Msg("Failed clusters")
		writer.WriteHeader(code)
		return
	}
	report, err := server.Storage.ReadReportForCluster(clusterName)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		log.Error().Err(err).Msg(unableToReadReportErrorMessage)
		handleServerError(err)
		return
	}
	writer.Header().Set(contentType, appJSON)

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

	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		log.Error().Err(err).Msg("dump error")
		return
	}

	log.Info().Str("dump", string(dump)).Msg("dump of request")

	err = json.NewDecoder(request.Body).Decode(&clusterList)

	if err != nil {
		log.Error().Err(err).Msg("getting list of clusters")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// server response has JSON format for this endpoint
	writer.Header().Set(contentType, appJSON)

	// construct reports for all clusters in a list
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

	// try to serialize all reports
	bytes, err := json.MarshalIndent(generatedReports, "", "\t")
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
		return
	}

	// send report back to client
	_, err = writer.Write(bytes)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

func (server *HTTPServer) readReportForOrganizationAndCluster(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(contentType, appJSON)

	organizationID, err := readOrganizationID(writer, request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	clusterName, err := readClusterName(writer, request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	report, err := server.Storage.ReadReportForOrganizationAndCluster(organizationID, clusterName)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
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
		log.Error().Err(err).Msg("parse rule selector")
		return types.Component(""), types.ErrorKey(""), err
	}

	IDValidator := regexp.MustCompile(`^[a-zA-Z_0-9.]+$`)

	isRuleIDValid := IDValidator.MatchString(splitedRuleID[0])
	isErrorKeyValid := IDValidator.MatchString(splitedRuleID[1])

	if !isRuleIDValid || !isErrorKeyValid {
		err := fmt.Errorf("invalid rule ID, each part of ID must contain only latin characters, number, underscores or dots")
		log.Error().Err(err).Msg("rule name validity check")
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

// ruleClusterDetailEndpoint method implements endpoint that should return a list of all the clusters IDs affected by this rule
func (server *HTTPServer) ruleClusterDetailEndpoint(writer http.ResponseWriter, request *http.Request) {
	// read the selector
	ruleSelector, err := readRuleSelector(writer, request)

	// check for missing/improper selector
	if err != nil {
		log.Error().Err(err).Msg("unable to read rule selector")
		err = responses.SendBadRequest(writer, err.Error())
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}
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
	writer.Header().Set(contentType, appJSON)
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

func logClusterName(clusterName types.ClusterName) {
	log.Info().Str("cluster_name", string(clusterName)).Msg(requestParameter)
}

func logRequestID(requestID types.RequestID) {
	log.Info().Str("request_id", string(requestID)).Msg(requestParameter)
}

// RequestStatus contains description about one request ID
type RequestStatus struct {
	RequestID string `json:"requestID"`
	Valid     bool   `json:"valid"`
	Received  string `json:"received"`
	Processed string `json:"processed"`
}

func constructRequestsList(requestIDs []types.RequestID) []RequestStatus {
	states := make([]RequestStatus, len(requestIDs))

	for i := range requestIDs {
		states[i].RequestID = string(requestIDs[i])
		states[i].Valid = true

		received := time.Date(2000, time.November, 1, 1, 0, 0, 999, time.UTC).Format(time.RFC3339Nano)
		states[i].Received = received

		processed := time.Now().UTC().Format(time.RFC3339Nano)
		states[i].Processed = processed
	}
	return states
}

// readListOfRequestIDs method implements endpoint that should return a list of
// all request IDs for given cluster
func (server *HTTPServer) readListOfRequestIDs(writer http.ResponseWriter, request *http.Request) {
	clusterName, err := readClusterName(writer, request)
	if err != nil {
		err = responses.SendBadRequest(writer, err.Error())
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}
		return
	}
	logClusterName(clusterName)

	requestIDs, found := data.RequestIDs[clusterName]
	if !found {
		err := responses.SendNotFound(writer, requestsForClusterNotFound)
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}
		return
	}

	// prepare data structure
	responseData := map[string]interface{}{"status": "ok"}
	responseData["cluster"] = string(clusterName)
	responseData["requests"] = constructRequestsList(requestIDs)

	err = responses.SendOK(writer, responseData)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

// RequestList represents sequence of request IDs
type RequestList []types.RequestID

func dumpRequest(request *http.Request) {
	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		log.Error().Err(err).Msg("dump error")
		return
	}
	log.Info().Str("dump", string(dump)).Msg("dump of request")
}

func readRequestList(writer http.ResponseWriter, request *http.Request) (RequestList, error) {
	var requestList RequestList
	err := json.NewDecoder(request.Body).Decode(&requestList)
	if err != nil {
		log.Error().Err(err).Msg("getting list of requests")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return requestList, nil
}

func logRequestIDs(message string, requestList RequestList) {
	log.Info().Msg(message)
	for _, requestID := range requestList {
		log.Info().Str("ID", string(requestID)).Msg("   ")
	}
}

func filterRequestIDs(requestIDs, requestList RequestList) RequestList {
	// filter request IDs using the rather inneficient way!
	var filteredIDs RequestList
	for _, requestID := range requestIDs {
		found := false
		for _, wantedID := range requestList {
			if requestID == wantedID {
				found = true
				break
			}
		}
		if found {
			filteredIDs = append(filteredIDs, requestID)
		}
	}
	return filteredIDs
}

// readListOfRequestIDsPostVariant method implements endpoint that should return a list of
// request IDs for given cluster
func (server *HTTPServer) readListOfRequestIDsPostVariant(writer http.ResponseWriter, request *http.Request) {
	dumpRequest(request)
	requestList, err := readRequestList(writer, request)
	if err != nil {
		return
	}
	logRequestIDs("List of requests send to service", requestList)

	clusterName, err := readClusterName(writer, request)
	if err != nil {
		err = responses.SendBadRequest(writer, err.Error())
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}
		return
	}
	logClusterName(clusterName)

	requestIDs, found := data.RequestIDs[clusterName]
	if !found {
		err := responses.SendNotFound(writer, requestsForClusterNotFound)
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}
		return
	}

	filteredIDs := filterRequestIDs(requestIDs, requestList)
	logRequestIDs("Filtered IDs", filteredIDs)

	// prepare data structure
	responseData := map[string]interface{}{"status": "ok"}
	responseData["cluster"] = string(clusterName)
	responseData["requests"] = constructRequestsList(filteredIDs)

	err = responses.SendOK(writer, responseData)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

// readStatusOfRequestID method implements endpoint that should return a status
// for given request ID. Currently the status is set to "processed" or
// "unknown" because we won't have information about "in-between" states.
func (server *HTTPServer) readStatusOfRequestID(writer http.ResponseWriter, request *http.Request) {
	clusterName, err := readClusterName(writer, request)
	if err != nil {
		err = responses.SendBadRequest(writer, err.Error())
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}
		return
	}
	logClusterName(clusterName)

	requestID, err := readRequestID(writer, request)
	if err != nil {
		err = responses.SendBadRequest(writer, err.Error())
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}
		return
	}
	logRequestID(requestID)

	requestIDs, found := data.RequestIDs[clusterName]
	if !found {
		err := responses.SendNotFound(writer, requestsForClusterNotFound)
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}
		return
	}

	// prepare data structure
	responseData := map[string]interface{}{}
	responseData["cluster"] = string(clusterName)
	responseData["requestID"] = requestID
	responseData["status"] = "unknown"

	// try to find the required request in requests IDs
	for _, storedRequestID := range requestIDs {
		if storedRequestID == requestID {
			// update data structure
			responseData["status"] = StatusProcessed
			break
		}
	}

	// send response back to user
	err = responses.SendOK(writer, responseData)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

// readRuleHitsForRequestID method implements endpoint that should return
// simplified result for given request ID
func (server *HTTPServer) readRuleHitsForRequestID(writer http.ResponseWriter, request *http.Request) {
	clusterName, err := readClusterName(writer, request)
	if err != nil {
		err = responses.SendBadRequest(writer, err.Error())
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}
		return
	}
	logClusterName(clusterName)

	requestID, err := readRequestID(writer, request)
	if err != nil {
		err = responses.SendBadRequest(writer, err.Error())
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}
		return
	}
	logRequestID(requestID)

	_, found := data.RequestIDs[clusterName]
	if !found {
		err := responses.SendNotFound(writer, requestsForClusterNotFound)
		if err != nil {
			log.Error().Err(err).Msg(responseDataError)
		}
		return
	}

	// prepare data structure
	var responseData types.SimplifiedReport
	responseData.Cluster = string(clusterName)
	responseData.RequestID = string(requestID)
	responseData.Status = StatusProcessed
	// can be nil
	responseData.RuleHits = data.SimplifiedRuleHits[clusterName][requestID]

	bytes, err := json.MarshalIndent(responseData, "", "\t")
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
		return
	}

	// send response back to user
	_, err = writer.Write(bytes)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

func (server *HTTPServer) exit(writer http.ResponseWriter, _ *http.Request) {
	err := responses.SendOK(writer, responses.BuildOkResponse())
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
	err = server.Stop(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Error stopping HTTP server")
	}
	os.Exit(0)
}

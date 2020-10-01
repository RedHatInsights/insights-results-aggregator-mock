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

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ClusterName represents name of cluster in format c8590f31-e97e-4b85-b506-c45ce1911a12
type ClusterName string

// ClusterReport represents cluster report
type ClusterReport string

const path = "../../../data"

const address = ":8080"

const contentType = "Content-Type"
const appJSON = "application/json; charset=utf-8"

// responseDataError is used as the error message when the responses functions return an error
const responseDataError = "Unexpected error during response data encoding"

const reportForClusterMessage = "report for cluster"

var reports map[string]string = make(map[string]string)

var clusters = []string{
	"34c3ecc5-624a-49a5-bab8-4fdc5e51a266",
	"34c3ecc5-624a-49a5-bab8-4fdc5e51a267",
	"34c3ecc5-624a-49a5-bab8-4fdc5e51a268",
	"34c3ecc5-624a-49a5-bab8-4fdc5e51a269",
	"34c3ecc5-624a-49a5-bab8-4fdc5e51a26a",
	"34c3ecc5-624a-49a5-bab8-4fdc5e51a26b",
	"34c3ecc5-624a-49a5-bab8-4fdc5e51a26c",
	"34c3ecc5-624a-49a5-bab8-4fdc5e51a26d",
	"34c3ecc5-624a-49a5-bab8-4fdc5e51a26e",
	"34c3ecc5-624a-49a5-bab8-4fdc5e51a26f",
	"74ae54aa-6577-4e80-85e7-697cb646ff37",
	"a7467445-8d6a-43cc-b82c-7007664bdf69",
	"ee7d2bf4-8933-4a3a-8634-3328fe806e08",
	"eeeeeeee-eeee-eeee-eeee-000000000001",
	"00000001-624a-49a5-bab8-4fdc5e51a266",
	"00000001-624a-49a5-bab8-4fdc5e51a267",
	"00000001-624a-49a5-bab8-4fdc5e51a268",
	"00000001-624a-49a5-bab8-4fdc5e51a269",
	"00000001-624a-49a5-bab8-4fdc5e51a26a",
	"00000001-624a-49a5-bab8-4fdc5e51a26b",
	"00000001-624a-49a5-bab8-4fdc5e51a26c",
	"00000001-624a-49a5-bab8-4fdc5e51a26d",
	"00000001-624a-49a5-bab8-4fdc5e51a26e",
	"00000001-624a-49a5-bab8-4fdc5e51a26f",
	"00000001-6577-4e80-85e7-697cb646ff37",
	"00000001-8933-4a3a-8634-3328fe806e08",
	"00000001-8d6a-43cc-b82c-7007664bdf69",
	"00000001-eeee-eeee-eeee-000000000001",
	"00000002-624a-49a5-bab8-4fdc5e51a266",
	"00000002-6577-4e80-85e7-697cb646ff37",
	"00000002-8933-4a3a-8634-3328fe806e08",
	"00000003-8933-4a3a-8634-3328fe806e08",
	"00000003-8d6a-43cc-b82c-7007664bdf69",
	"00000003-eeee-eeee-eeee-000000000001",
}

func readReport(path string, clusterName string) (string, error) {
	absPath, err := filepath.Abs(path + "/report_" + clusterName + ".json")
	if err != nil {
		return "", err
	}
	// disable "G304 (CWE-22): Potential file inclusion via variable"
	// #nosec G304
	report, err := ioutil.ReadFile(absPath)
	if err != nil {
		return "", err
	}
	return string(report), nil
}

func initStorage(path string) error {
	for _, cluster := range clusters {
		report, err := readReport(path, cluster)
		if err != nil {
			return err
		}
		log.Info().Str("cluster-1", cluster).Msg("Read cluster report")
		reports[cluster] = report
	}
	return nil
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("logger initialized")
}

func startHttpServer() error {
	log.Info().Msgf("Starting HTTP server at '%s'", address)
	router := mux.NewRouter().StrictSlash(true)

	log.Info().Msgf("Initializing HTTP server at '%s'", address)
	server := &http.Server{Addr: address, Handler: router}
	addEndpointsToRouter(router)
	log.Info().Msgf("Server has been initiliazed")

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Error().Err(err).Msg("Unable to start HTTP/S server")
		return err
	}

	log.Info().Msgf("Server has been started")
	return nil
}

func setDefaultContentType(w http.ResponseWriter) {
	w.Header().Set(contentType, appJSON)
}

func sendResponse(statusCode int, w http.ResponseWriter, data interface{}) error {
	setDefaultContentType(w)
	w.WriteHeader(statusCode)
	if status, ok := data.(string); ok {
		return json.NewEncoder(w).Encode(buildResponse(status))
	} else if rawData, ok := data.([]byte); ok {
		_, err := w.Write(rawData)
		return err
	}

	return json.NewEncoder(w).Encode(data)
}

func buildResponse(status string) map[string]interface{} {
	return map[string]interface{}{"status": status}
}

func buildOkResponse() map[string]interface{} {
	return map[string]interface{}{"status": "ok"}
}

func sendOKResponse(w http.ResponseWriter, data map[string]interface{}) error {
	return sendResponse(http.StatusOK, w, data)
}

// mainEndpoint will handle the requests for / endpoint
func mainEndpoint(writer http.ResponseWriter, _ *http.Request) {
	err := sendOKResponse(writer, buildOkResponse())
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

// readClusterName retrieves cluster name from request
// if it's not possible, it writes http error to the writer and returns error
func readClusterName(writer http.ResponseWriter, request *http.Request) (ClusterName, error) {
	clusterName, err := getRouterParam(request, "cluster")
	if err != nil {
		return "", err
	}

	return ClusterName(clusterName), nil
}

// readIndex retrieves cluster index from request
func readIndex(writer http.ResponseWriter, request *http.Request) (int, error) {
	index, err := getRouterParam(request, "n")
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(index)
	if err != nil {
		return 0, err
	}

	return i, nil
}

// getRouterParam retrieves parameter from URL like `/organization/{org_id}`
func getRouterParam(request *http.Request, paramName string) (string, error) {
	value, found := mux.Vars(request)[paramName]
	if !found {
		return "", errors.New("Missing param")
	}

	return value, nil
}

func nextCluster(clusterName ClusterName) (ClusterName, error) {
	for i, cluster := range clusters {
		if cluster == string(clusterName) {
			nextIndex := (i + 1) / len(cluster)
			return ClusterName(clusters[nextIndex]), nil
		}
	}
	return "", errors.New("Wrong cluster name '" + string(clusterName) + "'")
}

func firstReportEndpoint(writer http.ResponseWriter, request *http.Request) {
	clusterName := ClusterName(clusters[0])
	log.Info().Str("first cluster", string(clusterName)).Msg(reportForClusterMessage)

	reportForCluster(writer, clusterName)
}

func nextReportEndpoint(writer http.ResponseWriter, request *http.Request) {
	clusterName, err := readClusterName(writer, request)
	if err != nil {
		log.Error().Err(err).Msg("get cluster name operation")
		return
	}
	log.Info().Str("last cluster", string(clusterName)).Msg(reportForClusterMessage)

	nextClusterName, err := nextCluster(clusterName)
	if err != nil {
		log.Error().Err(err).Msg("get next cluster")
		return
	}
	log.Info().Str("next cluster", string(nextClusterName)).Msg(reportForClusterMessage)

	reportForCluster(writer, clusterName)
}

func nthReportEndpoint(writer http.ResponseWriter, request *http.Request) {
	index, err := readIndex(writer, request)
	if err != nil {
		log.Error().Err(err).Msg("get cluster index operation")
		return
	}
	log.Info().Int("index", index).Msg("index of cluster")

	index = index % len(clusters)
	clusterName := clusters[index]
	log.Info().Str("cluster with index", clusterName).Msg(reportForClusterMessage)

	reportForCluster(writer, ClusterName(clusterName))
}

func reportForCluster(writer http.ResponseWriter, clusterName ClusterName) {
	report, err := readReportForCluster(clusterName)
	if err != nil {
		log.Error().Err(err).Msg("reading report for selected cluster")
		return
	}

	r := []byte(report)
	_, err = writer.Write(r)
	if err != nil {
		log.Error().Err(err).Msg("can not write response data")
	}
}

func addEndpointsToRouter(router *mux.Router) {
	apiPrefix := "/"
	MainEndpoint := ""
	FirstReportEndpoint := "first-report/"
	ReportEndpoint := "report/{cluster}"
	NthReportEndpoint := "report-by-n/{n}"
	router.HandleFunc(apiPrefix+MainEndpoint, mainEndpoint).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+FirstReportEndpoint, firstReportEndpoint).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+ReportEndpoint, nextReportEndpoint).Methods(http.MethodGet)
	router.HandleFunc(apiPrefix+NthReportEndpoint, nthReportEndpoint).Methods(http.MethodGet)
}

func readReportForCluster(clusterName ClusterName) (ClusterReport, error) {
	report, ok := reports[string(clusterName)]
	if !ok {
		return "", errors.New("can not read report")
	}
	return ClusterReport(report), nil
}

func main() {
	err := initStorage(path)
	if err != nil {
		log.Error().Err(err).Msg("initStorage error")
	}

	err = startHttpServer()
	if err != nil {
		log.Error().Err(err).Msg("startHttpServer")
	}
}

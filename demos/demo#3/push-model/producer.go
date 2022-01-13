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
	"encoding/csv"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	clusterNameMessage = "clsn"
	apiURLBase         = "localhost:8080"
	resultEndpoint     = "/result"
	messagesToConsume  = 100000
)

// ClusterName represents name of cluster in format c8590f31-e97e-4b85-b506-c45ce1911a12
type ClusterName string

// ClusterReport represents cluster report
type ClusterReport string

const path = "../../../data"

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("push model producer: logger initialized")
}

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
	report, err := ioutil.ReadFile(absPath) // #nosec G304
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

func readReportForCluster(clusterName string) (ClusterReport, error) {
	report, ok := reports[clusterName]
	if !ok {
		return "", errors.New("can not read report")
	}
	return ClusterReport(report), nil
}

func reportForCluster(clusterName string) ([]byte, error) {
	report, err := readReportForCluster(clusterName)
	if err != nil {
		log.Error().Err(err).Msg("reading report for selected cluster")
		return nil, err
	}

	r := []byte(report)
	return r, nil
}

func prepareConnection(apiURL string, endpoint string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: apiURL, Path: endpoint}
	log.Info().Str("address", u.String()).Msg("connecting to")

	connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func sendMessages(connection *websocket.Conn, messagesToConsume int, writer *csv.Writer) error {
	for i := 1; i < messagesToConsume; i++ {
		startTime := time.Now()
		clusterName := clusters[i%len(clusters)]
		report, err := reportForCluster(clusterName)
		if err != nil {
			log.Error().Err(err).Str(clusterNameMessage, clusterName).Msg("read report for cluster")
			return err
		}
		err = connection.WriteMessage(websocket.TextMessage, report)
		if err != nil {
			log.Error().Err(err).Msg("write message")
			return err
		}
		log.Info().Int("message", i).Str(clusterNameMessage, clusterName).Msg("message sent")

		duration := time.Since(startTime)
		usec := int(duration / time.Microsecond)
		log.Info().Int("usec", usec).Msg("duration for processing")
		err = writer.Write([]string{strconv.Itoa(i + 1), strconv.Itoa(usec), "0"})
		if err != nil {
			log.Error().Err(err).Msg("can not write record into CSV")
		}
	}
	log.Info().Int("messages", messagesToConsume).Msg("all messages sent")
	return nil
}

func closeConnection(connection *websocket.Conn) error {
	// cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	err := connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return err
}

func main() {
	err := initStorage(path)
	if err != nil {
		log.Error().Err(err).Msg("initStorage error")
	}

	connection, err := prepareConnection(apiURLBase, resultEndpoint)
	if err != nil {
		log.Error().Err(err).Msg("dial")
		return
	}

	log.Info().Msg("connection established")
	defer connection.Close()

	file, err := os.Create("results.csv")
	if err != nil {
		log.Error().Err(err).Msg("can not open file")
	}

	writer := csv.NewWriter(file)

	err = writer.Write([]string{"#", "usec", "error"})
	if err != nil {
		log.Error().Err(err).Msg("can not write table header into CSV")
	}

	// flush all the buffers
	writer.Flush()

	// close the file with CSV data
	err = file.Close()
	if err != nil {
		log.Error().Err(err).Msg("can not close file")
	}

	err = sendMessages(connection, messagesToConsume, writer)
	if err != nil {
		log.Error().Err(err).Msg("sendMessages")
		return
	}

	err = closeConnection(connection)
	if err != nil {
		log.Error().Err(err).Msg("closeConnection")
		return
	}

}

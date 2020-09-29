package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ClusterReport represents cluster report
type ClusterReport interface{}

// ClusterResponse represents response containing cluster
type ClusterResponse struct {
	Report ClusterReport `json:"report"`
	Status string        `json:"status"`
}

const (
	apiURL              = "http://localhost:8080/"
	contentTypeHeader   = "Content-Type"
	contentLengthHeader = "Content-Length"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("logger initialized")
}

func performRequest(client http.Client, baseurl string, n int) error {
	url := baseurl + strconv.Itoa(n)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("NewRequest")
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		log.Error().Err(err).Msg("DoRequest")
		return err
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msg("ReadResponseBody")
		return err
	}

	report := ClusterResponse{}
	err = json.Unmarshal(body, &report)
	if err != nil {
		log.Error().Err(err).Msg("JSON decoding")
		return err
	}

	log.Info().Int("len", len(body)).Msg("report length")
	return nil
}

func main() {
	url := apiURL + "report-by-n/"

	client := http.Client{
		Timeout: time.Second * 2,
	}

	file, err := os.Create("results.csv")
	if err != nil {
		log.Error().Err(err).Msg("can not open file")
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Error().Err(err).Msg("can not close file")
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"#", "usec", "error"})
	if err != nil {
		log.Error().Err(err).Msg("can not write table header into CSV")
	}

	for i := 0; i < 10000; i++ {
		startTime := time.Now()
		err := performRequest(client, url, i)
		if err != nil {
			err = writer.Write([]string{strconv.Itoa(i + 1), "0", "1"})
			if err != nil {
				log.Error().Err(err).Msg("can not write record into CSV")
			}
		}
		duration := time.Since(startTime)
		usec := int(duration / time.Microsecond)
		log.Info().Int("usec", usec).Msg("duration for processing")
		err = writer.Write([]string{strconv.Itoa(i + 1), strconv.Itoa(usec), "0"})
		if err != nil {
			log.Error().Err(err).Msg("can not write record into CSV")
		}
	}
}

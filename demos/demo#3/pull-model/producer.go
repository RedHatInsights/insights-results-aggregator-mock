package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const path = "../../../data"

var reports map[string]string = make(map[string]string)

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
	clusters := []string{
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
	for _, cluster := range clusters {
		report, err := readReport(path, cluster)
		if err != nil {
			return err
		}
		log.Info().Str("cluster", cluster).Msg("Read cluster report")
		reports[cluster] = report
	}
	return nil
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("logger initialized")
}

func main() {
	err := initStorage(path)
	if err != nil {
		log.Error().Err(err).Msg("initStorage error")
	}
}

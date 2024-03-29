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

// Package conf contains definition of data type named ConfigStruct that
// represents configuration of the mock service. This package also contains
// function named LoadConfiguration that can be used to load configuration from
// provided configuration file and/or from environment variables. Additionally
// two specific functions named GetServerConfiguration and
// GetGroupsConfiguration are to be used to return specific configuration
// options.
//
// Generated documentation is available at:
// https://godoc.org/github.com/RedHatInsights/insights-results-aggregator-mock/conf
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-results-aggregator-mock/packages/conf/configuration.html
package conf

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/RedHatInsights/insights-results-aggregator-mock/content"
	"github.com/RedHatInsights/insights-results-aggregator-mock/groups"
	"github.com/RedHatInsights/insights-results-aggregator-mock/server"
)

const (
	configFileEnvVariableName = "INSIGHTS_RESULTS_AGGREGATOR_MOCK_CONFIG_FILE"
)

// PathsConfiguration is data structure that represents path to directory
// containing files with mock data.
type PathsConfiguration struct {
	MockDataPath string `mapstructure:"mock_data" toml:"mock_data"`
}

// ConfigStruct is a structure holding the whole service configuration
type ConfigStruct struct {
	Server  server.Configuration  `mapstructure:"server" toml:"server"`
	Content content.Configuration `mapstructure:"content" toml:"content"`
	Groups  groups.Configuration  `mapstructure:"groups" toml:"groups"`
	Paths   PathsConfiguration    `mapstructure:"paths" toml:"paths"`
}

// Config has exactly the same structure as *.toml file
var Config ConfigStruct

// LoadConfiguration loads configuration from defaultConfigFile, file set in
// configFileEnvVariableName or from env
func LoadConfiguration(defaultConfigFile string) (ConfigStruct, error) {
	configFile, specified := os.LookupEnv(configFileEnvVariableName)
	if specified {
		// we need to separate the directory name and filename without
		// extension
		directory, basename := filepath.Split(configFile)
		file := strings.TrimSuffix(basename, filepath.Ext(basename))
		// parse the configuration
		viper.SetConfigName(file)
		viper.AddConfigPath(directory)
	} else {
		// parse the configuration
		viper.SetConfigName(defaultConfigFile)
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()
	if _, isNotFoundError := err.(viper.ConfigFileNotFoundError); !specified && isNotFoundError {
		// viper is not smart enough to understand the structure of
		// config by itself
		fakeTomlConfigWriter := new(bytes.Buffer)

		err := toml.NewEncoder(fakeTomlConfigWriter).Encode(Config)
		if err != nil {
			return Config, err
		}

		fakeTomlConfig := fakeTomlConfigWriter.String()

		viper.SetConfigType("toml")

		err = viper.ReadConfig(strings.NewReader(fakeTomlConfig))
		if err != nil {
			return Config, err
		}
	} else if err != nil {
		return Config, fmt.Errorf("fatal error config file: %s", err)
	}

	// override config from env if there's variable in env

	const envPrefix = "INSIGHTS_RESULTS_AGGREGATOR_MOCK_"

	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "__"))

	err = viper.Unmarshal(&Config)
	return Config, err
}

// GetServerConfiguration returns server configuration
func GetServerConfiguration() server.Configuration {
	err := checkIfFileExists(Config.Server.APISpecFile)
	if err != nil {
		log.Fatal().Err(err).Msg("All customer facing APIs MUST serve the current OpenAPI specification")
	}

	return Config.Server
}

// GetGroupsConfiguration returns groups configuration
func GetGroupsConfiguration() groups.Configuration {
	err := checkIfFileExists(Config.Groups.ConfigPath)
	if err != nil {
		log.Fatal().Err(err).Msg("The groups configuration file is not defined")
	}

	return Config.Groups
}

// GetContentConfiguration returns groups configuration
func GetContentConfiguration() content.Configuration {
	err := checkIfFileExists(Config.Content.Path)
	if err != nil {
		log.Fatal().Err(err).Msg("The content file is not defined")
	}

	return Config.Content
}

// checkIfFileExists returns nil if path doesn't exist or isn't a file,
// otherwise it returns corresponding error
func checkIfFileExists(path string) error {
	if path == "" {
		return fmt.Errorf("Empty path provided")
	}
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("The following file path does not exist. Path: '%v'", path)
	} else if err != nil {
		return err
	}

	if fileMode := fileInfo.Mode(); !fileMode.IsRegular() {
		return fmt.Errorf("The following file path is not a regular file. Path: '%v'", path)
	}

	return nil
}

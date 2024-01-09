/*
Copyright © 2020, 2021, 2022 Red Hat, Inc.

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

// Entry point to the insights content service
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/RedHatInsights/insights-results-aggregator-mock/conf"
	"github.com/RedHatInsights/insights-results-aggregator-mock/content"
	"github.com/RedHatInsights/insights-results-aggregator-mock/groups"
	"github.com/RedHatInsights/insights-results-aggregator-mock/server"
	"github.com/RedHatInsights/insights-results-aggregator-mock/storage"
)

const (
	// ExitStatusOK means that the tool finished with success
	ExitStatusOK = iota

	// ExitStatusServerError is returned in case of any REST API server-related error
	ExitStatusServerError

	// ExitStatusOther represents other errors that might happen
	ExitStatusOther

	defaultConfigFilename = "config"
)

var (
	serverInstance *server.HTTPServer

	// BuildVersion contains the major.minor version of the CLI client
	BuildVersion = "*not set*"

	// BuildTime contains timestamp when the CLI client has been built
	BuildTime = "*not set*"

	// BuildBranch contains Git branch used to build this application
	BuildBranch = "*not set*"

	// BuildCommit contains Git commit used to build this application
	BuildCommit = "*not set*"
)

// startService starts service and returns error code
func startService(config *conf.ConfigStruct) int {
	serverCfg := conf.GetServerConfiguration()
	groupsCfg := conf.GetGroupsConfiguration()
	contentCfg := conf.GetContentConfiguration()

	ruleGroups, err := groups.ParseGroupConfigFile(groupsCfg.ConfigPath)
	if err != nil {
		log.Error().Err(err).Msg("Groups init error")
		return ExitStatusServerError
	}

	ruleContent, err := content.ParseContent(contentCfg.Path)
	if err != nil {
		log.Error().Err(err).Msg("Content init error")
		return ExitStatusServerError
	}
	log.Info().Int("count", len(ruleContent)).Msg("Content read")

	storageInstance, err := storage.New(config.Paths.MockDataPath)
	if err != nil {
		log.Error().Err(err).Msg("Storage construction error")
		return ExitStatusServerError
	}
	err = storageInstance.Init()
	if err != nil {
		log.Error().Err(err).Msg("Storage initialization error")
		return ExitStatusServerError
	}

	serverInstance = server.New(serverCfg, storageInstance, ruleGroups, ruleContent)

	err = serverInstance.Start()
	if err != nil {
		log.Error().Err(err).Msg("HTTP(s) start error")
		return ExitStatusServerError
	}

	return ExitStatusOK
}

func printInfo(msg, val string) {
	fmt.Printf("%s\t%s\n", msg, val)
}

func printVersionInfo() int {
	printInfo("Version:", BuildVersion)
	printInfo("Build time:", BuildTime)
	printInfo("Branch:", BuildBranch)
	printInfo("Commit:", BuildCommit)
	return ExitStatusOK
}

func initInfoLog(msg string) {
	log.Info().Str("type", "init").Msg(msg)
}

func logVersionInfo() {
	initInfoLog("Version: " + BuildVersion)
	initInfoLog("Build time: " + BuildTime)
	initInfoLog("Branch: " + BuildBranch)
	initInfoLog("Commit: " + BuildCommit)
}

const helpMessageTemplate = `
Service to provide content for OCP rules

Usage:

    %+v [command]

The commands are:

    <EMPTY>                      starts content service
    start-service                starts content service
    help     print-help          prints help
    config   print-config        prints current configuration set by files & env variables
    version  print-version-info  prints version info
    authors  print-authors       prints authors

`

const authorsList = `
Authors:
Pavel Tisnovsky <ptisnovs@redhat.com>
`

func printHelp() int {
	fmt.Printf(helpMessageTemplate, os.Args[0])
	return ExitStatusOK
}

func printAuthors() int {
	fmt.Print(authorsList)

	return ExitStatusOK
}

func printConfig(config *conf.ConfigStruct) int {
	configBytes, err := json.MarshalIndent(config, "", "    ")

	if err != nil {
		log.Error().Err(err).Msg("print config")
		return ExitStatusOther
	}

	fmt.Println(string(configBytes))

	return ExitStatusOK
}

func main() {
	config, err := conf.LoadConfiguration(defaultConfigFilename)
	if err != nil {
		panic(err)
	}

	command := "start-service"

	if len(os.Args) >= 2 {
		command = strings.ToLower(strings.TrimSpace(os.Args[1]))
	}

	os.Exit(handleCommand(&config, removeDashes(command)))
}

// function removeDashes removes one or two dashes from the beginning of a
// given string.
func removeDashes(command string) string {
	command = strings.TrimPrefix(command, "--")
	command = strings.TrimPrefix(command, "-")
	return command
}

func handleCommand(config *conf.ConfigStruct, command string) int {
	switch command {
	case "start-service":
		logVersionInfo()

		errCode := startService(config)
		if errCode != ExitStatusOK {
			return errCode
		}
		return ExitStatusOK
	case "help", "print-help":
		return printHelp()
	case "config", "print-config":
		return printConfig(&conf.Config)
	case "version", "print-version-info":
		return printVersionInfo()
	case "authors", "print-authors":
		return printAuthors()
	default:
		fmt.Printf("\nCommand '%v' not found\n", command)
		return printHelp()
	}
}

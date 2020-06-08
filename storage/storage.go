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
*/

package storage

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/RedHatInsights/insights-results-aggregator-mock/types"
)

// Storage represents an interface to almost any database or storage system
type Storage interface {
	Init() error
	Close() error
	ListOfOrgs() ([]types.OrgID, error)
	ListOfClustersForOrg(orgID types.OrgID) ([]types.ClusterName, error)
	ReadReportForCluster(orgID types.OrgID, clusterName types.ClusterName) (types.ClusterReport, error)
	ReadReportForClusterByClusterName(clusterName types.ClusterName) (types.ClusterReport, types.Timestamp, error)
	ReportsCount() (int, error)
	VoteOnRule(
		clusterID types.ClusterName,
		ruleID types.RuleID,
		userID types.UserID,
		userVote types.UserVote,
	) error
	AddOrUpdateFeedbackOnRule(
		clusterID types.ClusterName,
		ruleID types.RuleID,
		userID types.UserID,
		message string,
	) error
	GetUserFeedbackOnRule(
		clusterID types.ClusterName, ruleID types.RuleID, userID types.UserID,
	) (*UserFeedbackOnRule, error)
	GetContentForRules(
		rules types.ReportRules,
		userID types.UserID,
		clusterName types.ClusterName,
	) ([]types.RuleContentResponse, error)
	ToggleRuleForCluster(
		clusterID types.ClusterName,
		ruleID types.RuleID,
		userID types.UserID,
		ruleToggle RuleToggle,
	) error
	ListDisabledRulesForCluster(
		clusterID types.ClusterName,
		userID types.UserID,
	) ([]types.DisabledRuleResponse, error)
	GetFromClusterRuleToggle(
		types.ClusterName,
		types.RuleID,
		types.UserID,
	) (*ClusterRuleToggle, error)
	DeleteFromRuleClusterToggle(
		clusterID types.ClusterName,
		ruleID types.RuleID,
		userID types.UserID,
	) error
	GetRuleByID(ruleID types.RuleID) (*types.Rule, error)
	GetOrgIDByClusterID(cluster types.ClusterName) (types.OrgID, error)
	GetUserFeedbackOnRules(
		clusterID types.ClusterName,
		rulesContent []types.RuleContentResponse,
		userID types.UserID,
	) (map[types.RuleID]types.UserVote, error)
	GetRuleWithContent(ruleID types.RuleID, ruleErrorKey types.ErrorKey) (*types.RuleWithContent, error)
}

type MemoryStorage struct {
}

var reports map[string]string = make(map[string]string)

func readReport(path string, clusterName string) (string, error) {
	absPath, err := filepath.Abs(path + "/report_" + clusterName + ".json")
	if err != nil {
		return "", err
	}
	report, err := ioutil.ReadFile(absPath)
	if err != nil {
		return "", err
	}
	return string(report), nil
}

func initStorage(path string) error {
	clusters := []string{
		"34c3ecc5-624a-49a5-bab8-4fdc5e51a266",
		"74ae54aa-6577-4e80-85e7-697cb646ff37",
		"a7467445-8d6a-43cc-b82c-7007664bdf69",
		"ee7d2bf4-8933-4a3a-8634-3328fe806e08",
	}
	for _, cluster := range clusters {
		report, err := readReport(path, cluster)
		if err != nil {
			return err
		}
		reports[cluster] = report
	}
	return nil
}

// New function creates and initializes a new instance of Storage interface
func New(path string) (*MemoryStorage, error) {
	err := initStorage(path)
	return &MemoryStorage{}, err
}

// Init performs all database initialization
// tasks necessary for further service operation.
func (storage MemoryStorage) Init() error {
	log.Info().Msg("Initializing connection to data storage")
	return nil
}

// Close method closes the connection to database. Needs to be called at the end of application lifecycle.
func (storage MemoryStorage) Close() error {
	log.Info().Msg("Closing connection to data storage")
	return nil
}

// Report represents one (latest) cluster report.
//     Org: organization ID
//     Name: cluster GUID in the following format:
//         c8590f31-e97e-4b85-b506-c45ce1911a12
type Report struct {
	Org        types.OrgID         `json:"org"`
	Name       types.ClusterName   `json:"cluster"`
	Report     types.ClusterReport `json:"report"`
	ReportedAt types.Timestamp     `json:"reported_at"`
}

// ListOfOrgs reads list of all organizations that have at least one cluster report
func (storage MemoryStorage) ListOfOrgs() ([]types.OrgID, error) {
	orgs := []types.OrgID{
		11789772,
		11940171,
	}
	return orgs, nil
}

// ListOfClustersForOrg reads list of all clusters fro given organization
func (storage MemoryStorage) ListOfClustersForOrg(orgID types.OrgID) ([]types.ClusterName, error) {
	clusters := make([]types.ClusterName, 0)
	switch orgID {
	case 11940171:
		return clusters, errors.New("You have no permissions to get or change info about this organization")
	case 11789772:
		clusters = append(clusters, "34c3ecc5-624a-49a5-bab8-4fdc5e51a266")
		clusters = append(clusters, "74ae54aa-6577-4e80-85e7-697cb646ff37")
		clusters = append(clusters, "a7467445-8d6a-43cc-b82c-7007664bdf69")
		clusters = append(clusters, "ee7d2bf4-8933-4a3a-8634-3328fe806e08")
	}

	return clusters, nil
}

// GetOrgIDByClusterID reads OrgID for specified cluster
func (storage MemoryStorage) GetOrgIDByClusterID(cluster types.ClusterName) (types.OrgID, error) {
	var orgID uint64 = 42

	return types.OrgID(orgID), nil
}

func getReportForCluster(clusterName types.ClusterName) string {
	report, ok := reports[string(clusterName)]
	if !ok {
		return ""
	}
	return report
}

// ReadReportForCluster reads result (health status) for selected cluster for given organization
func (storage MemoryStorage) ReadReportForCluster(
	orgID types.OrgID, clusterName types.ClusterName,
) (types.ClusterReport, error) {
	var report string

	switch orgID {
	case 11940171:
		return types.ClusterReport(report), errors.New("You have no permissions to get or change info about this organization")
	case 11789772:
		report = getReportForCluster(clusterName)
	}

	return types.ClusterReport(report), nil
}

// ReadReportForClusterByClusterName reads result (health status) for selected cluster for given organization
func (storage MemoryStorage) ReadReportForClusterByClusterName(
	clusterName types.ClusterName,
) (types.ClusterReport, types.Timestamp, error) {
	var report string
	var lastChecked time.Time

	return types.ClusterReport(report), types.Timestamp(lastChecked.UTC().Format(time.RFC3339)), nil
}

// GetContentForRules retrieves content for rules that were hit in the report
func (storage MemoryStorage) GetContentForRules(
	reportRules types.ReportRules,
	userID types.UserID,
	clusterName types.ClusterName,
) ([]types.RuleContentResponse, error) {
	rules := make([]types.RuleContentResponse, 0)

	return rules, nil
}

// ReportsCount reads number of all records stored in database
func (storage MemoryStorage) ReportsCount() (int, error) {
	count := -1

	return count, nil
}

// GetRuleByID gets a rule by ID
func (storage MemoryStorage) GetRuleByID(ruleID types.RuleID) (*types.Rule, error) {
	var rule types.Rule

	return &rule, nil
}

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

// Package types contains declaration of various data types (usually structures)
// used elsewhere in the aggregator code.
//
// Generated documentation is available at:
// https://godoc.org/github.com/RedHatInsights/insights-results-aggregator-mock/types
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-results-aggregator-mock/packages/types/types.html
package types

import "time"

// OrgID represents organization ID
type OrgID uint32

// ClusterName represents name of cluster in format c8590f31-e97e-4b85-b506-c45ce1911a12
type ClusterName string

// ClusterReport represents cluster report
type ClusterReport string

// Timestamp represents any timestamp in a form gathered from database
// TODO: need to be improved
type Timestamp string

// UserVote is a type for user's vote
type UserVote int

// RequestID is used to store the request ID supplied in input Kafka records as
// a unique identifier of payloads. Empty string represents a missing request ID.
type RequestID string

// RuleOnReport represents a single (hit) rule of the string encoded report
type RuleOnReport struct {
	Module   string      `json:"component"`
	ErrorKey string      `json:"key"`
	Details  interface{} `json:"details"`
}

// ReportRules is a helper struct for easy JSON unmarshalling of string encoded report
type ReportRules struct {
	HitRules     []RuleOnReport `json:"reports"`
	SkippedRules []RuleOnReport `json:"skips"`
	PassedRules  []RuleOnReport `json:"pass"`
	TotalCount   int
}

// ReportResponse represents the response of /report endpoint
type ReportResponse struct {
	Meta  ReportResponseMeta    `json:"meta"`
	Rules []RuleContentResponse `json:"data"`
}

// ReportResponseMeta contains metadata about the report
type ReportResponseMeta struct {
	Count         int       `json:"count"`
	LastCheckedAt Timestamp `json:"last_checked_at"`
}

// RuleContentResponse represents a single rule in the response of /report endpoint
type RuleContentResponse struct {
	CreatedAt    string      `json:"created_at"`
	Description  string      `json:"description"`
	ErrorKey     string      `json:"-"`
	Generic      string      `json:"details"`
	Reason       string      `json:"reason"`
	Resolution   string      `json:"resolution"`
	TotalRisk    int         `json:"total_risk"`
	RiskOfChange int         `json:"risk_of_change"`
	RuleModule   string      `json:"rule_id"`
	TemplateData interface{} `json:"extra_data"`
	Tags         []string    `json:"tags"`
	UserVote     UserVote    `json:"user_vote"`
	Disabled     bool        `json:"disabled"`
}

// DisabledRuleResponse represents a single disabled rule displaying only identifying information
type DisabledRuleResponse struct {
	RuleModule  string `json:"rule_id"`
	Description string `json:"description"`
	Generic     string `json:"details"`
	DisabledAt  string `json:"disabled_at"`
}

// RuleID represents type for rule id
type RuleID string

// RuleSelector represents component + error key
type RuleSelector string

// Component represent name of component (of rule)
type Component string

// ErrorKey represents type for error key
type ErrorKey string

// UserID represents type for user id
type UserID string

// Rule represents the content of rule table
type Rule struct {
	Module     RuleID `json:"module"`
	Name       string `json:"name"`
	Summary    string `json:"summary"`
	Reason     string `json:"reason"`
	Resolution string `json:"resolution"`
	MoreInfo   string `json:"more_info"`
}

// RuleErrorKey represents the content of rule_error_key table
type RuleErrorKey struct {
	ErrorKey    ErrorKey  `json:"error_key"`
	RuleModule  RuleID    `json:"rule_module"`
	Condition   string    `json:"condition"`
	Description string    `json:"description"`
	Impact      int       `json:"impact"`
	Likelihood  int       `json:"likelihood"`
	PublishDate time.Time `json:"publish_date"`
	Active      bool      `json:"active"`
	Generic     string    `json:"generic"`
	Tags        []string  `json:"tags"`
}

// RuleWithContent represents a rule with content, basically the mix of rule and rule_error_key tables' content
type RuleWithContent struct {
	Module      RuleID    `json:"module"`
	Name        string    `json:"name"`
	Summary     string    `json:"summary"`
	Reason      string    `json:"reason"`
	Resolution  string    `json:"resolution"`
	MoreInfo    string    `json:"more_info"`
	ErrorKey    ErrorKey  `json:"error_key"`
	Condition   string    `json:"condition"`
	Description string    `json:"description"`
	TotalRisk   int       `json:"total_risk"`
	PublishDate time.Time `json:"publish_date"`
	Active      bool      `json:"active"`
	Generic     string    `json:"generic"`
	Tags        []string  `json:"tags"`
}

// RuleHit represents one rule hit for one defined cluster
type RuleHit struct {
	Component Component
	ErrorKey  ErrorKey
	Cluster   ClusterName
}

// KafkaOffset type for kafka offset
type KafkaOffset int64

// DBDriver type for db driver enum
type DBDriver int

// Acknowledge represents user acknowledgement of given rule
type Acknowledge struct {
	Acknowledged  bool   `json:"-"` // let's skip this one in responses
	Rule          string `json:"rule"`
	Justification string `json:"justification"`
	CreatedBy     string `json:"created_by"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// AcknowledgementsMetadata contains metadata about list of acknowledgements
type AcknowledgementsMetadata struct {
	Count int `json:"count"`
}

// AcknowledgementsResponse is structure returned to client in JSON
// serialization format
type AcknowledgementsResponse struct {
	Metadata AcknowledgementsMetadata `json:"meta"`
	Data     []Acknowledge            `json:"data"`
}

// AcknowledgementJustification data structure represents body of request with
// specified justification of given acknowledgement
type AcknowledgementJustification struct {
	Value string `json:"justification"`
}

// AcknowledgementRuleSelectorJustification data structure represents body of
// request with specified rule selector and justification of given
// acknowledgement
type AcknowledgementRuleSelectorJustification struct {
	RuleSelector RuleSelector `json:"rule_id"`
	Value        string       `json:"justification"`
}

// Alert data structure representing a single alert
type Alert struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Severity  string `json:"severity"`
	URL       string `json:"url"`
}

// OperatorCondition data structure representing a single operator condition
type OperatorCondition struct {
	Name      string `json:"name"`
	Condition string `json:"condition"`
	Reason    string `json:"reason"`
	URL       string `json:"url"`
}

// UpgradeRisksPredictors data structure represents the alerts and conditions
// that are contained in an upgrade-risk-predictions response
type UpgradeRisksPredictors struct {
	Alerts             []Alert             `json:"alerts"`
	OperatorConditions []OperatorCondition `json:"operator_conditions"`
}

// UpgradeRiskPrediction data structure represents body of the response
// for an upgrade-risk-predictions request
type UpgradeRiskPrediction struct {
	Recommended bool                   `json:"upgrade_recommended"`
	Predictors  UpgradeRisksPredictors `json:"upgrade_risks_predictors"`
}

const (
	// DBDriverSQLite3 shows that db driver is sqlite
	DBDriverSQLite3 DBDriver = iota
	// DBDriverPostgres shows that db driver is postgres
	DBDriverPostgres
	// DBDriverGeneral general sql(used for mock now)
	DBDriverGeneral
)

const (
	// UserVoteDislike shows user's dislike
	UserVoteDislike UserVote = -1
	// UserVoteNone shows no vote from user
	UserVoteNone UserVote = 0
	// UserVoteLike shows user's like
	UserVoteLike UserVote = 1
)

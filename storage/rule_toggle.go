// Copyright 2020 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"time"

	"github.com/RedHatInsights/insights-results-aggregator-mock/types"
)

// RuleToggle is a type for user's vote
type RuleToggle int

const (
	// RuleToggleDisable indicates the rule has been disabled
	RuleToggleDisable RuleToggle = 1
	// RuleToggleEnable indicates the rule has been (re)enabled
	RuleToggleEnable RuleToggle = 0
)

// ClusterRuleToggle represents a record from rule_cluster_toggle
type ClusterRuleToggle struct {
	ClusterID  types.ClusterName
	RuleID     types.RuleID
	UserID     types.UserID
	Disabled   RuleToggle
	DisabledAt time.Time
	EnabledAt  time.Time
	UpdatedAt  time.Time
}

// ToggleRuleForCluster toggles rule for specified cluster
func (storage MemoryStorage) ToggleRuleForCluster(
	_ types.ClusterName, _ types.RuleID, _ types.UserID, _ RuleToggle,
) error {

	return nil
}

// ListDisabledRulesForCluster retrieves disabled rules for specified cluster
func (storage MemoryStorage) ListDisabledRulesForCluster(
	_ types.ClusterName, _ types.UserID,
) ([]types.DisabledRuleResponse, error) {

	rules := make([]types.DisabledRuleResponse, 0)

	return rules, nil
}

// GetFromClusterRuleToggle gets a rule from cluster_rule_toggle
func (storage MemoryStorage) GetFromClusterRuleToggle(
	_ types.ClusterName, _ types.RuleID, _ types.UserID,
) (*ClusterRuleToggle, error) {
	var disabledRule ClusterRuleToggle

	return &disabledRule, nil
}

// DeleteFromRuleClusterToggle deletes a record from the table rule_cluster_toggle. Only exposed in debug mode.
func (storage MemoryStorage) DeleteFromRuleClusterToggle(
	_ types.ClusterName, _ types.RuleID, _ types.UserID,
) error {
	return nil
}

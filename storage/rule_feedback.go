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

// UserFeedbackOnRule shows user's feedback on rule
type UserFeedbackOnRule struct {
	ClusterID types.ClusterName
	RuleID    types.RuleID
	UserID    types.UserID
	Message   string
	UserVote  types.UserVote
	AddedAt   time.Time
	UpdatedAt time.Time
}

// VoteOnRule likes or dislikes rule for cluster by user. If entry exists, it overwrites it
func (storage MemoryStorage) VoteOnRule(
	clusterID types.ClusterName,
	ruleID types.RuleID,
	userID types.UserID,
	userVote types.UserVote,
) error {
	return nil
}

// AddOrUpdateFeedbackOnRule adds feedback on rule for cluster by user. If entry exists, it overwrites it
func (storage MemoryStorage) AddOrUpdateFeedbackOnRule(
	clusterID types.ClusterName,
	ruleID types.RuleID,
	userID types.UserID,
	message string,
) error {
	return nil
}

// addOrUpdateUserFeedbackOnRuleForCluster adds or updates feedback
// will update user vote and messagePtr if the pointers are not nil
func (storage MemoryStorage) addOrUpdateUserFeedbackOnRuleForCluster(
	clusterID types.ClusterName,
	ruleID types.RuleID,
	userID types.UserID,
	userVotePtr *types.UserVote,
	messagePtr *string,
) error {
	return nil
}

// GetUserFeedbackOnRule gets user feedback from DB
func (storage MemoryStorage) GetUserFeedbackOnRule(
	clusterID types.ClusterName, ruleID types.RuleID, userID types.UserID,
) (*UserFeedbackOnRule, error) {
	feedback := UserFeedbackOnRule{}

	return &feedback, nil
}

// GetUserFeedbackOnRules gets user feedbacks for defined array of rule IDs from DB
func (storage MemoryStorage) GetUserFeedbackOnRules(
	clusterID types.ClusterName, rulesContent []types.RuleContentResponse, userID types.UserID,
) (map[types.RuleID]types.UserVote, error) {
	feedbacks := make(map[types.RuleID]types.UserVote)

	return feedbacks, nil
}

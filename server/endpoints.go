/*
Copyright © 2020, 2021, 2022, 2023 Red Hat, Inc.

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

package server

const (
	// MainEndpoint defines suffix of the root endpoint
	MainEndpoint = ""

	// GroupsEndpoint defines suffix of the groups request endpoint
	GroupsEndpoint = "groups"

	// ContentEndpoint defines suffix of the content request endpoint
	ContentEndpoint = "content"

	// InfoEndpoint defines suffix for the endpoint to return services info
	InfoEndpoint = "info"

	// DeleteOrganizationsEndpoint deletes all {organizations}(comma separated array). DEBUG only
	DeleteOrganizationsEndpoint = "organizations/{organizations}"
	// DeleteClustersEndpoint deletes all {clusters}(comma separated array). DEBUG only
	DeleteClustersEndpoint = "clusters/{clusters}"
	// OrganizationsEndpoint returns all organizations
	OrganizationsEndpoint = "organizations"
	// ClustersEndpoint returns reports for selected clusters
	ClustersEndpoint = "clusters"
	// ClustersInOrgEndpoint returns reports for all clusters in selected organization
	ClustersInOrgEndpoint = "clusters/{organization}"
	// ReportEndpoint returns report for provided {organization} and {cluster}
	ReportEndpoint = "report/{organization}/{cluster}"
	// ReportForClusterEndpoint returns report for provided {cluster} (w/o organization)
	ReportForClusterEndpoint = "report/{cluster}"
	// ReportForClusterEndpoint2 returns report for provided {cluster} (w/o organization)
	ReportForClusterEndpoint2 = "clusters/{cluster}/report"
	// LikeRuleEndpoint likes rule with {rule_id} for {cluster} using current user(from auth header)
	LikeRuleEndpoint = "clusters/{cluster}/rules/{rule_id}/like"
	// DislikeRuleEndpoint dislikes rule with {rule_id} for {cluster} using current user(from auth header)
	DislikeRuleEndpoint = "clusters/{cluster}/rules/{rule_id}/dislike"
	// ResetVoteOnRuleEndpoint resets vote on rule with {rule_id} for {cluster} using current user(from auth header)
	ResetVoteOnRuleEndpoint = "clusters/{cluster}/rules/{rule_id}/reset_vote"
	// GetVoteOnRuleEndpoint is an endpoint to get vote on rule. DEBUG only
	GetVoteOnRuleEndpoint = "clusters/{cluster}/rules/{rule_id}/get_vote"
	// RuleEndpoint is an endpoint to create&delete a rule. DEBUG only
	RuleEndpoint = "rules/{rule_id}"
	// RuleErrorKeyEndpoint is for endpoints to create&delete a rule_error_key (DEBUG only)
	// and for endpoint to get a rule
	RuleErrorKeyEndpoint = "rules/{rule_id}/error_keys/{error_key}"
	// RuleGroupsEndpoint is a simple redirect endpoint to the insights-content-service API specified in configuration
	RuleGroupsEndpoint = "groups"
	// ClustersForOrganizationEndpoint returns all clusters for {organization}
	ClustersForOrganizationEndpoint = "organizations/{organization}/clusters"
	// DisableRuleForClusterEndpoint disables a rule for specified cluster
	DisableRuleForClusterEndpoint = "clusters/{cluster}/rules/{rule_id}/disable"
	// EnableRuleForClusterEndpoint re-enables a rule for specified cluster
	EnableRuleForClusterEndpoint = "clusters/{cluster}/rules/{rule_id}/enable"
	// RuleClusterDetailEndpoint should return a list of all the clusters IDs affected by this rule
	RuleClusterDetailEndpoint = "rule/{rule_selector}/clusters_detail/"

	// Endpoints to manipulate with simplified rule results stored
	// independently under "tracker_id" identifier

	// ListAllRequestIDs should return list of all request IDs detected for
	// given cluster. In reality the list is refreshing as old request IDs
	// are forgotten after 24 hours
	ListAllRequestIDs = "cluster/{cluster}/requests/"

	// StatusOfRequestID should return status of processing one given
	// request ID
	StatusOfRequestID = "cluster/{cluster}/request/{request_id}/status"

	// RuleHitsForRequestID should return simplified results for given
	// cluster and requestID
	RuleHitsForRequestID = "cluster/{cluster}/request/{request_id}/report"

	// Endpoints to acknowledge rule and to manipulate with
	// acknowledgements.

	// AckListEndpoint list acks from this account where the rule is
	// active. Will return an empty list if this account has no acks.
	AckListEndpoint = "ack"

	// AckAcknowledgePostEndpoint acknowledges (and therefore hides) a rule
	// from view in an account. If there's already an acknowledgement of
	// this rule by this account, then return that. Otherwise, a new ack is
	// created.
	AckAcknowledgePostEndpoint = "ack"

	// AckGetEndpoint acknowledges (and therefore hides) a rule from view
	// in an account. This view handles listing, retrieving, creating and
	// deleting acks. Acks are created and deleted by Insights rule ID, not
	// by their own ack ID.
	AckGetEndpoint = "ack/{rule_selector}"

	// AckUpdateEndpoint updates an acknowledgement for a rule, by rule ID.
	// A new justification can be supplied. The username is taken from the
	// authenticated request. The updated ack is returned.
	AckUpdateEndpoint = "ack/{rule_selector}"

	// AckDeleteEndpoint deletes an acknowledgement for a rule, by its rule
	// ID. If the ack existed, it is deleted and a 204 is returned.
	// Otherwise, a 404 is returned.
	AckDeleteEndpoint = "ack/{rule_selector}"

	// UpgradeRisksPredictionEndpoint returns the prediction about upgrading
	// the given cluster.
	UpgradeRisksPredictionEndpoint = "cluster/{cluster}/upgrade-risks-prediction"

	// MetricsEndpoint returns prometheus metrics
	MetricsEndpoint = "metrics"

	// ExitEndpoint perform server shutdown (in Debug mode only)
	ExitEndpoint = "exit"

	// AllDVONamespaces endpoint address.
	//
	// Returns the list of all DVO namespaces (i.e. array of objects) to
	// which this particular account has access.  Each object contains the
	// namespace ID, the namespace display name if available, the cluster
	// ID under which this namespace is created, and the number of
	// affecting recommendations for this namespace as well.
	//
	// BDD scenarios for this endpoint:
	// https://github.com/RedHatInsights/insights-behavioral-spec/blob/main/features/DVO_Recommendations/Smart_Proxy_REST_API.feature
	// - Accessing Smart Proxy REST API endpoint to retrieve list of all DVO namespaces for current organization
	// - Returning just active clusters in Smart Proxy REST API endpoint to retrieve list of all DVO namespaces for current organization
	// - Checking organization in Smart Proxy REST API endpoint to retrieve list of all DVO namespaces for current organization
	//
	AllDVONamespaces = "namespaces/dvo"

	// DVONamespaceForCluster1 endpoint address.
	//
	// Returns the list of all namespaces (i.e. array of objects) to which
	// this particular account has access filtered by {cluster_name}.  Each
	// object contains the namespace ID, the namespace display name if
	// available, the cluster ID under which this namespace is created
	// (repeated input), and the number of affecting recommendations for
	// this namespace as well.
	//
	// BDD scenarios for this endpoint:
	// https://github.com/RedHatInsights/insights-behavioral-spec/blob/main/features/DVO_Recommendations/Smart_Proxy_REST_API.feature
	// - Accessing Smart Proxy REST API endpoint to retrieve DVO namespaces for not known cluster
	// - Accessing Smart Proxy REST API endpoint to retrieve DVO namespaces for known cluster without DVO namespaces
	// - Accessing Smart Proxy REST API endpoint to retrieve DVO namespaces for improper cluster UUID
	// - Accessing Smart Proxy REST API endpoint to retrieve DVO namespaces for improper organization
	//
	DVONamespaceForCluster1 = "cluster/{cluster_name}/namespaces/dvo/{namespace}"

	// DVONamespaceForCluster2 endpoint address.
	//
	// This endpoint has the same meaning as DVONamespaceForCluster1, but
	// order of selectors is different.
	//
	// BDD scenarios for this endpoint:
	// https://github.com/RedHatInsights/insights-behavioral-spec/blob/main/features/DVO_Recommendations/Smart_Proxy_REST_API.feature
	// - Accessing Smart Proxy REST API endpoint to retrieve DVO namespaces for selected existing cluster
	// - Accessing Smart Proxy REST API endpoint to retrieve DVO namespaces for not known cluster
	// - Accessing Smart Proxy REST API endpoint to retrieve DVO namespaces for known cluster without DVO namespaces
	// - Accessing Smart Proxy REST API endpoint to retrieve DVO namespaces for improper cluster UUID
	// - Accessing Smart Proxy REST API endpoint to retrieve DVO namespaces for improper organization
	// - Accessing Smart Proxy REST API endpoint to retrieve DVO namespaces when cluster is not specified
	//
	DVONamespaceForCluster2 = "namespaces/dvo/{namespace}/cluster/{cluster_name}"

	// DVONamespaceInfo endpoint address.
	//
	// Returns information about the requested namespace. Contains the display name,
	// associated cluster ID. Probably, some other metadata like last seen (but not
	// needed according to current UX pre-design).
	//
	// BDD scenarios for this endpoint:
	// https://github.com/RedHatInsights/insights-behavioral-spec/blob/main/features/DVO_Recommendations/Smart_Proxy_REST_API.feature
	// - Accessing Smart Proxy REST API endpoint to retrieve information about selected DVO namespace
	// - Accessing Smart Proxy REST API endpoint to retrieve information about selected DVO namespace when no such namespace exists
	// - Accessing Smart Proxy REST API endpoint to retrieve information about selected DVO namespace for improper organization
	// - Accessing Smart Proxy REST API endpoint to retrieve information about selected DVO namespace for improper user
	//
	DVONamespaceInfo = "namespaces/dvo/{namespace_id}/info"

	// DVONamespaceReports endpoint address.
	//
	// Returns the list of all recommendations affecting this namespace. It
	// is basically an array with objects meeting the
	// https://github.com/RedHatInsights/insights-results-smart-proxy/blob/master/server/api/v2/openapi.json#L1537
	//
	// BDD scenarios for this endpoint:
	// https://github.com/RedHatInsights/insights-behavioral-spec/blob/main/features/DVO_Recommendations/Smart_Proxy_REST_API.feature
	// - Accessing Smart Proxy REST API endpoint to retrieve list of all recommendations affecting the selected namespace when no rules are hitting
	// - Accessing Smart Proxy REST API endpoint to retrieve list of all recommendations affecting the selected namespace when just one rule is hitting
	// - Accessing Smart Proxy REST API endpoint to retrieve list of all recommendations affecting the selected namespace when two rules are hitting
	// - Accessing Smart Proxy REST API endpoint to retrieve list of all recommendations affecting the selected namespace when no such namespace exists
	// - Accessing Smart Proxy REST API endpoint to retrieve list of all recommendations affecting the selected namespace for improper organization
	// - Accessing Smart Proxy REST API endpoint to retrieve list of all recommendations affecting the selected namespace for improper user
	//
	DVONamespaceReports = "namespaces/dvo/{namespace_id}/reports"
)

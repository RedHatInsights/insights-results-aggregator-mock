/*
Copyright © 2021 Red Hat, Inc.

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

// Already acked rules to be used by client code.
import "github.com/RedHatInsights/insights-results-aggregator-mock/types"

// I don't like too long string literals in code
const (
	rule1 = "ccx_rules_ocp.external.rules.nodes_requirements_check.report|NODES_MINIMUM_REQUIREMENTS_NOT_MET"
	rule2 = "ccx_rules_ocp.external.bug_rules.bug_1766907.report|BUGZILLA_BUG_1766907"
	rule3 = "ccx_rules_ocp.external.rules.nodes_kubelet_version_check.report|NODE_KUBELET_VERSION"
	rule4 = "ccx_rules_ocp.external.rules.samples_op_failed_image_import_check.report|SAMPLES_FAILED_IMAGE_IMPORT_ERR"
	rule5 = "ccx_rules_ocp.external.rules.cluster_wide_proxy_auth_check.report|AUTH_OPERATOR_PROXY_ERROR"
)

// the only data structure we have to deal with
var acks map[types.RuleSelector]types.Acknowledge = make(map[types.RuleSelector]types.Acknowledge)

// initialize the acks data structure
func init() {
	acks[rule1] = types.Acknowledge{
		Acknowledged:  true,
		Rule:          rule1,
		Justification: "Justification1",
		CreatedBy:     "tester1",
		CreatedAt:     "2021-09-04T17:11:35.130Z",
		UpdatedAt:     "2021-09-04T17:11:35.130Z"}

	acks[rule2] = types.Acknowledge{
		Acknowledged:  true,
		Rule:          rule2,
		Justification: "Justification2",
		CreatedBy:     "tester2",
		CreatedAt:     "2021-09-04T17:11:35.130Z",
		UpdatedAt:     "2021-09-04T17:11:35.130Z"}

	acks[rule3] = types.Acknowledge{
		Acknowledged:  true,
		Rule:          rule3,
		Justification: "Justification3",
		CreatedBy:     "tester3",
		CreatedAt:     "2021-09-04T17:11:35.130Z",
		UpdatedAt:     "2021-09-04T17:11:35.130Z"}

	acks[rule4] = types.Acknowledge{
		Acknowledged:  true,
		Rule:          rule4,
		Justification: "Justification4",
		CreatedBy:     "tester4",
		CreatedAt:     "2021-09-04T17:11:35.130Z",
		UpdatedAt:     "2021-09-04T17:11:35.130Z"}

	acks[rule5] = types.Acknowledge{
		Acknowledged:  true,
		Rule:          rule5,
		Justification: "Justification5",
		CreatedBy:     "tester5",
		CreatedAt:     "2021-09-04T17:11:35.130Z",
		UpdatedAt:     "2021-09-04T17:11:35.130Z"}
}

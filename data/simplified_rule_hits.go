/*
Copyright © 2023 Red Hat, Inc.

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

package data

import (
	"github.com/RedHatInsights/insights-results-aggregator-mock/types"
)

// SimplifiedRuleHits contains mapping between cluster name+request ID and sequence rule hits
var SimplifiedRuleHits map[types.ClusterName]map[types.RequestID][]types.SimplifiedRuleHit

// init is called before the program enters the main function, so it is perfect
// time to initialize maps etc.
func init() {
	SimplifiedRuleHits = make(map[types.ClusterName]map[types.RequestID][]types.SimplifiedRuleHit, 10)

	ruleHit1 := types.SimplifiedRuleHit{
		RuleFQDN:    "ccx_rules_ocp.external.rules.nodes_requirements_check.report",
		ErrorKey:    "NODES_MINIMUM_REQUIREMENTS_NOT_MET",
		Description: "Lorem ipsum...",
		TotalRisk:   1,
	}

	ruleHit2 := types.SimplifiedRuleHit{
		RuleFQDN:    "samples_op_failed_image_import_check.report",
		ErrorKey:    "SAMPLES_FAILED_IMAGE_IMPORT_ERR",
		Description: "Lorem ipsum...",
		TotalRisk:   2,
	}

	ruleHit3 := types.SimplifiedRuleHit{
		RuleFQDN:    "ccx_rules_ocp.external.bug_rules.bug_1766907.report",
		ErrorKey:    "BUGZILLA_BUG_1766907",
		Description: "Lorem ipsum...",
		TotalRisk:   3,
	}

	ruleHit4 := types.SimplifiedRuleHit{
		RuleFQDN:    "ccx_rules_ocp.external.rules.nodes_kubelet_version_check.report",
		ErrorKey:    "NODE_KUBELET_VERSION",
		Description: "Lorem ipsum...",
		TotalRisk:   4,
	}

	ruleHit5 := types.SimplifiedRuleHit{
		RuleFQDN:    "ccx_rules_ocp.external.rules.samples_op_failed_image_import_check.report",
		ErrorKey:    "SAMPLES_FAILED_IMAGE_IMPORT_ERR",
		Description: "Lorem ipsum...",
		TotalRisk:   5,
	}

	rm1 := map[types.RequestID][]types.SimplifiedRuleHit{
		"3nl2vda87ld6e3s25jlk7n2dna": {
			ruleHit1,
			ruleHit2,
		},
		"18njbjudvkc521w8buicx2clri": {
			ruleHit2,
			ruleHit3,
		},
		"38584huk209q82uhl8md5gsdxr": {
			ruleHit1,
			ruleHit2,
			ruleHit3,
			ruleHit4,
			ruleHit5,
		},
	}

	rm2 := map[types.RequestID][]types.SimplifiedRuleHit{
		"1zlcewj4kjtsp37x0yyr6cwhgr": {
			ruleHit1,
		},
		"3m3imli92shw225d4c3glzycxq": {
			ruleHit2,
		},
		"13yqlst6dmdji2z717w2v5fwcp": {
			ruleHit1,
			ruleHit2,
		},
		"271w1b53jlfjq2axaetgpe0yrd": {
			ruleHit2,
			ruleHit3,
		},
		"32zr43d2a4cbq1ogi1eu3hrti1": {
			ruleHit1,
			ruleHit2,
			ruleHit3,
			ruleHit4,
			ruleHit5,
		},
		"3pyjpvp4umqwx1xnhdq3mwgzkh": {
			ruleHit5,
		},
	}
	SimplifiedRuleHits["34c3ecc5-624a-49a5-bab8-4fdc5e51a267"] = rm1
	SimplifiedRuleHits["74ae54aa-6577-4e80-85e7-697cb646ff37"] = rm2
	SimplifiedRuleHits["eeeeeeee-eeee-eeee-eeee-000000000001"] = nil
}

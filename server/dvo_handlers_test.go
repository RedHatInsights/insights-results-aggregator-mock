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

package server_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/insights-results-aggregator-mock/server"
	"github.com/RedHatInsights/insights-results-aggregator-mock/types"
)

func getTestDVOWorkloads() []types.DVOWorkload {
	// just static data ATM are needed
	return []types.DVOWorkload{
		types.DVOWorkload{Rule: "host_network", CheckDescription: "Alert on pods/deployment-likes with sharing host's network namespace", CheckRemediation: "Ensure the host's network namespace is not shared.", Kind: "DaemonSet", NamespaceUID: "fbcbe2d3-e398-4b40-9d5e-4eb46fe8286f", UID: "be466de5-12fb-4710-bf70-62deb38ae563"},
		types.DVOWorkload{Rule: "host_network", CheckDescription: "Alert on pods/deployment-likes with sharing host's network namespace", CheckRemediation: "Ensure the host's network namespace is not shared.", Kind: "DaemonSet", NamespaceUID: "e6ed9bb3-efc3-46a6-b3ae-3f1a6e59546c", UID: "da5a07e1-3273-4056-8914-2732beb41b4c"},
		types.DVOWorkload{Rule: "host_pid", CheckDescription: "Alert on pods/deployment-likes with sharing host's process namespace", CheckRemediation: "Ensure the host's process namespace is not shared.", Kind: "DaemonSet", NamespaceUID: "e6ed9bb3-efc3-46a6-b3ae-3f1a6e59546c", UID: "da5a07e1-3273-4056-8914-2732beb41b4c"},
		types.DVOWorkload{Rule: "host_pid", CheckDescription: "Alert on pods/deployment-likes with sharing host's process namespace", CheckRemediation: "Ensure the host's process namespace is not shared.", Kind: "DaemonSet", NamespaceUID: "d00b47da-fc6f-4c72-abc1-94f525441c75", UID: "fec695db-b904-4865-b8e6-068f491c9a3b"},
		types.DVOWorkload{Rule: "non_isolated_pod", CheckDescription: "Alert on deployment-like objects that are not selected by any NetworkPolicy.", CheckRemediation: "Ensure pod does not accept unsafe traffic by isolating it with a NetworkPolicy. See https://cloud.redhat.com/blog/gUID:e-to-kubernetes-ingress-network-policies for more details.", Kind: "CronJob", NamespaceUID: "4354a80c-a7a6-405b-bfa6-9666b24e3b48", UID: "7b97edf7-8627-4f0e-a36f-822ccab0a0ae"},
		types.DVOWorkload{Rule: "non_isolated_pod", CheckDescription: "Alert on deployment-like objects that are not selected by any NetworkPolicy.", CheckRemediation: "Ensure pod does not accept unsafe traffic by isolating it with a NetworkPolicy. See https://cloud.redhat.com/blog/gUID:e-to-kubernetes-ingress-network-policies for more details.", Kind: "CronJob", NamespaceUID: "4354a80c-a7a6-405b-bfa6-9666b24e3b48", UID: "d641d1b5-a574-469e-82a1-3a4f985e2ddb"},
	}
}

func TestGetNamespaces(t *testing.T) {
	workloads := getTestDVOWorkloads()
	namespaces := server.GetNamespaces(workloads)

	// maintain expected order
	sort.Sort(sort.StringSlice(namespaces))

	expected := []string{
		"4354a80c-a7a6-405b-bfa6-9666b24e3b48",
		"d00b47da-fc6f-4c72-abc1-94f525441c75",
		"e6ed9bb3-efc3-46a6-b3ae-3f1a6e59546c",
		"fbcbe2d3-e398-4b40-9d5e-4eb46fe8286f"}

	assert.Equal(t, expected, namespaces)
}

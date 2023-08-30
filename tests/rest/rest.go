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

// Package tests contains REST API tests for following endpoints:
//
// apiPrefix
// apiPrefix + "clusters"
// apiPrefix + "groups"
// apiPrefix + "organizations"
package tests

// ServerTests run all tests for basic REST API endpoints
func ServerTests() {
	BasicTests()
}

// BasicTests implements basic tests for REST API apiPrefix
func BasicTests() {
	// tests for OpenAPI specification that is accessible via its endpoint as well
	// implementation of these tests is stored in openapi.go
	checkOpenAPISpecification()

	// implementations of these tests are stored in entrypoint.go
	checkRestAPIEntryPoint()
	checkNonExistentEntryPoint()
	checkWrongEntryPoint()
	checkWrongMethodsForEntryPoint()

	// implementations of these tests are stored in groups.go
	checkGroupsEndpoint()
	checkWrongMethodsForGroupsEndpoint()

	// implementations of these tests are stored in organizations.go
	checkOrganizationsEndpoint()
	checkWrongMethodsForOrganizationsEndpoint()

	// implementations of these tests are stored in clusters.go
	checkClustersEndpointForOrganization1()
	checkClustersEndpointForOrganization2()
	checkWrongMethodsForClustersEndpoint()

	// implementations of these tests are stored in reports.go
	checkReportForKnownOrganizationKnownCluster()
	checkReportForUnknownOrganization()
	checkReportForImproperOrganization()
	checkReportForKnownOrganizationUnknownCluster()
	checkReportForKnownOrganizationWrongCluster()
	checkWrongMethodsForReportForOrgAndClusterEndpoint()

	checkReportForKnownCluster()
	checkReportForUnknownCluster()
	checkReportForImproperCluster()
	checkReportForFailedCluster200()
	checkReportForFailedCluster400()
	checkReportForFailedCluster500()
	checkReportForFailedClusterNegativeTestCase()
	checkWrongMethodsForClusterReportEndpoint()

	checkReportsForAllClustersInOrganizationPositiveTestCase()

	// implementations of these tests is stored in multiple_reports.go
	checkMultipleReportsForKnownOrganizationAnd1KnownClusterUsingPostMethod()
	checkMultipleReportsForKnownOrganizationAnd2KnownClustersUsingPostMethod()
	checkMultipleReportsForKnownOrganizationAnd3KnownClustersUsingPostMethod()
	checkMultipleReportsForKnownOrganizationAndUnknownClusterUsingPostMethod()
	checkMultipleReportsForKnownOrganizationAndKnownAndUnknownClusterUsingPostMethod()
	checkMultipleReportsForKnownOrganizationAnd2KnownAndUnknownClusterUsingPostMethod()
	checkMultipleReportsForKnownOrganizationAnd3KnownAndUnknownClusterUsingPostMethod()

	// implementations of these tests are stored in content.go
	checkContentEndpoint()
	checkWrongMethodsForContentEndpoint()

	// implementations of these tests are stored in requests.go
	checkListAllRequestIDsForKnownCluster()
	checkListAllRequestIDsEmptyList()
	checkListAllRequestIDsForUnknownCluster()
	checkListSelectedRequestIDsForKnownCluster()
	checkListSelectedRequestIDsForUnknownCluster()
	checkListSelectedRequestIDsEmptyList()
	checkRetrieveRequestStatusForKnownClusterAndKnownRequest()
	checkRetrieveRequestStatusForKnownClusterAndUnknownRequest()
	checkRetrieveRequestStatusForUnknownCluster()
	checkRetrieveRequestReportForKnownClusterAndKnownRequest()
	checkRetrieveRequestReportForKnownClusterAndUnknownRequest()
	checkRetrieveRequestReportForUnknownCluster()

	// implementations of these tests are stored in rules.go
	checkRetrieveClusterDetailsForKnownRule()
	checkRetrieveClusterDetailsForUnknownRule()
	checkRetrieveClusterDetailsForImproperRule()

	// implementation of these tests are stored in acks.go
	checkRetrieveListOfAcks()
	checkAckRule()
	checkAckRuleWithIncorrectName()
	checkAckRuleViaGetEndpoint()
	checkAckRuleSeveralTimes()
	checkAckIncorrectRule()
	checkUpdateExistingRule()
	checkUpdateNonExistingRule()
	checkUpdateIncorrectRule()
	checkDeleteExistingRule()
	checkDeleteNonExistingRule()
	checkDeleteIncorrectRule()

	// implementation of these tests are stored in dvo.go
	checkListOfDVONamespaces()
	checkListOfDVONamespacesOtherMethods()
}

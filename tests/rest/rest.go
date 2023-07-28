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

// Package tests contains REST API tests for following endpoints:
//
// apiPrefix
package tests

// ServerTests run all tests for basic REST API endpoints
func ServerTests() {
	BasicTests()
}

// BasicTests implements basic tests for REST API apiPrefix
func BasicTests() {
	// implementation of these tests is stored in entrypoint.go
	checkRestAPIEntryPoint()
	checkNonExistentEntryPoint()
	checkWrongEntryPoint()
	checkWrongMethodsForEntryPoint()

	// implementation of these tests is stored in groups.go
	checkGroupsEndpoint()
	checkWrongMethodsForGroupsEndpoint()
}

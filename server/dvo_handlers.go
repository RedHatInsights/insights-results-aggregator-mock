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

package server

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// allDVONamespaces handler returns list of all DVO namespaces. Currently it
// does not depend on Organization ID as this information is passed through
// Bearer token in real Smart Proxy service. The format of output should be:
//
//           {
//             "status": "ok",
//             "workloads": [
//                 {
//                     "cluster": {
//                         "uuid": "{cluster UUID}",
//                         "display_name": "{cluster UUID or displayable name}",
//                     },
//                     "namespace": {
//                         "uuid": "{namespace UUID}",
//                         "name": "{namespace real name}", // optional, might be null
//                     },
//                     "reports": [
//                         {
//                             "check": "{for example no_anti_affinity}", // taken from the original full name deploment_validation_operator_no_anti_affinity
//                             "kind": "{kind attribute}",
//                             "description": {description}",
//                             "remediation": {remediation}",
//                         },
//                     ]
//             ]
//         }
func (server *HTTPServer) allDVONamespaces(writer http.ResponseWriter, request *http.Request) {
	log.Info().Msg("All DVO namespaces handler")

	// prepare response structure
	var responseData AllDVONamespacesResponse
	responseData.Status = "ok"

	// transform response structure into proper JSON payload
	bytes, err := json.MarshalIndent(responseData, "", "\t")
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
		return
	}

	// and send the response to client
	_, err = writer.Write(bytes)
	if err != nil {
		log.Error().Err(err).Msg(responseDataError)
	}
}

# insights-results-aggregator-mock
Mock service mimicking Insights Results Aggregator

[![forthebadge made-with-go](http://ForTheBadge.com/images/badges/made-with-go.svg)](https://go.dev/)

[![GoDoc](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator-mock?status.svg)](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator-mock)
[![GitHub Pages](https://img.shields.io/badge/%20-GitHub%20Pages-informational)](https://redhatinsights.github.io/insights-results-aggregator-mock/)
[![Go Report Card](https://goreportcard.com/badge/github.com/RedHatInsights/insights-results-aggregator-mock)](https://goreportcard.com/report/github.com/RedHatInsights/insights-results-aggregator-mock)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/RedHatInsights/insights-results-aggregator-mock)
[![License](https://img.shields.io/badge/license-Apache-blue)](https://github.com/RedHatInsights/insights-results-aggregator-mock/blob/master/LICENSE)


<!-- vim-markdown-toc GFM -->

* [Description](#description)
* [Howto build the service](#howto-build-the-service)
* [Howto start the service](#howto-start-the-service)
* [Generate the image for Docker](#generate-the-image-for-docker)
* [Running in Docker](#running-in-docker)
* [Makefile targets](#makefile-targets)
* [Usage](#usage)
* [Accessing results](#accessing-results)
    * [Settings for localhost](#settings-for-localhost)
    * [Basic endpoints](#basic-endpoints)
    * [Rule content](#rule-content)
    * [Groups](#groups)
    * [Clusters per organization](#clusters-per-organization)
    * [Report for organization + cluster](#report-for-organization--cluster)
    * [Report for one particular cluster](#report-for-one-particular-cluster)
    * [Getting report for several clusters](#getting-report-for-several-clusters)
* [List of cluster IDs that can be accesses by this service](#list-of-cluster-ids-that-can-be-accesses-by-this-service)
    * [Clusters that return 'static' rule results](#clusters-that-return-static-rule-results)
        * [Organization ID `11789772`](#organization-id-11789772)
        * [Organization ID `1`](#organization-id-1)
        * [Organization ID `2`](#organization-id-2)
        * [Organization ID `3`](#organization-id-3)
    * [Cluster that returns no results (ie just empty report)](#cluster-that-returns-no-results-ie-just-empty-report)
    * [Clusters that return rules that change every 15 minutes](#clusters-that-return-rules-that-change-every-15-minutes)
    * [List of clusters that return improper results and/or failure](#list-of-clusters-that-return-improper-results-andor-failure)
* [List of clusters hitting specified rule](#list-of-clusters-hitting-specified-rule)
    * [An example of response:](#an-example-of-response)
* [Endpoint to ack rule](#endpoint-to-ack-rule)
    * [List of acked rules](#list-of-acked-rules)
    * [Ack rule with specified justification](#ack-rule-with-specified-justification)
        * [Ack new rule](#ack-new-rule)
        * [Ack existing rule](#ack-existing-rule)
    * [Ack rule](#ack-rule)
        * [Ack new rule](#ack-new-rule-1)
        * [Ack existing rule](#ack-existing-rule-1)
    * [Update rule](#update-rule)
        * [Update existing rule](#update-existing-rule)
        * [Update non existing rule](#update-non-existing-rule)
    * [Delete rule](#delete-rule)
        * [Delete existing rule](#delete-existing-rule)
        * [Delete nonexisting rule](#delete-nonexisting-rule)
* [Package manifest](#package-manifest)

<!-- vim-markdown-toc -->


## Description

Mock service mimicking Insights Results Aggregator / SmartProxy REST API

## Howto build the service

```
make build
```

## Howto start the service

```
make run
```

## Generate the image for Docker

```
docker build -t insights-results-aggregator-mock:latest .
```

## Running in Docker

```
docker run --rm insights-results-aggregator-mock:latest
```

## Makefile targets

```
Usage: make <OPTIONS> ... <TARGETS>

Available targets are:

build                Build binary containing service executable
fmt                  Run go fmt -w for all sources
lint                 Run golint
vet                  Run go vet. Report likely mistakes in source code
cyclo                Run gocyclo
ineffassign          Run ineffassign checker
shellcheck           Run shellcheck
errcheck             Run errcheck
goconst              Run goconst checker
gosec                Run gosec checker
abcgo                Run ABC metrics checker
style                Run all the formatting related commands (fmt, vet, lint, cyclo) + check shell scripts
run                  Build the project and executes the binary
test                 Run the unit tests
cover                Generate HTML pages with code coverage
coverage             Display code coverage on terminal
integration_tests    Run all integration tests
help                 Show this help screen
function_list        List all functions in generated binary file
```

## Usage

```
Usage:

    ./insights-results-aggregator-mock [command]

The commands are:

    <EMPTY>                      starts content service
    start-service                starts content service
    help     print-help          prints help
    config   print-config        prints current configuration set by files & env variables
    version  print-version-info  prints version info
    authors  print-authors       prints authors
```

Note: it is possible to use single dash or double dashes for all commands.

## Accessing results

### Settings for localhost

```
ADDRESS=localhost:8080/api/insights-results-aggregator/v1
```

### Basic endpoints

```
curl -k -v $ADDRESS/
curl -k -v $ADDRESS/groups
curl -k -v $ADDRESS/content
curl -k -v $ADDRESS/organizations
curl -k -v $ADDRESS/clusters
```

### Rule content

Returns rule content and also group info:

```
curl -k -v $ADDRESS/content
```

An example of response (shortened):

```
{
  "content": [
    {
      "plugin": {
        "name": "",
        "node_id": "",
        "product_code": "",
        "python_module": "foo.bar.baz"
      },
  ],
  "groups": [
    {
      "title": "Performance",
      "description": "High utilization, proposed tuned profiles, storage issues",
      "tags": [
        "performance"
      ]
    },
  ],
  "status": "ok"
}
```

### Groups

```
curl -k -v $ADDRESS/groups
```

An example of response (further formatted by `jq`):

```
{
  "groups": [
    {
      "title": "Security",
      "description": "Issues related to certificates, user management, security groups, specific port usage, storage permissions, usage of kubeadmin account, exposed keys etc.",
      "tags": [
        "security"
      ]
    },
    {
      "title": "Fault Tolerance",
      "description": "Load balancer issues, machine api and autoscaler issues, failover issues, nodes down, cluster api/cluster provider issues.",
      "tags": [
        "fault_tolerance"
      ]
    },
    {
      "title": "Performance",
      "description": "High utilization, proposed tuned profiles, storage issues",
      "tags": [
        "performance"
      ]
    },
    {
      "title": "Service Availability",
      "description": "Operator degraded, missing functionality due to misconfiguration or resource constraints.",
      "tags": [
        "service_availability"
      ]
    }
  ],
  "status": "ok"
}
```

### Clusters per organization

```
curl -k -v $ADDRESS/organizations/11789772/clusters
curl -k -v $ADDRESS/organizations/11940171/clusters
```

### Report for organization + cluster

```
curl -k -v $ADDRESS/report/11789772/34c3ecc5-624a-49a5-bab8-4fdc5e51a266
```

or alternatively

```
curl -k -v $ADDRESS/clusters/11789772/34c3ecc5-624a-49a5-bab8-4fdc5e51a266/report
```

In this case `11789772` is organization ID and `34c3ecc5-624a-49a5-bab8-4fdc5e51a266` is cluster ID

### Report for one particular cluster

```
curl -k -v $ADDRESS/report/34c3ecc5-624a-49a5-bab8-4fdc5e51a266
```

### Getting report for several clusters

List of clusters has to be provided in payload in JSON format:

```
curl -k -v $ADDRESS/clusters -d @cluster_list.json
```

Format of the payload:

```json
{
        "clusters" : [
                "34c3ecc5-624a-49a5-bab8-4fdc5e51a266",
                "74ae54aa-6577-4e80-85e7-697cb646ff37",
                "a7467445-8d6a-43cc-b82c-7007664bdf69",
                "ee7d2bf4-8933-4a3a-8634-3328fe806e08"
        ]
}
```

Format of response:

```json
{
        "clusters": [
                "34c3ecc5-624a-49a5-bab8-4fdc5e51a266",
                "74ae54aa-6577-4e80-85e7-697cb646ff37",
                "a7467445-8d6a-43cc-b82c-7007664bdf69",
                "ee7d2bf4-8933-4a3a-8634-3328fe806e08"
        ],
        "errors": null,
        "reports": {
                "34c3ecc5-624a-49a5-bab8-4fdc5e51a266": {
                        "report": {
                            // ...
                            // ...
                            // ...
                        }
                }
        },
        "generated_at": "2020-08-11T10:17:29Z"
}
```

Response format in case it is not possible to return result for some cluster:

```json
{
        "clusters": [
                "ee7d2bf4-8933-4a3a-8634-3328fe806e08"
        ],
        "errors": [
                "00000000-0000-0000-0000-000000000000"
        ],
        "reports": {
                "ee7d2bf4-8933-4a3a-8634-3328fe806e08": {
                        "report": {
                                "data": [
                                    // ...
                                    // ...
                                    // ...
                                ]
                        },
                        "status": "ok"
                }
        },
        "generated_at": "2020-08-11T10:17:29Z"
}
```

## List of cluster IDs that can be accesses by this service

### Clusters that return 'static' rule results

#### Organization ID `11789772`

```
34c3ecc5-624a-49a5-bab8-4fdc5e51a266
34c3ecc5-624a-49a5-bab8-4fdc5e51a267
34c3ecc5-624a-49a5-bab8-4fdc5e51a268
34c3ecc5-624a-49a5-bab8-4fdc5e51a269
34c3ecc5-624a-49a5-bab8-4fdc5e51a26a
34c3ecc5-624a-49a5-bab8-4fdc5e51a26b
34c3ecc5-624a-49a5-bab8-4fdc5e51a26c
34c3ecc5-624a-49a5-bab8-4fdc5e51a26d
34c3ecc5-624a-49a5-bab8-4fdc5e51a26e
34c3ecc5-624a-49a5-bab8-4fdc5e51a26f
74ae54aa-6577-4e80-85e7-697cb646ff37
a7467445-8d6a-43cc-b82c-7007664bdf69
ee7d2bf4-8933-4a3a-8634-3328fe806e08
```

#### Organization ID `1`

```
00000001-624a-49a5-bab8-4fdc5e51a266
00000001-624a-49a5-bab8-4fdc5e51a267
00000001-624a-49a5-bab8-4fdc5e51a268
00000001-624a-49a5-bab8-4fdc5e51a269
00000001-624a-49a5-bab8-4fdc5e51a26a
00000001-624a-49a5-bab8-4fdc5e51a26b
00000001-624a-49a5-bab8-4fdc5e51a26c
00000001-624a-49a5-bab8-4fdc5e51a26d
00000001-624a-49a5-bab8-4fdc5e51a26e
00000001-624a-49a5-bab8-4fdc5e51a26f
00000001-6577-4e80-85e7-697cb646ff37
00000001-8933-4a3a-8634-3328fe806e08
00000001-8d6a-43cc-b82c-7007664bdf69
```

#### Organization ID `2`

```
00000002-624a-49a5-bab8-4fdc5e51a266
00000002-6577-4e80-85e7-697cb646ff37
00000002-8933-4a3a-8634-3328fe806e08
```

#### Organization ID `3`

```
00000003-8933-4a3a-8634-3328fe806e08
00000003-8d6a-43cc-b82c-7007664bdf69
```

### Cluster that returns no results (ie just empty report)

```
eeeeeeee-eeee-eeee-eeee-000000000001
00000001-eeee-eeee-eeee-000000000001
00000003-eeee-eeee-eeee-000000000001
```

**Mnemotechnic**: `e` means "empty"

### Clusters that return rules that change every 15 minutes

```
Cluster ID                            Returns results that are similar to:

cccccccc-cccc-cccc-cccc-000000000001  34c3ecc5-624a-49a5-bab8-4fdc5e51a266
                                      74ae54aa-6577-4e80-85e7-697cb646ff37
                                      a7467445-8d6a-43cc-b82c-7007664bdf69
cccccccc-cccc-cccc-cccc-000000000002  74ae54aa-6577-4e80-85e7-697cb646ff37
                                      a7467445-8d6a-43cc-b82c-7007664bdf69
                                      ee7d2bf4-8933-4a3a-8634-3328fe806e08
cccccccc-cccc-cccc-cccc-000000000003  ee7d2bf4-8933-4a3a-8634-3328fe806e08
                                      ee7d2bf4-8933-4a3a-8634-3328fe806e08
                                      34c3ecc5-624a-49a5-bab8-4fdc5e51a266
cccccccc-cccc-cccc-cccc-000000000004  eeeeeeee-eeee-eeee-eeee-000000000001
                                      eeeeeeee-eeee-eeee-eeee-000000000001
                                      34c3ecc5-624a-49a5-bab8-4fdc5e51a266
```

**Mnemotechnic**: `c` means "changing"

### List of clusters that return improper results and/or failure

```
ffffffff-ffff-ffff-ffff-000000000xxx'
```

Returns HTTP code xxx taken directly from the last three digits of cluster ID.
It means that devels/testers could use this functionality to check the
behaviour on client side.

**Mnemotechnic**: `f` means "failure"

Example:

```
ADDRESS=localhost:8080/api/insights-results-aggregator/v1

clusters="ffffffff-ffff-ffff-ffff-000000000200
ffffffff-ffff-ffff-ffff-000000000201
ffffffff-ffff-ffff-ffff-000000000404
ffffffff-ffff-ffff-ffff-000000000405
ffffffff-ffff-ffff-ffff-000000000201"

for cluster in $clusters
do
    curl -k -v $ADDRESS/report/${cluster}
done
```

## List of clusters hitting specified rule

```
curl 'localhost:8080/api/insights-results-aggregator/v1/rule/ccx_rules_ocp.external.rules.nodes_requirements_check.report|NODES_MINIMUM_REQUIREMENTS_NOT_MET/clusters_detail/'
```

### An example of response:

```
{
        "meta": {
                "count": 24,
                "component": "ccx_rules_ocp.external.rules.nodes_requirements_check.report",
                "error_key": "NODES_MINIMUM_REQUIREMENTS_NOT_MET",
                "generated_at": "2021-08-27T12:12:18Z"
        },
        "data": [
                "00000001-624a-49a5-bab8-4fdc5e51a266",
                "00000001-6577-4e80-85e7-697cb646ff37",
                "00000001-8933-4a3a-8634-3328fe806e08",
                "00000001-8d6a-43cc-b82c-7007664bdf69",
                "00000001-0000-0000-0000-000000000000",
                "00000001-1111-1111-1111-000000000000",
                "00000001-2222-2222-2222-000000000000",
                "00000001-3333-3333-3333-000000000000",
                "00000001-4444-4444-4444-000000000000",
                "00000001-5555-5555-5555-000000000000",
                "00000001-6666-6666-6666-000000000000",
                "00000001-7777-7777-7777-000000000000",
                "00000001-8888-8888-8888-000000000000",
                "00000001-9999-9999-9999-000000000000",
                "00000001-aaaa-aaaa-aaaa-000000000000",
                "00000001-bbbb-bbbb-bbbb-000000000000",
                "00000001-cccc-cccc-cccc-000000000000",
                "00000001-dddd-dddd-dddd-000000000000",
                "00000001-ffff-ffff-ffff-000000000000",
                "00000001-ffff-ffff-ffff-000000000000",
                "34c3ecc5-624a-49a5-bab8-4fdc5e51a266",
                "74ae54aa-6577-4e80-85e7-697cb646ff37",
                "a7467445-8d6a-43cc-b82c-7007664bdf69",
                "ee7d2bf4-8933-4a3a-8634-3328fe806e08"
        ]
}
```

## Endpoint to ack rule



### List of acked rules

List acks from this account where the rule is active. Will return an empty list if this account has no acks.

Request to the service:

```
curl localhost:8080/api/insights-results-aggregator/v1/ack
```

Response from the service:

```
{
        "meta": {
                "count": 5
        },
        "data": [
                {
                        "rule": "ccx_rules_ocp.external.rules.nodes_kubelet_version_check.report|NODE_KUBELET_VERSION",
                        "justification": "Justification3",
                        "created_by": "tester3",
                        "created_at": "2021-09-04T17:11:35.130Z",
                        "updated_at": "2021-09-04T17:11:35.130Z"
                },
                {
                        "rule": "ccx_rules_ocp.external.rules.samples_op_failed_image_import_check.report|SAMPLES_FAILED_IMAGE_IMPORT_ERR",
                        "justification": "Justification4",
                        "created_by": "tester4",
                        "created_at": "2021-09-04T17:11:35.130Z",
                        "updated_at": "2021-09-04T17:11:35.130Z"
                },
                {
                        "rule": "ccx_rules_ocp.external.rules.cluster_wide_proxy_auth_check.report|AUTH_OPERATOR_PROXY_ERROR",
                        "justification": "Justification5",
                        "created_by": "tester5",
                        "created_at": "2021-09-04T17:11:35.130Z",
                        "updated_at": "2021-09-04T17:11:35.130Z"
                },
                {
                        "rule": "ccx_rules_ocp.external.rules.nodes_requirements_check.report|NODES_MINIMUM_REQUIREMENTS_NOT_MET",
                        "justification": "Justification1",
                        "created_by": "tester1",
                        "created_at": "2021-09-04T17:11:35.130Z",
                        "updated_at": "2021-09-04T17:11:35.130Z"
                },
                {
                        "rule": "ccx_rules_ocp.external.bug_rules.bug_1766907.report|BUGZILLA_BUG_1766907",
                        "justification": "Justification2",
                        "created_by": "tester2",
                        "created_at": "2021-09-04T17:11:35.130Z",
                        "updated_at": "2021-09-04T17:11:35.130Z"
                }
        ]
}
```

### Ack rule with specified justification

Acknowledges (and therefore hides) a rule from view in an account. If there's
already an acknowledgement of this rule by this account, then return that.
Otherwise, a new ack is created.

#### Ack new rule

Request to the service:

```
curl -v -X POST -H "Content-Type: application/json" --data '{"rule_id":"foo|bar", "justification":"xyzzy"}' "localhost:8080/api/insights-results-aggregator/v1/ack"
```

Response from the service:

```
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Date: Sun, 05 Sep 2021 14:29:33 GMT
< Content-Length: 168
< 
{
        "rule": "foo|bar",
        "justification": "xyzzy",
        "created_by": "onlineTester",
        "created_at": "2021-09-05T16:29:33+02:00",
        "updated_at": "2021-09-05T16:29:33+02:00"
}
```

#### Ack existing rule

Request to the service:

```
curl -v -X POST -H "Content-Type: application/json" --data '{"rule_id":"existing|rule", "justification":"xyzzy"}' "localhost:8080/api/insights-results-aggregator/v1/ack"
```

Response from the service:

```
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Sun, 05 Sep 2021 14:35:51 GMT
< Content-Length: 168
< 
{
        "rule": "foo|bar",
        "justification": "xyzzy",
        "created_by": "onlineTester",
        "created_at": "2021-09-05T16:35:25+02:00",
        "updated_at": "2021-09-05T16:35:25+02:00"
}

```



### Ack rule

Acks acknowledge (and therefore hide) a rule from view in an account. This view
handles listing, retrieving, creating and deleting acks. Acks are created and
deleted by Insights rule ID, not by their own ack ID.

#### Ack new rule

Request to the service:

```
curl -v "localhost:8080/api/insights-results-aggregator/v1/ack/new|rule"
```

Response from the service:

```
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Sat, 04 Sep 2021 18:46:20 GMT
< Content-Length: 165
< 
{
        "rule": "new|rule",
        "justification": "?",
        "created_by": "onlineTester",
        "created_at": "2021-09-04T20:46:20+02:00",
        "updated_at": "2021-09-04T20:46:20+02:00"
}
```

#### Ack existing rule

Request to the service:

```
curl -v "localhost:8080/api/insights-results-aggregator/v1/ack/ccx_rules_ocp.external.rules.cluster_wide_proxy_auth_check.report|AUTH_OPERATOR_PROXY_ERROR"
```

Response from the service:

```
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Sat, 04 Sep 2021 18:32:58 GMT
< Content-Length: 260
< 
{
        "rule": "ccx_rules_ocp.external.rules.cluster_wide_proxy_auth_check.report|AUTH_OPERATOR_PROXY_ERROR",
        "justification": "Justification5",
        "created_by": "onlineTester",
        "created_at": "2021-09-04T17:11:35.130Z",
        "updated_at": "2021-09-04T20:32:58+02:00"

```

Please note that just `updated_at` attribute is changed in this situation.

### Update rule

#### Update existing rule

Request to the service:

```
curl -v-X PUT -H "Content-Type: application/json" --data '{"justification":"xyzzy"}' "localhost:8080/api/insights-results-aggregator/v1/ack/existing|rule"
```

Response from the service:

```
< HTTP/1.1 200 OK
< Date: Sun, 05 Sep 2021 05:47:46 GMT
< Content-Length: 169
< Content-Type: text/plain; charset=utf-8
< 
{
        "rule": "existing|rule",
        "justification": "xyzzy",
        "created_by": "onlineTester",
        "created_at": "2021-09-05T07:45:00+02:00",
        "updated_at": "2021-09-05T07:47:46+02:00"
}
```

#### Update non existing rule

Update an acknowledgement for a rule, by rule ID. A new justification can be
supplied. The username is taken from the authenticated request. The updated ack
is returned.

Request to the service:

```
curl -v -X PUT -H "Content-Type: application/json" --data '{"justification":"xyzzy"}' "localhost:8080/api/insights-results-aggregator/v1/ack/new|rule"
```

Response from the service:

```
< HTTP/1.1 404 Not Found
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Sun, 05 Sep 2021 06:13:27 GMT
< Content-Length: 51
< 
rule not found -> justification can not be changed
```

### Delete rule

Delete an acknowledgement for a rule, by its rule ID. If the ack existed, it is
deleted and a 204 is returned. Otherwise, a 404 is returned.

#### Delete existing rule

Request to the service:

```
curl-v -X DELETE "localhost:8080/api/insights-results-aggregator/v1/ack/ccx_rules_ocp.external.rules.cluster_wide_proxy_auth_check.report|AUTH_OPERATOR_PROXY_ERROR"
```

Response from the service:

```
< HTTP/1.1 204 No Content
< Date: Sat, 04 Sep 2021 17:44:32 GMT
< 
```

#### Delete nonexisting rule

Request to the service:

```
curl-v -X DELETE "localhost:8080/api/insights-results-aggregator/v1/ack/foobar|foobar"
```

Response from the service:

```
< HTTP/1.1 404 Not Found
< Date: Sat, 04 Sep 2021 17:44:35 GMT
< Content-Length: 0
< 
```

## Package manifest

Package manifest is available at [docs/manifest.txt](docs/manifest.txt).

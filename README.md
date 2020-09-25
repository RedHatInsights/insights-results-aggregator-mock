# insights-results-aggregator-mock
Mock service mimicking Insights Results Aggregator

[![GoDoc](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator-mock?status.svg)](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator-mock)
[![GitHub Pages](https://img.shields.io/badge/%20-GitHub%20Pages-informational)](https://redhatinsights.github.io/insights-results-aggregator-mock/)
[![Go Report Card](https://goreportcard.com/badge/github.com/RedHatInsights/insights-results-aggregator-mock)](https://goreportcard.com/report/github.com/RedHatInsights/insights-results-aggregator-mock)

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

## Accessing results

### Settings for localhost

```
ADDRESS=localhost:8080/api/v1
```

### Basic endpoints

```
curl -k -v $ADDRESS/
curl -k -v $ADDRESS/groups
curl -k -v $ADDRESS/organizations
curl -k -v $ADDRESS/clusters
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
ADDRESS=localhost:8080/api/v1

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

#!/usr/bin/env bash

ADDRESS=https://ira-mock-ccx-dev.apps.ocp.prod.psi.redhat.com/api/insights-results-aggregator

curl -k -v -u "insights-qa:redhatqa $ADDRESS/v1/" > main.json
curl -k -v -u "insights-qa:redhatqa $ADDRESS/v1/groups" > groups.json
curl -k -v -u "insights-qa:redhatqa $ADDRESS/v1/organizations" > organizations.json

curl -k -v -u "insights-qa:redhatqa $ADDRESS/v1/organizations/11789772/clusters" > clusters_11789772.json
curl -k -v -u "insights-qa:redhatqa $ADDRESS/v1/organizations/11940171/clusters" > clusters_11940171.json

clusters="34c3ecc5-624a-49a5-bab8-4fdc5e51a266 74ae54aa-6577-4e80-85e7-697cb646ff37 a7467445-8d6a-43cc-b82c-7007664bdf69 ee7d2bf4-8933-4a3a-8634-3328fe806e08"
clusters="34c3ecc5-624a-49a5-bab8-4fdc5e51a266 74ae54aa-6577-4e80-85e7-697cb646ff37 a7467445-8d6a-43cc-b82c-7007664bdf69 ee7d2bf4-8933-4a3a-8634-3328fe806e08 eeeeeeee-eeee-eeee-eeee-000000000001"

for cluster in $clusters
do
    curl -k -v -u insights-qa:redhatqa "$ADDRESS/v1/report/11789772/${cluster}" > "report_11789772_org_${cluster}.json"
done

for cluster in $clusters
do
    curl -k -v "$ADDRESS/report/${cluster}" > "localhost/report_${cluster}.json"
done

diff -r expected localhost

# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
    echo "Error!"
else
    echo "OK"
fi

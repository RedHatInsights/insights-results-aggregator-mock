#!/usr/bin/env bash

ADDRESS=localhost:8080//api/insights-results-aggregator/v2

mkdir -p localhost

curl -k -v $ADDRESS/ > localhost/main.json
curl -k -v $ADDRESS/groups > localhost/groups.json
curl -k -v $ADDRESS/organizations > localhost/organizations.json

curl -k -v $ADDRESS/organizations/11789772/clusters > localhost/clusters_11789772.json
curl -k -v $ADDRESS/organizations/11940171/clusters > localhost/clusters_11940171.json
curl -k -v $ADDRESS/organizations/1/clusters > localhost/clusters_00000001.json
curl -k -v $ADDRESS/organizations/2/clusters > localhost/clusters_00000002.json
curl -k -v $ADDRESS/organizations/3/clusters > localhost/clusters_00000003.json

curl -k -v $ADDRESS/clusters -d @cluster_list_1.json | jq --tab 'del(.generated_at)' > localhost/report_list1.json
curl -k -v $ADDRESS/clusters -d @cluster_list_2.json | jq --tab 'del(.generated_at)' > localhost/report_list2.json
curl -k -v $ADDRESS/clusters -d @cluster_list_3.json | jq --tab 'del(.generated_at)' > localhost/report_list3.json
curl -k -v $ADDRESS/clusters -d @cluster_list_4.json | jq --tab 'del(.generated_at)' > localhost/report_list4.json
curl -k -v $ADDRESS/clusters -d @cluster_list_5.json | jq --tab 'del(.generated_at)' > localhost/report_list5.json
curl -k -v $ADDRESS/clusters -d @cluster_list_6.json | jq --tab 'del(.generated_at)' > localhost/report_list6.json
curl -k -v $ADDRESS/clusters -d @cluster_list_7.json | jq --tab 'del(.generated_at)' > localhost/report_list7.json

clusters="34c3ecc5-624a-49a5-bab8-4fdc5e51a266 74ae54aa-6577-4e80-85e7-697cb646ff37 a7467445-8d6a-43cc-b82c-7007664bdf69 ee7d2bf4-8933-4a3a-8634-3328fe806e08 eeeeeeee-eeee-eeee-eeee-000000000001"

for cluster in $clusters
do
    curl -k -v "$ADDRESS/report/11789772/${cluster}" > "localhost/report_11789772_org_${cluster}.json"
done

for cluster in $clusters
do
    curl -k -v "$ADDRESS/report/${cluster}" > "localhost/report_${cluster}.json"
done

RED_BG=$(tput setab 1)
GREEN_BG=$(tput setab 2)
NC=$(tput sgr0) # No Color

diff -r expected localhost

# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
    echo "${RED_BG}Error!${NC}"
else
    echo "${GREEN_BG}OK${NC}"
fi

#!/usr/bin/env bash

ADDRESS=https://ira-mock-ccx-dev.apps.ocp.prod.psi.redhat.com/api/insights-results-aggregator

clusters="ffffffff-ffff-ffff-ffff-000000000200
ffffffff-ffff-ffff-ffff-000000000201
ffffffff-ffff-ffff-ffff-000000000404
ffffffff-ffff-ffff-ffff-000000000405
ffffffff-ffff-ffff-ffff-000000000201"

for cluster in $clusters
do
    curl -k -v -u insights-qa:redhatqa "$ADDRESS/v1/report/11789772/${cluster}" > "${cluster}.json"
done

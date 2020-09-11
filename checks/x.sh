#!/usr/bin/env bash

ADDRESS=localhost:8080/api/v1

clusters="ffffffff-ffff-ffff-ffff-000000000200
ffffffff-ffff-ffff-ffff-000000000201
ffffffff-ffff-ffff-ffff-000000000404
ffffffff-ffff-ffff-ffff-000000000405
ffffffff-ffff-ffff-ffff-000000000201"

for cluster in $clusters
do
    curl -k -v "$ADDRESS/report/${cluster}"
done

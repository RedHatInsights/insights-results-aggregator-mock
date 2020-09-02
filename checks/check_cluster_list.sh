#!/bin/bash

ADDRESS=localhost:8080/api/v1

mkdir -p localhost

curl -k -v $ADDRESS/clusters -d @cluster_list_1.json > localhost/report_list1.json
curl -k -v $ADDRESS/clusters -d @cluster_list_2.json > localhost/report_list2.json
curl -k -v $ADDRESS/clusters -d @cluster_list_3.json > localhost/report_list3.json
curl -k -v $ADDRESS/clusters -d @cluster_list_4.json > localhost/report_list4.json
curl -k -v $ADDRESS/clusters -d @cluster_list_5.json > localhost/report_list5.json
curl -k -v $ADDRESS/clusters -d @cluster_list_6.json > localhost/report_list6.json
curl -k -v $ADDRESS/clusters -d @cluster_list_7.json > localhost/report_list7.json

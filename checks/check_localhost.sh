ADDRESS=localhost:8080/api/v1

mkdir -p localhost

curl -k -v $ADDRESS/ > localhost/main.json
curl -k -v $ADDRESS/groups > localhost/groups.json
curl -k -v $ADDRESS/organizations > localhost/organizations.json

curl -k -v $ADDRESS/organizations/11789772/clusters > localhost/clusters_11789772.json
curl -k -v $ADDRESS/organizations/11940171/clusters > localhost/clusters_11940171.json

clusters="34c3ecc5-624a-49a5-bab8-4fdc5e51a266 74ae54aa-6577-4e80-85e7-697cb646ff37 a7467445-8d6a-43cc-b82c-7007664bdf69 ee7d2bf4-8933-4a3a-8634-3328fe806e08"

for cluster in $clusters
do
    curl -k -v $ADDRESS/report/11789772/${cluster} > localhost/report_${cluster}.json
done

diff -r expected localhost

if [ $? -ne 0 ]; then
    echo "Error!"
else
    echo "OK"
fi

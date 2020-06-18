(ns ACM-demo-#1)
(require '[clojure.data.json :as json])
(require '[clojure.pprint :as pprint])


; base address of mock service that run on PSI
(def base-address "https://ira-mock-ccx-dev.apps.ocp.prod.psi.redhat.com/api/insights-results-aggregator")


; prefix for all REST API endpoints
(def REST-API-prefix "/v1/")


; function to retrieve payload via REST API, parse it, and display
(defn retrieve-data
  [endpoint]
  (->
    (str base-address REST-API-prefix endpoint)
    slurp
    json/read-str
    pprint/pprint))


; function to print cluster report
(defn print-cluster-report
  [cluster]
  (retrieve-data (str "report/" cluster)))


; clusters with prepared reports
(def cluster1 "34c3ecc5-624a-49a5-bab8-4fdc5e51a266")
(def cluster2 "74ae54aa-6577-4e80-85e7-697cb646ff37")
(def cluster3 "a7467445-8d6a-43cc-b82c-7007664bdf69")
(def cluster4 "ee7d2bf4-8933-4a3a-8634-3328fe806e08")


; this is the shortest report
(print-cluster-report cluster4)


; let's try longer one
(print-cluster-report cluster1)


; done!
(println "Finito")


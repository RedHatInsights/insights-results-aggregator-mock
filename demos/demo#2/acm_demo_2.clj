(ns ACM-demo-#2)

(require '[clojure.pprint :as pprint])
(require '[clojure.data.json :as json])
(require '[clj-http.client :as client])

; base address of mock service that run on PSI or locally
(def base-address "http://localhost:8080")

; prefix for all REST API endpoints
(def REST-API-prefix "/api/v1/")

; function to retrieve results for selected clusters, parse it, and display
(defn get-cluster-results
  [filename]
  (let [payload (-> filename load-file json/write-str)]
    (->
      (str base-address REST-API-prefix "clusters")
      (client/post {:body payload})
      (:body)
      (println))))

; load original data, convert them to formatted JSON, and display them
(defn show-payload
  [filename]
  (let [payload (-> filename load-file)]
  (println "--------------------------------------")
  (json/pprint payload)
  (println "--------------------------------------")))



; show the payload that is to be send to REST API server
(show-payload "resources/cluster_list_1.clj")

; read and display results for four clusters
(get-cluster-results "resources/cluster_list_1.clj")



; three clusters
(show-payload "resources/cluster_list_2.clj")

(get-cluster-results "resources/cluster_list_2.clj")



; empty input list
(show-payload "resources/cluster_list_6.clj")

; response for empty input list
(get-cluster-results "resources/cluster_list_6.clj")



; one existing cluster and one non-existing one
(show-payload "resources/cluster_list_7.clj")

(get-cluster-results "resources/cluster_list_7.clj")



; one non-existing cluster
(show-payload "resources/cluster_list_8.clj")

(get-cluster-results "resources/cluster_list_8.clj")

; done!
(println "Finito")


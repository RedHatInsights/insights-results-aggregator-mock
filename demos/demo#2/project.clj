(defproject demo#2 "0.1.0-SNAPSHOT"
  :description "FIXME: write description"
  :url "http://example.com/FIXME"
  :license {:name "Eclipse Public License"
            :url "http://www.eclipse.org/legal/epl-v10.html"}
  :dependencies [[org.clojure/clojure "1.8.0"]
                 [org.clojure/data.json "1.0.0"]
                 [clj-http "3.10.1"]]
  :main ^:skip-aot demo#2.core
  :target-path "target/%s"
  :profiles {:uberjar {:aot :all}})

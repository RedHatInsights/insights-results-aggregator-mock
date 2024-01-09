;
;  (C) Copyright 2023  Pavel Tisnovsky
;
;  All rights reserved. This program and the accompanying materials
;  are made available under the terms of the Eclipse Public License v1.0
;  which accompanies this distribution, and is available at
;  http://www.eclipse.org/legal/epl-v10.html
;
;  Contributors:
;      Pavel Tisnovsky
;

(ns dvo-workloads.core
  "Core module containing the -main function and the startup code.

    Author: Pavel Tisnovsky"
  (:gen-class))


(require '[clojure.tools.logging   :as log])

(require '[ring.adapter.jetty      :as jetty])
(require '[ring.middleware.params  :as http-params])
(require '[ring.middleware.cookies :as cookies])
(require '[ring.middleware.session :as session])

(require '[clj-middleware.middleware  :as middleware])

(require '[dvo-workloads.server :as server])


(def app
  "Definition of a Ring-based application behaviour."
  (-> server/handler            ; handle all events
      ;(middleware/inject-configuration configuration) ; inject configuration
      ;                                                ; structure into the
      ;                                                ; parameter
      session/wrap-session
      cookies/wrap-cookies      ; we need to work with cookies
      http-params/wrap-params)) ; and to process request parameters, of course


(defn start-server
  "Start the HTTP server on the specified port.
     The port is specified as string."
  [port]
  (log/info "Starting the server at the port: " port)
  (jetty/run-jetty app {:port (read-string port)}))


(defn -main
  "Entry point to the service."
  [& args]
  (start-server "8080"))

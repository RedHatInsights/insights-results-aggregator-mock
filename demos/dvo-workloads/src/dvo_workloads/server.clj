(ns dvo-workloads.server
  "Server module with functions to accept requests and send response back to users via HTTP.")

(require '[ring.util.response      :as http-response])
(require '[clojure.tools.logging   :as log])
(require '[clojure.pprint          :as pprint])

(require '[clj-fileutils.fileutils :as fileutils])

(require '[clj-http-utils.http-utils :as http-utils])

(require '[dvo-workloads.html-renderer  :as html-renderer])

(use '[clj-utils.utils])


(defn finish-processing
  ([request response-html]
   ; use the previous session as the new one
   (finish-processing request response-html (:session request)))
  ([request response-html session]
   (let [cookies (:cookies request)
         user-id 42] ; not needed ATM
     (log/info "Incoming cookies: " cookies)
     (log/info "user-id:          " user-id)
     (-> (http-response/response response-html)
         (http-response/set-cookie "user-id" user-id {:max-age 36000000})
         ; use the explicitly specified new session with a map of values
         (assoc :session session)
         (http-response/content-type "text/html; charset=utf-8")))))


(defn process-index-page
  "Render index page."
  [request]
    (finish-processing
      request
      (html-renderer/render-index-page)))


(defn uri->file-name
  [uri]
  (subs uri (inc (.indexOf uri "/"))))


(defn gui-call-handler
  "This function is used to handle all GUI calls. Three parameters are expected:
   data structure containing HTTP request, string with URI, and the HTTP method."
  [request uri method]
  (cond (.endsWith uri ".gif")  (http-utils/return-file "www" (uri->file-name uri) "image/gif")
        (.endsWith uri ".png")  (http-utils/return-file "www" (uri->file-name uri) "image/png")
        (.endsWith uri ".jpg")  (http-utils/return-file "www" (uri->file-name uri) "image/jpeg")
        (.endsWith uri ".ico")  (http-utils/return-file "www" (uri->file-name uri) "image/x-icon")
        (.endsWith uri ".css")  (http-utils/return-file "www" (uri->file-name uri) "text/css")
        (.endsWith uri ".js")   (http-utils/return-file "www" (uri->file-name uri) "application/javascript")
        (.endsWith uri ".htm")  (http-utils/return-file "www" (uri->file-name uri) "text/html")
        (.endsWith uri ".html") (http-utils/return-file "www" (uri->file-name uri) "text/html")
        :else
      (condp = uri
          ; common pages
          "/"                           (process-index-page request)
          )))

(defn handler
  "Handler that is called by Ring for all requests received from user(s)."
  [request]
  (log/info "request URI:   " (:uri request))
  ;(log/info "configuration: " (:configuration request))
  (let [uri             (:uri request)
        method          (:request-method request)]
    ;(println uri)
    (gui-call-handler request uri method)))

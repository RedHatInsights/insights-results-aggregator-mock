(ns dvo-workloads.html-renderer-widgets
  "Module that contains common utility functions for the html-renderer and html-renderer help modules

    Author: Pavel Tisnovsky")


(require '[hiccup.core :as hiccup])
(require '[hiccup.page :as page])
(require '[hiccup.form :as form])


(defn header
  "Renders part of HTML page - the header."
  [url-prefix]
  [:head
   [:title "Insights Advisor Mock"]
   [:meta {:name "Author"    :content "Pavel Tisnovsky"}]
   [:meta {:name "Generator" :content "Clojure"}]
   [:meta {:http-equiv "Content-type" :content "text/html; charset=utf-8"}]
   (page/include-css (str url-prefix "bootstrap/bootstrap.min.css"))
   (page/include-css (str url-prefix "style.css"))
   (page/include-js (str url-prefix "dvo_workloads.js"))
  ] ; head
)


(defn navigation-bar
  "Renders whole navigation bar."
  [url-prefix]
  [:nav {:class "navbar navbar-inverse navbar-fixed-top" :role "navigation"} ; use navbar-default instead of navbar-inverse
      [:div {:class "container-fluid"}
          [:div {:class "row" :style "margin-bottom:8px;"}
              [:div {:class "col-md-7"}
                  [:div {:class "navbar-header"}
                      [:a {:href url-prefix :class "navbar-brand"} "Insights Advisor Mock"]
                  ] ; ./navbar-header
                  [:div {:class "navbar-header"}
                      [:ul {:class "nav navbar-nav"}
                          ;[:li (tab-class :whitelist mode) [:a {:href (str url-prefix "whitelist")} "Whitelist"]]
                      ]
                  ]
              ] ; col-md-7 ends
              ;[:div {:class "col-md-3"}
              ;    (render-name-field user-name (remember-me-href url-prefix mode))
              ;]
              ;[:div {:class "col-md-2"}
              ;    [:div {:class "navbar-header"}
              ;        [:a {:href (users-href url-prefix mode) :class "navbar-brand"} "Users"]
               ;   ] ; ./navbar-header
              ;] ; col ends
          ] ; row ends
      ] ; </div .container-fluid>
]); </nav>


(def footer
  "Renders part of HTML page - the footer."
  [:div "<br /><br />&copy; Pavel Tisnovsky, Red Hat"])


(def back-button
  "Back button widget."
  [:button {:class "btn btn-primary",
            :onclick "window.history.back()",
            :type "button",
            :style "width:12em"} "Back"])


(defn submit-button
  "Submit button widget."
  [text name value]
  [:button {:type "submit",
            :name name,
            :value value,
            :class "btn btn-success",
            :style "width:15em"} text])


(defn add-button
  "Add button widget."
  [language configuration drop-down-id]
  (let [onclick (str "onAddApplicationPart('" language "', '" configuration "', '" drop-down-id "')")]
    [:button {:type "button" :class "add_button" :style "width:7em" :onclick onclick} "Add"]))


(defn remove-button
  "Add button widget."
  [language configuration drop-down-id]
  (let [onclick (str "onRemoveApplicationPart('" language "', '" configuration "', '" drop-down-id "')")]
    [:button {:type "button" :class "remove_button" :style "width:7em" :onclick onclick} "Remove"]))


(defn disabled-submit-button
  "Disabled submit button widget."
  [text name value]
  [:button {:type "submit",
            :id name,
            :name name,
            :value value,
            :class "btn btn-success",
            :style "width:12em",
            :disabled "disabled"} text])


(defn radio-button
  "Radio button widget."
  ([group checked value label]
   [:span (form/radio-button group checked value)
    " " label
    [:br]])
  ([group checked value label on-click]
   [:span (form/radio-button {:onclick on-click} group checked value)
    " " label
    [:br]]))

(defn help-button
  "Help button widget."
  [help-page-url]
  [:a {:href help-page-url} [:img {:src "icons/help.gif"}]])


(def canvas
  [:div {:id "canvas_container"}])

(defn drop-down
  [drop-down-id drop-down-values]
  (form/drop-down {:id drop-down-id, :class "select"}
                  drop-down-id
                  drop-down-values))

(def short-vertical-separator
  [:div {:style "height: 2ex"}])

(def tall-vertical-separator
  [:div {:style "height: 10ex"}])

(defn warning-div
  [text]
  [:div {:class "alert alert-danger" :role "alert"} text])

(defn sidebar
  []
  [:div {:class "col-sm-2 navbar navbar-inverse" :style "height:1000px"}
   [:br]
   (form/form-to
     {:name "inputForm1"}
     [:get "/recommendations"]
     (submit-button "Recommendations" "recommendations" "recommendations")
     [:br][:br])
   (form/form-to
     {:name "inputForm2"}
     [:get "/clusters"]
     (submit-button "Clusters" "clusters" "clusters")
     [:br][:br])
   (form/form-to
     {:name "inputForm3"}
     [:get "/workloads"]
     (submit-button "Workloads" "workloads" "workloads")
     [:br][:br]) ])


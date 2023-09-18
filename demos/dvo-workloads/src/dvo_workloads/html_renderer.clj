(ns dvo-workloads.html-renderer
  "Module that contains functions used to render HTML pages sent back to the browser.

    Author: Pavel Tisnovsky")


(require '[clojure.string :as str])
(require '[clojure.pprint :as pprint])

(require '[hiccup.page :as page])
(require '[hiccup.form :as form])

(require '[dvo-workloads.html-renderer-widgets :as widgets])



(defn render-index-page
  "Render index page."
  []
  (page/xhtml
    (widgets/header "/")
    [:body
     [:div {:class "container"}
          (widgets/navigation-bar "/")
          [:h3 "Insights Advisor Mock"]
          [:div {:style "height: 10ex"}]
          (form/form-to
                {:name "inputForm1"}
                [:get "/select-app-type"]
                   (widgets/submit-button "Recommendations" "app-type" "app-type")
                [:br][:br])
          (form/form-to
                {:name "inputForm1"}
                [:get "/select-app-type"]
                   (widgets/submit-button "Clusters" "app-type" "app-type")
                [:br][:br])
          (form/form-to
                {:name "inputForm1"}
                [:get "/select-app-type"]
                   (widgets/submit-button "Workloads" "app-type" "app-type")
                [:br][:br])
          [:br]
          widgets/footer
     ] ; </div class="container">
    ] ; </body>
  ))

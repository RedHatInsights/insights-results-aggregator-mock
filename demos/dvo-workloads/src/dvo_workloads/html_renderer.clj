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
     [:div {:class "container-fluid"}
          (widgets/navigation-bar "/")
          [:h3 "Insights Advisor Mock"]
          [:div {:class "container-fluid" :style "padding:0px; margin:0px;"}
            [:div {:class "row justify-content-start"}
             (widgets/sidebar)
             [:div {:class "col-sm-2" :style "padding:20px"}
              "Text text"
             ]
            ]
          ]
          [:div {:style "height: 10ex"}]
          [:br]
          widgets/footer
     ] ; </div class="container">
    ] ; </body>
  ))


(defn render-recommendations-page
  []
  (page/xhtml
    (widgets/header "/")
    [:body
     [:div {:class "container-fluid"}
          (widgets/navigation-bar "/")
          [:h3 "Insights Advisor Mock"]
          [:div {:class "container-fluid" :style "padding:0px; margin:0px;"}
            [:div {:class "row justify-content-start"}
             (widgets/sidebar)
             [:div {:class "col-sm-2" :style "padding:20px"}
              "Text text"
             ]
            ]
          ]
          [:div {:style "height: 10ex"}]
          [:br]
          widgets/footer
     ] ; </div class="container">
    ] ; </body>
  ))


(defn render-clusters-page
  []
  (page/xhtml
    (widgets/header "/")
    [:body
     [:div {:class "container-fluid"}
          (widgets/navigation-bar "/")
          [:h3 "Insights Advisor Mock"]
          [:div {:class "container-fluid" :style "padding:0px; margin:0px;"}
            [:div {:class "row justify-content-start"}
             (widgets/sidebar)
             [:div {:class "col-sm-2" :style "padding:20px"}
              "Text text"
             ]
            ]
          ]
          [:div {:style "height: 10ex"}]
          [:br]
          widgets/footer
     ] ; </div class="container">
    ] ; </body>
  ))


(defn render-workloads-page
  []
  (page/xhtml
    (widgets/header "/")
    [:body
     [:div {:class "container-fluid"}
          (widgets/navigation-bar "/")
          [:h3 "Insights Advisor Mock"]
          [:div {:class "container-fluid" :style "padding:0px; margin:0px;"}
            [:div {:class "row justify-content-start"}
             (widgets/sidebar)
             [:div {:class "col-sm-2" :style "padding:20px"}
              "Text text"
             ]
            ]
          ]
          [:div {:style "height: 10ex"}]
          [:br]
          widgets/footer
     ] ; </div class="container">
    ] ; </body>
  ))

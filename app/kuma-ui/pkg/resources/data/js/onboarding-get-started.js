(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["onboarding-get-started"],{1390:function(t,a){t.exports="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAHgAAAB4CAYAAAA5ZDbSAAAABGdBTUEAALGPC/xhBQAABFVJREFUeAHtnM9O1EAYwKfLyhOYECJqEEw8evDiAcLRGOUtyMIDeNF48eDFB4ANb4GJZwwHLx48GY0KMWIIiU+AsLXfSrGUKdvszp/t56/JhnamO9/M79evnd2dYAwbBCAAAQhAAAIQgAAEIAABCEAAAhCAAAQgAAEIQAACEIAABCAAAQhAAAIQgAAEIAABCEAAAhCAgDGJbwjdbjf1HUND+51Ox4uLlgY4jKGaQLu6ym3N1ptttw0qaW354ZLXkZDBXvHGbxzB8R147QGCveKN3ziC4zvw2gMEe8Ubv/Fgs+j4Qw3Tg1YrMXduz5qZa1P9gPs/D82nL3um14vzdQCCHXsXufNzN85azfc/ft49Kwu5wy3aMe08c4vN2sqK9T73EeyT7hi0jWDHEuSZW95sZeVzfB3zDHZMViZUsuW35XyS5ThM7eYQXBtVvRNltiwTqliTqnIvuUWXiSg7RrAyoeXhILhMRNkxgpUJLQ8HwWUiyo69z6J7vV4f2aMHi8rQuRlOzsdNaxdbIYMvMlFV4j2DW62/1xBrsuzXDWuy7FworUmAW3RNUE09DcFNNVez3wiuCaqppyG4qeZq9tv7LLpmP9ScxposNSrtA2FNlp2LmtL8h/7igGxlxXqf+zyDfdIdg7YR7FiCbf2Vrcxx2MrmmGRVohmugjVZw3FrzLtYk9UYVTo6yjNYh8fKUSC4Eo2OCgTr8Fg5CgRXotFR4f1jUr7miDVZ9gsm52OvHb2UDB6d4Vi34D2DWZN1uX/WZF3Oh9oBBLhFDwDU9GoEN93ggP4jeACgplcjuOkGB/Tf+yx6QHx11azJUqf0/IBYk3Weh7oj2/orW1mogfMMDkU6UhwEOwZvW39lK3MctrI5JlmVaIarYE3WcNwa8y7WZDVGlY6O8gzW4bFyFAiuRKOjItgky/fvnjp0uB8FGeyeKS3GJrC+vp7KK3Y/QsQng0NQjhgDwRHhhwiN4BCUI8ZAcET4IUIjOATliDEQHBF+iNAIDkE5YgwER4QfIjSCQ1COGAPBEeGHCI3gEJQjxkgixg4aenNzc+Hk5ORJFvRemqbTErzdbi+urKzsBO1I4GD/RQZ3u92nmdy3mdjHuVzhfHx8vC11gZkHDac+gyVzRWSamuTb3n5ycPirD3h66qqZm51Jk8SkWSYvac1k9Rl8eltuidzvPw7M0dHv/kv2pSyz3To9J2hmhQrmLIPH/ffVnXcf+mKLYCcnr5iF+3eLRWO3v7a2NpIj9Rk8dsYCd2ikqyNwX4cKt7GxsSWTq6+7+0Zuy8Xt5vVpM39rxiRJ8np1dXW5WKdlX30GT0xMvMpk9WRCJULltiwv2ZcyqTs9R4vTc+NQn8Ey2mx+8Cz78yJ7lS/oXvZfgJ53Op2Xcp7G7b8QLOKKX3Scinwvmav145HGi5UxQQACEIAABCAAAQhAAAIQgAAEIAABCEAAAhCAAAQgAAEIQAACEIDAPwJ/ACqW/whzap9yAAAAAElFTkSuQmCC"},"22c7":function(t,a,e){},2486:function(t,a,e){"use strict";e.r(a);var s=function(){var t=this,a=t.$createElement,s=t._self._c||a;return s("div",{staticClass:"welcome welcome__step-1"},[t.title?s("p",{staticClass:"type-lg"},[t._v(" "+t._s(t.title)+" has been successfully installed but the cluster is currently empty. You are only a few steps away from deploying a modern distributed service mesh! ")]):t._e(),s("div",{staticClass:"app-setup"},[s("h3",{staticClass:"type-xl"},[t._v(" Let's set up your app ")]),s("div",{staticClass:"app-source-check"},[t.appSource&&"universal"===t.appSource||"kubernetes"===t.appSource||"k8s"===t.appSource?s("div",{staticClass:"app-source-check__inner flex items-center"},[s("div",{staticClass:"app-source-check__icon mr-4"},["universal"===t.appSource?s("img",{attrs:{src:e("1390"),alt:"Universal Icon"}}):"kubernetes"===t.appSource||"k8s"===t.appSource?s("img",{attrs:{src:e("9e14"),alt:"Kubernetes Icon"}}):t._e()]),s("div",{staticClass:"app-source-check__content px-4"},[s("p",[t._v(t._s(t.title)+" is running on "),s("span",{staticClass:"capitalize"},[t._v(t._s(t.appSource))])])]),t._m(0)]):s("div",{staticClass:"app-source-check--error"},[s("p",[t._v("The app was unable to determine "+t._s(t.productName)+"'s environment.")])])]),!1===t.tableDataLoadAttempted?s("div",{staticClass:"dataplane-loading-state flex -mx-2 mt-8"},[s("div",{staticClass:"px-2"},[s("KIcon",{attrs:{icon:"spinner",size:"36",color:"black"}})],1),t._m(1)]):t.tableData&&!1===t.tableDataIsEmpty?s("div",{staticClass:"mt-8"},[s("h2",{staticClass:"type-xl mb-2 pb-2"},[1===t.dataplaneCountForTitle?s("span",[t._v(" "+t._s(t.dataplaneCountForTitle)+" data plane proxy found: ")]):t.dataplaneCountForTitle<=10?s("span",[t._v(" "+t._s(t.dataplaneCountForTitle)+" data plane proxies found: ")]):s("span",[t._v(" "+t._s(t.dataplaneCountForTitle)+" data plane proxies found, including: ")])]),s("div",{staticClass:"data-table-wrapper"},[s("KTable",{attrs:{options:t.tableData},scopedSlots:t._u([{key:"status",fn:function(a){var e=a.rowValue;return[s("div",{staticClass:"entity-status",class:{"is-offline":"offline"===e.toLowerCase()||!1===e}},[s("span",{staticClass:"entity-status__dot"}),s("span",{staticClass:"entity-status__label"},[t._v(t._s(e))])])]}}])})],1),s("div",{staticClass:"md:flex items-center mt-4"},[s("div",{staticClass:"md:flex items-center md:mr-2 dataplane-global-status"},[s("div",[s("KButton",{staticClass:"mr-2",attrs:{appearance:"primary"},on:{click:t.reScanForDataplanes}},[t._v(" Refresh ")])],1),t.overallDpStatus?s("KAlert",{staticClass:"dataplane-status-alert",attrs:{appearance:"danger","alert-message":"Some data plane proxies appear to be offline."}}):t._e()],1),s("div",{staticClass:"md:ml-auto"},[s("KButton",{attrs:{to:{name:"setup-complete"},appearance:"primary"}},[t._v(" Next Step ")])],1)]),t.overallDpStatus?s("div",{staticClass:"dataplane-global-status__helper-text mt-8"},[s("h3",{staticClass:"type-xl mb-2 mt-4"},[t._v(" Offline data plane proxies ")]),s("p",[t._v(" This means your data plane proxy is not connected to the control plane at the moment. This might be due to a scheduled downtime or a network partitioning problem. ")])]):t._e()]):s("div",[s("div",{staticClass:"dataplane-fallback"},[s("div",{staticClass:"dataplane-fallback__inner flex -mx-4"},[t._m(2),s("div",{staticClass:"dataplane-fallback__content px-4"},[s("h3",{staticClass:"type-lg dataplane-fallback__title mb-2 pb-2"},[t._v(" No data plane proxies detected. ")]),s("p",{staticClass:"mb-2"},[t._v(" To bring your applications into "+t._s(t.productName)+" Service Mesh, you need to deploy data plane proxies (also known as sidecar proxies) next to them. ")])])])]),s("div",{staticClass:"dataplane-walkthrough my-4"},[s("div",[s("h3",{staticClass:"type-xl mb-2"},[t._v(" Next steps ")]),s("p",{staticClass:"mb-2"},[t._v(" The data plane proxy wizard will walk you through the creation of new sidecar proxies. ")]),s("div",{staticClass:"cols"},[s("KButton",{attrs:{to:{name:t.dataplaneWizardRoute},appearance:"primary"}},[t._v(" Get started ")])],1)])])]),s("footer",{staticClass:"extra-controls"},[s("KButton",{attrs:{to:{name:"global-overview",params:{mesh:"all",expandSidebar:!0}},appearance:"secondary",size:"small"},nativeOn:{click:function(a){return t.completeOnboarding()}}},[t._v(" Skip to Dashboard ")])],1)])])},i=[function(){var t=this,a=t.$createElement,s=t._self._c||a;return s("div",{staticClass:"px-4"},[s("img",{attrs:{src:e("57b2"),alt:"Checkmark Icon"}})])},function(){var t=this,a=t.$createElement,e=t._self._c||a;return e("div",{staticClass:"px-2"},[e("p",[t._v(" Waiting for data plane proxies to connect… ")])])},function(){var t=this,a=t.$createElement,s=t._self._c||a;return s("div",{staticClass:"dataplane-fallback__icon px-4"},[s("img",{attrs:{src:e("6ec4"),alt:"Dataplane Icon"}})])}],n=(e("07ac"),e("f3f3")),l=e("2f62"),o=e("8d3d"),c=e("c6ec"),r={name:"OnboardingStep1",metaInfo:function(){return{title:"Welcome to ".concat(this.productName,"!")}},data:function(){return{appSourceError:!1,productName:c["h"],tableDataLoadDelay:1500,tableDataIsEmpty:!0,tableDataLoadAttempted:!1,tableDataDataplaneCount:null,tableData:{headers:[{label:"Status",key:"status"},{label:"Name",key:"name"},{label:"Mesh",key:"mesh"}],data:[]},pageSize:10}},computed:Object(n["a"])(Object(n["a"])({},Object(l["c"])({title:"config/getTagline",appSource:"config/getEnvironment"})),{},{dataplaneCountForTitle:function(){var t=this.tableDataDataplaneCount;return t&&t>10?"10+":t},overallDpStatus:function(){return this.$store.getters.getAnyDpOffline},dataplaneWizardRoute:function(){return"universal"===this.appSource?"universal-dataplane":"kubernetes-dataplane"}}),beforeMount:function(){this.bootstrap()},methods:{bootstrap:function(){this.isLoading=!0,this.isEmpty=!1,this.getDataplaneTableData(),this.completeOnboarding()},reScanForDataplanes:function(){this.tableDataIsEmpty=!1,this.tableDataLoadAttempted=!1,this.getDataplaneTableData()},getDataplaneTableData:function(){var t=this,a={size:this.pageSize};this.$store.dispatch("getAllDataplanes",a).then((function(){var a=Object.values(t.$store.getters.getDataplanesList);a&&a.length>0?(t.tableDataDataplaneCount=a.length,t.tableData.data=a,t.tableDataLoadAttempted=!1,t.tableDataIsEmpty=!1,setTimeout((function(){t.tableDataLoadAttempted=!0}),t.tableDataLoadDelay)):(t.tableDataLoadAttempted=!0,t.tableDataIsEmpty=!0)}))},completeOnboarding:function(){this.$store.dispatch("updateOnboardingStatus",!0),Object(o["b"])("kumaOnboardingComplete",!0)}}},p=r,g=(e("45f7"),e("2877")),u=Object(g["a"])(p,s,i,!1,null,null,null);a["default"]=u.exports},"45f7":function(t,a,e){"use strict";e("22c7")},"57b2":function(t,a){t.exports="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PHBhdGggZD0iTTYuNzUgMTVMMjEuMDQ3LjcwM2EuOTk2Ljk5NiAwIDAxMS40MTUuMDA5bC44MjYuODI2Yy4zOTMuMzkzLjM5OCAxLjAyNi4wMDQgMS40Mkw3LjQ1OCAxOC43OTJhMS4wMDIgMS4wMDIgMCAwMS0xLjQxOC0uMDAyTC43MSAxMy40NmExIDEgMCAwMS4wMDItMS40MjJsLjgyNi0uODI2YTEuMDA5IDEuMDA5IDAgMDExLjQxNS0uMDA5TDYuNzUgMTV6IiBmaWxsPSIjMTE1NUNCIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiLz48L3N2Zz4K"},"6ec4":function(t,a){t.exports="data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPHN2ZyB3aWR0aD0iNDBweCIgaGVpZ2h0PSIzMnB4IiB2aWV3Qm94PSIwIDAgNDAgMzIiIHZlcnNpb249IjEuMSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayI+CiAgICA8IS0tIEdlbmVyYXRvcjogU2tldGNoIDYxICg4OTU4MSkgLSBodHRwczovL3NrZXRjaC5jb20gLS0+CiAgICA8dGl0bGU+aWNuLS1kYXRhcGxhbmU8L3RpdGxlPgogICAgPGRlc2M+Q3JlYXRlZCB3aXRoIFNrZXRjaC48L2Rlc2M+CiAgICA8ZyBpZD0iUGFnZS0xIiBzdHJva2U9Im5vbmUiIHN0cm9rZS13aWR0aD0iMSIgZmlsbD0ibm9uZSIgZmlsbC1ydWxlPSJldmVub2RkIj4KICAgICAgICA8ZyBpZD0iMXN0LXN0ZXAiIHRyYW5zZm9ybT0idHJhbnNsYXRlKC00MjQuMDAwMDAwLCAtNDMyLjAwMDAwMCkiIGZpbGw9IiM5NTlGQTYiPgogICAgICAgICAgICA8ZyBpZD0iQS1sZXQncy1zZXR1cC15b3VyLWFwcCIgdHJhbnNmb3JtPSJ0cmFuc2xhdGUoNDAwLjAwMDAwMCwgMjUwLjAwMDAwMCkiPgogICAgICAgICAgICAgICAgPGcgaWQ9ImJhbm5lci0tbm8tZGF0YS1wbGFuZXMiIHRyYW5zZm9ybT0idHJhbnNsYXRlKDAuMDAwMDAwLCAxNTguMDAwMDAwKSI+CiAgICAgICAgICAgICAgICAgICAgPHBhdGggZD0iTTMwLDM4IEwzMCw1NCBMNTgsNTQgTDU4LDQ4IEw2MCw0OCBMNjAsNTQgQzYwLDU1LjEwNDU2OTUgNTkuMTA0NTY5NSw1NiA1OCw1NiBMMzAsNTYgQzI4Ljg5NTQzMDUsNTYgMjgsNTUuMTA0NTY5NSAyOCw1NCBMMjgsMzggTDMwLDM4IFogTTMyLDM4IEwzMiw1MiBMNTYsNTIgTDU2LDQ4IEw1Nyw0OCBMNTcsNTMgTDMxLDUzIEwzMSwzOCBMMzIsMzggWiBNNDIsNDEgTDQyLDQ0IEw2NCw0NCBMNjQsNDYgTDQyLDQ2IEw0Miw0OSBMMzcsNDUgTDQyLDQxIFogTTU3LDI3IEw1Nyw0MiBMNTYsNDIgTDU2LDI4IEwzMiwyOCBMMzIsMzIgTDMxLDMyIEwzMSwyNyBMNTcsMjcgWiBNNTgsMjQgQzU5LjEwNDU2OTUsMjQgNjAsMjQuODk1NDMwNSA2MCwyNiBMNjAsNDIgTDU4LDQyIEw1OCwyNiBMMzAsMjYgTDMwLDMyIEwyOCwzMiBMMjgsMjYgQzI4LDI0Ljg5NTQzMDUgMjguODk1NDMwNSwyNCAzMCwyNCBMNTgsMjQgWiBNNDYsMzEgTDUxLDM1IEw0NiwzOSBMNDYsMzYgTDI0LDM2IEwyNCwzNCBMNDYsMzQgTDQ2LDMxIFoiIGlkPSJpY24tLWRhdGFwbGFuZSI+PC9wYXRoPgogICAgICAgICAgICAgICAgPC9nPgogICAgICAgICAgICA8L2c+CiAgICAgICAgPC9nPgogICAgPC9nPgo8L3N2Zz4="},"8d3d":function(t,a,e){"use strict";function s(t){var a=localStorage.getItem(t)||null;if(a)return JSON.parse(a)}function i(t,a){return localStorage.setItem(t,a)}e.d(a,"a",(function(){return s})),e.d(a,"b",(function(){return i}))},"9e14":function(t,a,e){t.exports=e.p+"img/icon-k8s.fb248510.png"}}]);
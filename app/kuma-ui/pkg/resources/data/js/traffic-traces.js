(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["traffic-traces"],{"1d3a":function(t,e,n){"use strict";n.d(e,"a",(function(){return c}));n("b0c0"),n("d3b7"),n("96cf");var a=n("c964"),r=n("d0ff"),i=n("c6ec");function o(t){return Object(r["a"])(t).sort((function(t,e){return t.name>e.name||t.name===e.name&&t.mesh>e.mesh?1:-1}))}var s=function(t){return 0!==t.total&&t.items&&t.items.length>0?o(t.items):[]};function l(t){var e=t.getSingleEntity,n=t.getAllEntities,a=t.getAllEntitiesFromMesh,r=t.mesh,i=t.query,o=t.size,s=t.offset,l={size:o,offset:s};return r&&"all"!==r?e&&i&&i.length&&"all"!==r?e(r,i,l):a?a(r,l):Promise.resolve():n(l)}function c(t){return u.apply(this,arguments)}function u(){return u=Object(a["a"])(regeneratorRuntime.mark((function t(e){var n,a,r,o,c,u,p,f,d;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return n=e.getSingleEntity,a=e.getAllEntities,r=e.getAllEntitiesFromMesh,o=e.mesh,c=e.query,u=e.size,p=void 0===u?i["f"]:u,f=e.offset,t.next=3,l({getSingleEntity:n,getAllEntities:a,getAllEntitiesFromMesh:r,mesh:o,query:c,size:p,offset:f});case 3:if(d=t.sent,d){t.next=6;break}return t.abrupt("return",{data:[],next:!1});case 6:return t.abrupt("return",{data:d.items?s(d):[d],next:Boolean(d.next)});case 7:case"end":return t.stop()}}),t)}))),u.apply(this,arguments)}},"23d6":function(t,e,n){"use strict";n("b91b")},3489:function(t,e,n){"use strict";n.r(e);var a=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"traffic-traces"},[n("FrameSkeleton",[n("DataOverview",{attrs:{"page-size":t.pageSize,"has-error":t.hasError,"is-loading":t.isLoading,"empty-state":t.empty_state,"table-data":t.tableData,"table-data-is-empty":t.tableDataIsEmpty,"table-data-function-text":"View","table-data-row":"name",next:t.next},on:{tableAction:t.tableAction,loadData:function(e){return t.loadData(e)}},scopedSlots:t._u([{key:"additionalControls",fn:function(){return[t.$route.query.ns?n("KButton",{staticClass:"back-button",attrs:{appearance:"primary",size:"small",to:{name:"traffic-traces"}}},[n("span",{staticClass:"custom-control-icon"},[t._v(" ← ")]),t._v(" View All ")]):t._e()]},proxy:!0}])}),!1===t.isEmpty?n("Tabs",{attrs:{"has-error":t.hasError,"is-loading":t.isLoading,tabs:t.tabs,"initial-tab-override":"overview"},scopedSlots:t._u([{key:"tabHeader",fn:function(){return[n("div",[n("h3",[t._v(t._s(t.tabGroupTitle))])]),n("div",[n("EntityURLControl",{attrs:{url:t.shareUrl}})],1)]},proxy:!0},{key:"overview",fn:function(){return[n("LabelList",{attrs:{"has-error":t.entityHasError,"is-loading":t.entityIsLoading,"is-empty":t.entityIsEmpty}},[n("div",[n("ul",t._l(t.entity,(function(e,a){return n("li",{key:a},[n("h4",[t._v(t._s(a))]),n("p",[t._v(" "+t._s(e)+" ")])])})),0)])])]},proxy:!0},{key:"yaml",fn:function(){return[n("YamlView",{attrs:{title:t.entityOverviewTitle,"has-error":t.entityHasError,"is-loading":t.entityIsLoading,"is-empty":t.entityIsEmpty,content:t.rawEntity}})]},proxy:!0}],null,!1,4098209877)}):t._e()],1)],1)},r=[],i=(n("99af"),n("b0c0"),n("d3b7"),n("96cf"),n("c964")),o=n("bc1e"),s=n("0f82"),l=n("1d3a"),c=n("6663"),u=n("1d10"),p=n("2778"),f=n("251b"),d=n("ff9d"),y=n("0ada"),h=n("c6ec"),m={name:"TrafficTraces",metaInfo:{title:"Traffic Traces"},components:{EntityURLControl:c["a"],FrameSkeleton:u["a"],DataOverview:p["a"],Tabs:f["a"],YamlView:d["a"],LabelList:y["a"]},data:function(){return{isLoading:!0,isEmpty:!1,hasError:!1,entityIsLoading:!0,entityIsEmpty:!1,entityHasError:!1,tableDataIsEmpty:!1,empty_state:{title:"No Data",message:"There are no Traffic Traces present."},tableData:{headers:[{key:"actions",hideLabel:!0},{label:"Name",key:"name"},{label:"Mesh",key:"mesh"},{label:"Type",key:"type"}],data:[]},tabs:[{hash:"#overview",title:"Overview"},{hash:"#yaml",title:"YAML"}],entity:[],rawEntity:null,firstEntity:null,pageSize:h["f"],next:null}},computed:{tabGroupTitle:function(){var t=this.entity;return t?"Traffic Trace: ".concat(t.name):null},entityOverviewTitle:function(){var t=this.entity;return t?"Entity Overview for ".concat(t.name):null},shareUrl:function(){var t=this,e="".concat(window.location.origin,"#"),n=this.entity,a=function(){return t.$route.query.ns?t.$route.fullPath:"".concat(e).concat(t.$route.fullPath,"?ns=").concat(n.name)};return a()}},watch:{$route:function(t,e){this.init()}},beforeMount:function(){this.init()},methods:{init:function(){this.loadData()},tableAction:function(t){var e=t;this.getEntity(e)},loadData:function(){var t=arguments,e=this;return Object(i["a"])(regeneratorRuntime.mark((function n(){var a,r,i,c,u,p,f,d;return regeneratorRuntime.wrap((function(n){while(1)switch(n.prev=n.next){case 0:return a=t.length>0&&void 0!==t[0]?t[0]:"0",e.isLoading=!0,r=e.$route.query.ns||null,i=e.$route.params.mesh||null,n.prev=4,n.next=7,Object(l["a"])({getSingleEntity:s["a"].getTrafficTrace.bind(s["a"]),getAllEntities:s["a"].getAllTrafficTraces.bind(s["a"]),getAllEntitiesFromMesh:s["a"].getAllTrafficTracesFromMesh.bind(s["a"]),mesh:i,query:r,size:e.pageSize,offset:a});case 7:c=n.sent,u=c.data,p=c.next,e.next=p,u.length?(e.tableData.data=u,e.tableDataIsEmpty=!1,e.isEmpty=!1,f=["type","name","mesh"],d=u[0],e.entity=Object(o["d"])(d,f),e.rawEntity=Object(o["j"])(d)):(e.tableData.data=[],e.tableDataIsEmpty=!0,e.isEmpty=!0,e.entityIsEmpty=!0),n.next=19;break;case 14:n.prev=14,n.t0=n["catch"](4),e.hasError=!0,e.isEmpty=!0,console.error(n.t0);case 19:return n.prev=19,e.isLoading=!1,e.entityIsLoading=!1,n.finish(19);case 23:case"end":return n.stop()}}),n,null,[[4,14,19,23]])})))()},getEntity:function(t){var e=this;this.entityIsLoading=!0,this.entityIsEmpty=!1,this.entityHasError=!1;var n=this.$route.params.mesh;if(t){var a="all"===n?t.mesh:n;return s["a"].getTrafficTrace(a,t.name).then((function(t){if(t){var n=["type","name","mesh"];e.entity=Object(o["d"])(t,n),e.rawEntity=Object(o["j"])(t)}else e.entity=null,e.entityIsEmpty=!0})).catch((function(t){e.entityHasError=!0,console.error(t)})).finally((function(){setTimeout((function(){e.entityIsLoading=!1}),"500")}))}setTimeout((function(){e.entityIsEmpty=!0,e.entityIsLoading=!1}),"500")}}},v=m,b=n("2877"),g=Object(b["a"])(v,a,r,!1,null,null,null);e["default"]=g.exports},"62e5":function(t,e){var n;n=function(){function t(t,e){var n,a,r,i,o,s,l,c,u;null==e&&(e=""),r="",o=t.length,s=null,a=0,i=0;while(i<o){if(n=t.charAt(i),"\\"===n)r+=t.slice(i,+(i+1)+1||9e9),i++;else if("("===n)if(i<o-2)if(c=t.slice(i,+(i+2)+1||9e9),"(?:"===c)i+=2,r+=c;else if("(?<"===c){a++,i+=2,l="";while(i+1<o){if(u=t.charAt(i+1),">"===u){r+="(",i++,l.length>0&&(null==s&&(s={}),s[l]=a);break}l+=u,i++}}else r+=n,a++;else r+=n;else r+=n;i++}this.rawRegex=t,this.cleanedRegex=r,this.regex=new RegExp(this.cleanedRegex,"g"+e.replace("g","")),this.mapping=s}return t.prototype.regex=null,t.prototype.rawRegex=null,t.prototype.cleanedRegex=null,t.prototype.mapping=null,t.prototype.exec=function(t){var e,n,a,r;if(this.regex.lastIndex=0,n=this.regex.exec(t),null==n)return null;if(null!=this.mapping)for(a in r=this.mapping,r)e=r[a],n[a]=n[e];return n},t.prototype.test=function(t){return this.regex.lastIndex=0,this.regex.test(t)},t.prototype.replace=function(t,e){return this.regex.lastIndex=0,t.replace(this.regex,e)},t.prototype.replaceAll=function(t,e,n){var a;null==n&&(n=0),this.regex.lastIndex=0,a=0;while(this.regex.test(t)&&(0===n||a<n))this.regex.lastIndex=0,t=t.replace(this.regex,e),a++;return[t,a]},t}(),t.exports=n},"64cff":function(t,e,n){"use strict";n("7ac1")},6663:function(t,e,n){"use strict";var a=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"entity-url-control"},[t.shouldDisplay?n("KClipboardProvider",{scopedSlots:t._u([{key:"default",fn:function(e){var a=e.copyToClipboard;return[n("KPop",{attrs:{placement:"bottom"},scopedSlots:t._u([{key:"content",fn:function(){return[n("div",[n("p",[t._v(t._s(t.confirmationText))])])]},proxy:!0}],null,!0)},[n("KButton",{attrs:{appearance:"secondary",size:"small"},on:{click:function(){a(t.url)}},scopedSlots:t._u([{key:"icon",fn:function(){return[n("KIcon",{attrs:{"view-box":"0 0 16 16",icon:"externalLink"}})]},proxy:!0}],null,!0)},[t._v(" "+t._s(t.copyButtonText)+" ")])],1)]}}],null,!1,244323940)}):t._e()],1)},r=[],i=n("a026"),o=i["a"].extend({name:"EntityURLControl",props:{url:{type:String,required:!0},copyButtonText:{type:String,default:"Copy URL"},confirmationText:{type:String,default:"URL copied to clipboard!"}},computed:{shouldDisplay:function(){var t=this.$route.params.mesh||null;return!(!t||"all"===t)}}}),s=o,l=n("2877"),c=Object(l["a"])(s,a,r,!1,null,null,null);e["a"]=c.exports},"6d8a":function(t,e,n){var a,r;r=n("62e5"),a=function(){var t;function e(){}return e.LIST_ESCAPEES=["\\","\\\\",'\\"','"',"\0","","","","","","","","\b","\t","\n","\v","\f","\r","","","","","","","","","","","","","","","","","","",(t=String.fromCharCode)(133),t(160),t(8232),t(8233)],e.LIST_ESCAPED=["\\\\",'\\"','\\"','\\"',"\\0","\\x01","\\x02","\\x03","\\x04","\\x05","\\x06","\\a","\\b","\\t","\\n","\\v","\\f","\\r","\\x0e","\\x0f","\\x10","\\x11","\\x12","\\x13","\\x14","\\x15","\\x16","\\x17","\\x18","\\x19","\\x1a","\\e","\\x1c","\\x1d","\\x1e","\\x1f","\\N","\\_","\\L","\\P"],e.MAPPING_ESCAPEES_TO_ESCAPED=function(){var t,n,a,r;for(a={},t=n=0,r=e.LIST_ESCAPEES.length;0<=r?n<r:n>r;t=0<=r?++n:--n)a[e.LIST_ESCAPEES[t]]=e.LIST_ESCAPED[t];return a}(),e.PATTERN_CHARACTERS_TO_ESCAPE=new r("[\\x00-\\x1f]|Â|Â |â¨|â©"),e.PATTERN_MAPPING_ESCAPEES=new r(e.LIST_ESCAPEES.join("|").split("\\").join("\\\\")),e.PATTERN_SINGLE_QUOTING=new r("[\\s'\":{}[\\],&*#?]|^[-?|<>=!%@`]"),e.requiresDoubleQuoting=function(t){return this.PATTERN_CHARACTERS_TO_ESCAPE.test(t)},e.escapeWithDoubleQuotes=function(t){var e;return e=this.PATTERN_MAPPING_ESCAPEES.replace(t,function(t){return function(e){return t.MAPPING_ESCAPEES_TO_ESCAPED[e]}}(this)),'"'+e+'"'},e.requiresSingleQuoting=function(t){return this.PATTERN_SINGLE_QUOTING.test(t)},e.escapeWithSingleQuotes=function(t){return"'"+t.replace(/'/g,"''")+"'"},e}(),t.exports=a},"7ac1":function(t,e,n){},b91b:function(t,e,n){},e80b:function(t,e,n){var a=n("6d8a"),r="  ";function i(t){var e=typeof t;return t instanceof Array?"array":"string"==e?"string":"boolean"==e?"boolean":"number"==e?"number":"undefined"==e||null===t?"null":"hash"}function o(t,e){var n=i(t);switch(n){case"array":s(t,e);break;case"hash":l(t,e);break;case"string":u(t,e);break;case"null":e.push("null");break;case"number":e.push(t.toString());break;case"boolean":e.push(t?"true":"false");break}}function s(t,e){for(var n=0;n<t.length;n++){var a=t[n],i=[];o(a,i);for(var s=0;s<i.length;s++)e.push((0==s?"- ":r)+i[s])}}function l(t,e){for(var n in t){var a=[];if(t.hasOwnProperty(n)){var s=t[n];o(s,a);var l=i(s);if("string"==l||"null"==l||"number"==l||"boolean"==l)e.push(c(n)+": "+a[0]);else{e.push(c(n)+": ");for(var u=0;u<a.length;u++)e.push(r+a[u])}}}}function c(t){return t.match(/^[\w]+$/)?t:a.requiresDoubleQuoting(t)?a.escapeWithDoubleQuotes(t):a.requiresSingleQuoting(t)?a.escapeWithSingleQuotes(t):t}function u(t,e){e.push(c(t))}var p=function(t){"string"==typeof t&&(t=JSON.parse(t));var e=[];return o(t,e),e.join("\n")};t.exports=p},ff9d:function(t,e,n){"use strict";var a=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"yaml-view"},[t.isReady?n("div",{staticClass:"yaml-view-content"},[t.isLoading||t.isEmpty?t._e():n("KCard",{attrs:{title:t.title,"border-variant":"noBorder"},scopedSlots:t._u([{key:"body",fn:function(){return[n("KTabs",{key:t.environment,attrs:{tabs:t.tabs},scopedSlots:t._u([{key:"universal",fn:function(){return[n("KClipboardProvider",{scopedSlots:t._u([{key:"default",fn:function(e){var a=e.copyToClipboard;return[n("KPop",{attrs:{placement:"bottom"},scopedSlots:t._u([{key:"content",fn:function(){return[n("div",[n("p",[t._v("Entity copied to clipboard!")])])]},proxy:!0}],null,!0)},[n("KButton",{staticClass:"copy-button",attrs:{appearance:"primary",size:"small"},on:{click:function(){a(t.yamlContent.universal)}}},[t._v(" Copy Universal YAML ")])],1)]}}],null,!1,1536634960)}),n("prism",{staticClass:"code-block",attrs:{language:"yaml",code:t.yamlContent.universal}})]},proxy:!0},{key:"kubernetes",fn:function(){return[n("KClipboardProvider",{scopedSlots:t._u([{key:"default",fn:function(e){var a=e.copyToClipboard;return[n("KPop",{attrs:{placement:"bottom"},scopedSlots:t._u([{key:"content",fn:function(){return[n("div",[n("p",[t._v("Entity copied to clipboard!")])])]},proxy:!0}],null,!0)},[n("KButton",{staticClass:"copy-button",attrs:{appearance:"primary",size:"small"},on:{click:function(){a(t.yamlContent.kubernetes)}}},[t._v(" Copy Kubernetes YAML ")])],1)]}}],null,!1,2265429040)}),n("prism",{staticClass:"code-block",attrs:{language:"yaml",code:t.yamlContent.kubernetes}})]},proxy:!0}],null,!1,661975406),model:{value:t.activeTab.hash,callback:function(e){t.$set(t.activeTab,"hash",e)},expression:"activeTab.hash"}})]},proxy:!0}],null,!1,2034136171)})],1):t._e(),!0===t.loaders?n("div",[t.isLoading?n("KEmptyState",{attrs:{"cta-is-hidden":""},scopedSlots:t._u([{key:"title",fn:function(){return[n("div",{staticClass:"card-icon mb-3"},[n("KIcon",{attrs:{icon:"spinner",color:"rgba(0, 0, 0, 0.1)",size:"42"}})],1),t._v(" Data Loading... ")]},proxy:!0}],null,!1,3263214496)}):t._e(),t.isEmpty&&!t.isLoading?n("KEmptyState",{attrs:{"cta-is-hidden":""},scopedSlots:t._u([{key:"title",fn:function(){return[n("div",{staticClass:"card-icon mb-3"},[n("KIcon",{staticClass:"kong-icon--centered",attrs:{color:"var(--yellow-200)",icon:"warning","secondary-color":"var(--black-75)",size:"42"}})],1),t._v(" There is no data to display. ")]},proxy:!0}],null,!1,1612658095)}):t._e(),t.hasError?n("KEmptyState",{attrs:{"cta-is-hidden":""},scopedSlots:t._u([{key:"title",fn:function(){return[n("div",{staticClass:"card-icon mb-3"},[n("KIcon",{staticClass:"kong-icon--centered",attrs:{color:"var(--yellow-200)",icon:"warning","secondary-color":"var(--black-75)",size:"42"}})],1),t._v(" An error has occurred while trying to load this data. ")]},proxy:!0}],null,!1,822917942)}):t._e()],1):t._e()])},r=[],i=(n("caad"),n("a15b"),n("b0c0"),n("4fad"),n("ac1f"),n("2532"),n("1276"),n("f3f3")),o=n("2f62"),s=n("2ccf"),l=n.n(s),c=n("e80b"),u=n.n(c),p={name:"YamlView",components:{prism:l.a},props:{title:{type:String,default:null},content:{type:Object,default:null},loaders:{type:Boolean,default:!0},isLoading:{type:Boolean,default:!1},hasError:{type:Boolean,default:!1},isEmpty:{type:Boolean,default:!1}},data:function(){return{tabs:[{hash:"#universal",title:"Universal"},{hash:"#kubernetes",title:"Kubernetes"}]}},computed:Object(i["a"])(Object(i["a"])({},Object(o["c"])({environment:"config/getEnvironment"})),{},{isReady:function(){return!this.isEmpty&&!this.hasError&&!this.isLoading},activeTab:{get:function(){var t=this.environment;return{hash:"#".concat(t),nohash:t}},set:function(t){return{hash:"#".concat(t),nohash:t}}},yamlContent:function(){var t=this,e=this.content,n=function(){var e={},n=Object.assign({},t.content),a=n.name,r=n.mesh,i=n.type,o=function(){var e=Object.assign({},t.content);return delete e.type,delete e.mesh,delete e.name,!!(e&&Object.entries(e).length>0)&&e};if(e.apiVersion="kuma.io/v1alpha1",e.kind=i,void 0!==r&&(e.mesh=n.mesh),null!==a&&void 0!==a&&a.includes(".")){var s=a.split("."),l=s.pop(),c=s.join(".");e.metadata={name:c,namespace:l}}else e.metadata={name:a};return o()&&(e.spec=o()),e},a={universal:u()(e),kubernetes:u()(n())};return a}})},f=p,d=(n("23d6"),n("64cff"),n("2877")),y=Object(d["a"])(f,a,r,!1,null,"d189da22",null);e["a"]=y.exports}}]);
(function(){"use strict";var e={6141:function(e,t,r){var a=r(144),n=r(998),s=r(9582),o=r(266),i=r(2150),l=r(2480),c=r(2928),u=r(3305),p=r(1713),d=r(8781),v=function(){var e=this,t=e._self._c;return t(n.Z,[t(c.Z,[t(i.Z,[t(p.Z,{attrs:{dense:""}},[t(o.Z,{staticClass:"ma-2",attrs:{align:"center"}},[t(i.Z,{staticClass:"text-h4 mb-2"},[e._v(" Photovoltaik ")]),t(s.Z,{staticClass:"ma-2 px-6 py-4"},[t(p.Z,[t(o.Z,[t(u.Z,{attrs:{color:e.percentColor,width:"10",rotate:"-90",size:"100",value:e.power.batteryPercent,indeterminate:null==e.power.batteryPercent}},[null!==e.power.batteryPercent?t("div",{staticClass:"text-h4"},[e._v(" "+e._s(e.power.batteryPercent+"%")+" ")]):e._e()])],1),null!=e.power.batteryPercent?t(p.Z,[t(o.Z,{staticClass:"text-left",attrs:{"align-self":"center"}},[t("StatsText",{attrs:{value:e.power.batteryPower,negResponse:"Akku entlädt",posResponse:"Akku lädt"}}),t("StatsText",{attrs:{value:e.power.gridPower,negResponse:"Netzbezug",posResponse:"Netzeinspeisung"}}),t("div",[e._v("PV Erzeugung")]),t("div",[e._v("Stromverbrauch")])],1),t(o.Z,{staticClass:"text-right",attrs:{"align-self":"center"}},[e.power.batteryPower?t("div",[e._v(e._s(Math.abs(e.power.batteryPower))+" W")]):e._e(),t("div",[e._v(e._s(Math.abs(e.power.gridPower))+" W")]),t("div",[e._v(e._s(Math.abs(e.power.pvPower))+" W")]),t("div",[e._v(e._s(Math.abs(e.power.powerConsumption))+" W")])])],1):e._e()],1)],1),t("br"),t(s.Z,{staticClass:"ma-2 pa-2"},[t("div",{staticClass:"text-h6"},[e._v("Prognose")]),t(d.Z,{staticClass:"ma-2",attrs:{padding:"10",value:e.pvForecast.pv_estimate,labels:e.pvForecast.period_end,smooth:"10",height:"150","auto-draw":"","stroke-linecap":"round"}}),t(l.Z,{staticClass:"ma-2",attrs:{dense:"","mobile-breakpoint":"0","hide-default-footer":"",headers:e.headers,items:e.dailyForecast}})],1)],1)],1)],1)],1)],1)},f=[],h=(r(7658),function(){var e=this,t=e._self._c;return e.response?t("div",[e._v(" "+e._s(e.response)+" ")]):t("div",[e.value<0?t("div",{staticClass:"red--text"},[e._v(" "+e._s(e.negResponse)+" ")]):e.value>0?t("div",{staticClass:"green--text"},[e._v(" "+e._s(e.posResponse)+" ")]):e._e()])}),w=[],y={name:"StatsText",props:{value:null,negResponse:String,posResponse:String,response:null}},g=y,b=r(1001),_=(0,b.Z)(g,h,w,!1,null,null,null),m=_.exports,P={name:"App",components:{StatsText:m},data(){return{power:{batteryPercent:null,gridPower:null,pvPower:null,powerConsumption:null,batteryPower:null},interval:0,pvForecast:{pv_estimate:[],period_end:[],period:""},dailyForecast:[],headers:[{text:"Tag",value:"day"},{text:"Prognose (kWh)",value:"estimate",align:"end"}]}},computed:{percentColor(){return this.power.batteryPercent>65?"green":this.power.batteryPercent>35?"orange":"red"}},created(){this.fetchData(),setInterval(this.fetchData,2e3),this.fetchForecast(),this.fetchDailyForecast()},methods:{async fetchData(){let e=await fetch("api/data");this.power=await e.json()},async fetchForecast(){let e=await fetch("api/forecast");const t=await e.json();this.pvForecast={pv_estimate:t.pv_estimate,period_end:t.period_end.map((e=>{const t=new Date(e);return t.getHours()+":"+(0==t.getMinutes()?"00":t.getMinutes())})),period:t.period}},async fetchDailyForecast(){const e=await fetch("api/forecast/daily");let t=await e.json();t=t.dailyForecast,Object.keys(t).forEach(((e,r)=>{0==r?this.dailyForecast.push({day:"Heute",estimate:t[e].estimate}):1==r?this.dailyForecast.push({day:"Morgen",estimate:t[e].estimate}):this.dailyForecast.push(t[e])}))}}},Z=P,x=(0,b.Z)(Z,v,f,!1,null,null,null),C=x.exports,F=r(1705);a.ZP.use(F.Z);var k=new F.Z({});a.ZP.config.productionTip=!1,new a.ZP({vuetify:k,render:e=>e(C)}).$mount("#app")}},t={};function r(a){var n=t[a];if(void 0!==n)return n.exports;var s=t[a]={exports:{}};return e[a](s,s.exports,r),s.exports}r.m=e,function(){var e=[];r.O=function(t,a,n,s){if(!a){var o=1/0;for(u=0;u<e.length;u++){a=e[u][0],n=e[u][1],s=e[u][2];for(var i=!0,l=0;l<a.length;l++)(!1&s||o>=s)&&Object.keys(r.O).every((function(e){return r.O[e](a[l])}))?a.splice(l--,1):(i=!1,s<o&&(o=s));if(i){e.splice(u--,1);var c=n();void 0!==c&&(t=c)}}return t}s=s||0;for(var u=e.length;u>0&&e[u-1][2]>s;u--)e[u]=e[u-1];e[u]=[a,n,s]}}(),function(){r.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return r.d(t,{a:t}),t}}(),function(){r.d=function(e,t){for(var a in t)r.o(t,a)&&!r.o(e,a)&&Object.defineProperty(e,a,{enumerable:!0,get:t[a]})}}(),function(){r.g=function(){if("object"===typeof globalThis)return globalThis;try{return this||new Function("return this")()}catch(e){if("object"===typeof window)return window}}()}(),function(){r.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)}}(),function(){r.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})}}(),function(){var e={143:0};r.O.j=function(t){return 0===e[t]};var t=function(t,a){var n,s,o=a[0],i=a[1],l=a[2],c=0;if(o.some((function(t){return 0!==e[t]}))){for(n in i)r.o(i,n)&&(r.m[n]=i[n]);if(l)var u=l(r)}for(t&&t(a);c<o.length;c++)s=o[c],r.o(e,s)&&e[s]&&e[s][0](),e[s]=0;return r.O(u)},a=self["webpackChunkvarta"]=self["webpackChunkvarta"]||[];a.forEach(t.bind(null,0)),a.push=t.bind(null,a.push.bind(a))}();var a=r.O(void 0,[998],(function(){return r(6141)}));a=r.O(a)})();
//# sourceMappingURL=app.2bfb3c43.js.map
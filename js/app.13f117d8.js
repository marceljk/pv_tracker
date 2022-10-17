(function(){"use strict";var t={3268:function(t,e,r){var n=r(144),o=r(9306),a=r(266),i=r(2150),u=r(2928),c=r(1253),l=r(1713),s=function(){var t=this,e=t._self._c;return e(o.Z,[e(u.Z,[e(i.Z,[e(l.Z,{attrs:{dense:"",justify:"center"}},[e(a.Z,{staticClass:"ma-4",attrs:{align:"center"}},[e(i.Z,{staticClass:"text-h2 mb-5"},[t._v(" Photovoltaik ")]),e(c.Z,{attrs:{color:t.percentColor,width:"20",rotate:"-90",size:"300",value:t.power.batteryPercent,indeterminate:null==t.power.batteryPercent}},[null!==t.power.batteryPercent?e("div",{staticClass:"text-h3"},[t._v(" "+t._s(t.power.batteryPercent+"%")+" ")]):t._e()]),t.power.batteryPower<0?e(i.Z,{staticClass:"text-h5 red--text"},[t._v(" Batterie entlädt: "+t._s(Math.abs(t.power.batteryPower))+" W ")]):t.power.batteryPower>0?e(i.Z,{staticClass:"text-h5 green--text"},[t._v(" Batterie ladet mit: "+t._s(Math.abs(t.power.batteryPower))+" W ")]):t._e()],1)],1)],1)],1)],1)},f=[],p={name:"App",components:{},data(){return{power:{batteryPercent:null,gridPower:null,pvPower:null,powerConsumption:null,batteryPower:null},interval:0}},computed:{percentColor(){return this.power.batteryPercent>65?"green":this.power.batteryPercent>35?"orange":"red"}},created(){this.fetchData(),setInterval(this.fetchData,2e3)},methods:{async fetchData(){let t=await fetch("http://192.168.2.155/api/data");this.power=await t.json()}}},v=p,d=r(1001),h=(0,d.Z)(v,s,f,!1,null,null,null),b=h.exports,w=r(8864);n.ZP.use(w.Z);var y=new w.Z({});n.ZP.config.productionTip=!1,new n.ZP({vuetify:y,render:t=>t(b)}).$mount("#app")}},e={};function r(n){var o=e[n];if(void 0!==o)return o.exports;var a=e[n]={exports:{}};return t[n](a,a.exports,r),a.exports}r.m=t,function(){var t=[];r.O=function(e,n,o,a){if(!n){var i=1/0;for(s=0;s<t.length;s++){n=t[s][0],o=t[s][1],a=t[s][2];for(var u=!0,c=0;c<n.length;c++)(!1&a||i>=a)&&Object.keys(r.O).every((function(t){return r.O[t](n[c])}))?n.splice(c--,1):(u=!1,a<i&&(i=a));if(u){t.splice(s--,1);var l=o();void 0!==l&&(e=l)}}return e}a=a||0;for(var s=t.length;s>0&&t[s-1][2]>a;s--)t[s]=t[s-1];t[s]=[n,o,a]}}(),function(){r.n=function(t){var e=t&&t.__esModule?function(){return t["default"]}:function(){return t};return r.d(e,{a:e}),e}}(),function(){r.d=function(t,e){for(var n in e)r.o(e,n)&&!r.o(t,n)&&Object.defineProperty(t,n,{enumerable:!0,get:e[n]})}}(),function(){r.g=function(){if("object"===typeof globalThis)return globalThis;try{return this||new Function("return this")()}catch(t){if("object"===typeof window)return window}}()}(),function(){r.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)}}(),function(){r.r=function(t){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(t,"__esModule",{value:!0})}}(),function(){var t={143:0};r.O.j=function(e){return 0===t[e]};var e=function(e,n){var o,a,i=n[0],u=n[1],c=n[2],l=0;if(i.some((function(e){return 0!==t[e]}))){for(o in u)r.o(u,o)&&(r.m[o]=u[o]);if(c)var s=c(r)}for(e&&e(n);l<i.length;l++)a=i[l],r.o(t,a)&&t[a]&&t[a][0](),t[a]=0;return r.O(s)},n=self["webpackChunkvarta"]=self["webpackChunkvarta"]||[];n.forEach(e.bind(null,0)),n.push=e.bind(null,n.push.bind(n))}();var n=r.O(void 0,[998],(function(){return r(3268)}));n=r.O(n)})();
//# sourceMappingURL=app.13f117d8.js.map
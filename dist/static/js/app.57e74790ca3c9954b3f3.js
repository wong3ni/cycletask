webpackJsonp([1],{NHnr:function(t,e,l){"use strict";Object.defineProperty(e,"__esModule",{value:!0});var a=l("7+uW"),o={render:function(){var t=this.$createElement,e=this._self._c||t;return e("div",{attrs:{id:"app"}},[e("router-view")],1)},staticRenderFns:[]};var r=l("VU/8")({name:"App"},o,!1,function(t){l("gsu9")},null,null).exports,n=l("/ocq"),i=l("PQ22").ip(),s={all:"http://"+i+":22222"},u=(l("nNlw"),l("mtWM")),c=function(t,e){return s[t]+e},f={post:function(t,e,l,a,o){return u.post(c(t,e),l,o).then(function(t){return t.data}).catch(function(t){a?a(t):console.error(t)})},get:function(t,e,l,a){return u.get(c(t,e),a).then(function(t){return t.data}).catch(function(t){l?l(t):console.log(t)})},eventsource:function(t,e,l){new EventSource(c(t,e)).onmessage=function(t){l(t)}}},d={data:function(){return{tablebar:[{label:"标签",prop:"tag"},{label:"请求地址",prop:"rurl"},{label:"转发地址",prop:"durl"},{label:"状态",prop:"state"}],tableData:[],dialogFormVisible:!1,form:{id:"",name:"",tag:"",rurl:"",durl:"",time:2,direction:""},formLabelWidth:"120px"}},mounted:function(){var t=this;f.eventsource("all","/cycletask/list",function(e){var l=JSON.parse(e.data);null!==l&&(l.forEach(function(t){t.state=!0===t.state?"正在运行":"已停止"}),t.tableData=l)})},methods:{controller:function(t,e){"已停止"===this.tableData[t].state?this.start(e,t):this.stop(e,t)},start:function(t,e){var l,a=this;console.log(this.tableData),(l=t,f.get("all","/api/start?tag="+l)).then(function(t){console.log(t.Code),0!==Number(t.Code)?alert("服务器异常"):a.tableData[e].state="正在运行"})},stop:function(t,e){var l,a=this;(l=t,f.get("all","/api/stop?tag="+l)).then(function(t){0!==Number(t.Code)?alert("服务器异常"):a.tableData[e].state="已停止"})},del:function(t){var e;this.tableData=this.tableData.filter(function(e){return e.tag!==t}),(e=t,f.get("all","/api/del?tag="+e)).then(function(t){0!==Number(t.Code)&&alert("服务器异常")})},add:function(){var t,e=this;console.log(this.form),(t=this.form,f.get("all","/api/add?id="+t.id+"&name="+t.name+"&tag="+t.tag+"&rurl="+t.rurl+"&durl="+t.durl+"&time="+t.time+"&direction="+t.direction)).then(function(t){0!==Number(t.Code)?(alert("服务器异常"),e.dialogFormVisible=!1):e.dialogFormVisible=!1})}}},m={render:function(){var t=this,e=t.$createElement,l=t._self._c||e;return l("div",[l("el-table",{staticStyle:{width:"100%"},attrs:{data:t.tableData}},[t._l(t.tablebar,function(t,e){return[l("el-table-column",{key:e,attrs:{align:"center",prop:t.prop,label:t.label}})]}),t._v(" "),l("el-table-column",{attrs:{label:"操作"},scopedSlots:t._u([{key:"default",fn:function(e){return[l("el-button",{attrs:{type:"text",size:"small"},on:{click:function(l){return t.controller(e.$index,e.row.tag)}}},[t._v(t._s("已停止"===e.row.state?"运行":"停止"))]),t._v(" "),l("el-button",{attrs:{type:"text",size:"small"},on:{click:function(l){return t.del(e.row.tag)}}},[t._v("删除")])]}}])})],2),t._v(" "),l("el-button",{attrs:{type:"text"},on:{click:function(e){t.dialogFormVisible=!0}}},[t._v("添加")]),t._v(" "),l("el-dialog",{attrs:{title:"添加信息",visible:t.dialogFormVisible},on:{"update:visible":function(e){t.dialogFormVisible=e}}},[l("el-form",{attrs:{model:t.form}},[l("el-form-item",{attrs:{label:"编号","label-width":t.formLabelWidth}},[l("el-input",{attrs:{autocomplete:"off"},model:{value:t.form.id,callback:function(e){t.$set(t.form,"id",e)},expression:"form.id"}})],1),t._v(" "),l("el-form-item",{attrs:{label:"名称","label-width":t.formLabelWidth}},[l("el-input",{attrs:{autocomplete:"off"},model:{value:t.form.name,callback:function(e){t.$set(t.form,"name",e)},expression:"form.name"}})],1),t._v(" "),l("el-form-item",{attrs:{label:"标签","label-width":t.formLabelWidth}},[l("el-input",{attrs:{autocomplete:"off"},model:{value:t.form.tag,callback:function(e){t.$set(t.form,"tag",e)},expression:"form.tag"}})],1),t._v(" "),l("el-form-item",{attrs:{label:"请求地址","label-width":t.formLabelWidth}},[l("el-input",{attrs:{autocomplete:"off"},model:{value:t.form.rurl,callback:function(e){t.$set(t.form,"rurl",e)},expression:"form.rurl"}})],1),t._v(" "),l("el-form-item",{attrs:{label:"转发地址","label-width":t.formLabelWidth}},[l("el-input",{attrs:{autocomplete:"off"},model:{value:t.form.durl,callback:function(e){t.$set(t.form,"durl",e)},expression:"form.durl"}})],1),t._v(" "),l("el-form-item",{attrs:{label:"上下游信息","label-width":t.formLabelWidth}},[l("el-select",{attrs:{placeholder:"上/下游"},model:{value:t.form.direction,callback:function(e){t.$set(t.form,"direction",e)},expression:"form.direction"}},[l("el-option",{attrs:{label:"上游",value:"0"}}),t._v(" "),l("el-option",{attrs:{label:"下游",value:"1"}})],1)],1)],1),t._v(" "),l("div",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[l("el-button",{on:{click:function(e){t.dialogFormVisible=!1}}},[t._v("取 消")]),t._v(" "),l("el-button",{attrs:{type:"primary"},on:{click:t.add}},[t._v("确 定")])],1)],1)],1)},staticRenderFns:[]},p=l("VU/8")(d,m,!1,null,null,null).exports;a.default.use(n.a);var b=new n.a({mode:"history",routes:[{path:"/",name:"HelloWorld",component:p}]}),v=l("zL8q"),h=l.n(v),g=(l("tvR6"),l("mtWM")),_=l.n(g),k=l("DWlv"),w=l.n(k);a.default.prototype.http=f,a.default.use(w.a,_.a),a.default.use(h.a),new a.default({el:"#app",router:b,components:{App:r},template:"<App/>"})},PQ22:function(t,e,l){let a=l("gAs1"),o="";t.exports={ip(){try{const t=a.networkInterfaces();t[Object.keys(t)[0]].forEach(function(t){if(""===o&&"IPv4"===t.family&&!t.internal)return o=t.address})}catch(t){console.log(124),o="127.0.0.1"}return o}}},gsu9:function(t,e){},tvR6:function(t,e){}},["NHnr"]);
//# sourceMappingURL=app.57e74790ca3c9954b3f3.js.map
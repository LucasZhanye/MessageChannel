import{j as X,_ as Y,d as L,n as Z,H as F,r as c,p as w,K as ee,x as A,t as j,o as k,A as W,w as G,a2 as te,e as $,N as H,g as m,T as D,a3 as ae,W as le,aB as oe,c as J,f as K,S as se,k as re,an as E,J as U,v as V,aa as ne,D as q,q as ie,a as ce,a0 as ue,aV as ve,y as fe,Q as me,O as de,a9 as pe,U as he}from"./index-04406ec4.js";import{t as be}from"./use-form-item-687dd32a.js";const g=4,ye={vertical:{offset:"offsetHeight",scroll:"scrollTop",scrollSize:"scrollHeight",size:"height",key:"vertical",axis:"Y",client:"clientY",direction:"top"},horizontal:{offset:"offsetWidth",scroll:"scrollLeft",scrollSize:"scrollWidth",size:"width",key:"horizontal",axis:"X",client:"clientX",direction:"left"}},ge=({move:d,size:u,bar:r})=>({[r.size]:u,transform:`translate${r.axis}(${d}%)`}),Q=Symbol("scrollbarContextKey"),we=X({vertical:Boolean,size:String,move:Number,ratio:{type:Number,required:!0},always:Boolean}),Se="Thumb",ze=L({__name:"thumb",props:we,setup(d){const u=d,r=Z(Q),l=F("scrollbar");r||be(Se,"can not inject scrollbar context");const i=c(),v=c(),s=c({}),f=c(!1);let a=!1,h=!1,b=oe?document.onselectstart:null;const t=w(()=>ye[u.vertical?"vertical":"horizontal"]),S=w(()=>ge({size:u.size,move:u.move,bar:t.value})),z=w(()=>i.value[t.value.offset]**2/r.wrapElement[t.value.scrollSize]/u.ratio/v.value[t.value.offset]),T=o=>{var e;if(o.stopPropagation(),o.ctrlKey||[1,2].includes(o.button))return;(e=window.getSelection())==null||e.removeAllRanges(),O(o);const n=o.currentTarget;n&&(s.value[t.value.axis]=n[t.value.offset]-(o[t.value.client]-n.getBoundingClientRect()[t.value.direction]))},P=o=>{if(!v.value||!i.value||!r.wrapElement)return;const e=Math.abs(o.target.getBoundingClientRect()[t.value.direction]-o[t.value.client]),n=v.value[t.value.offset]/2,p=(e-n)*100*z.value/i.value[t.value.offset];r.wrapElement[t.value.scroll]=p*r.wrapElement[t.value.scrollSize]/100},O=o=>{o.stopImmediatePropagation(),a=!0,document.addEventListener("mousemove",N),document.addEventListener("mouseup",y),b=document.onselectstart,document.onselectstart=()=>!1},N=o=>{if(!i.value||!v.value||a===!1)return;const e=s.value[t.value.axis];if(!e)return;const n=(i.value.getBoundingClientRect()[t.value.direction]-o[t.value.client])*-1,p=v.value[t.value.offset]-e,_=(n-p)*100*z.value/i.value[t.value.offset];r.wrapElement[t.value.scroll]=_*r.wrapElement[t.value.scrollSize]/100},y=()=>{a=!1,s.value[t.value.axis]=0,document.removeEventListener("mousemove",N),document.removeEventListener("mouseup",y),C(),h&&(f.value=!1)},M=()=>{h=!1,f.value=!!u.size},x=()=>{h=!0,f.value=a};ee(()=>{C(),document.removeEventListener("mouseup",y)});const C=()=>{document.onselectstart!==b&&(document.onselectstart=b)};return A(j(r,"scrollbarElement"),"mousemove",M),A(j(r,"scrollbarElement"),"mouseleave",x),(o,e)=>(k(),W(le,{name:m(l).b("fade"),persisted:""},{default:G(()=>[te($("div",{ref_key:"instance",ref:i,class:H([m(l).e("bar"),m(l).is(m(t).key)]),onMousedown:P},[$("div",{ref_key:"thumb",ref:v,class:H(m(l).e("thumb")),style:D(m(S)),onMousedown:T},null,38)],34),[[ae,o.always||f.value]])]),_:1},8,["name"]))}});var I=Y(ze,[["__file","/home/runner/work/element-plus/element-plus/packages/components/scrollbar/src/thumb.vue"]]);const _e=X({always:{type:Boolean,default:!0},width:String,height:String,ratioX:{type:Number,default:1},ratioY:{type:Number,default:1}}),Ee=L({__name:"bar",props:_e,setup(d,{expose:u}){const r=d,l=c(0),i=c(0);return u({handleScroll:s=>{if(s){const f=s.offsetHeight-g,a=s.offsetWidth-g;i.value=s.scrollTop*100/f*r.ratioY,l.value=s.scrollLeft*100/a*r.ratioX}}}),(s,f)=>(k(),J(se,null,[K(I,{move:l.value,ratio:s.ratioX,size:s.width,always:s.always},null,8,["move","ratio","size","always"]),K(I,{move:i.value,ratio:s.ratioY,size:s.height,vertical:"",always:s.always},null,8,["move","ratio","size","always"])],64))}});var ke=Y(Ee,[["__file","/home/runner/work/element-plus/element-plus/packages/components/scrollbar/src/bar.vue"]]);const He=X({height:{type:[String,Number],default:""},maxHeight:{type:[String,Number],default:""},native:{type:Boolean,default:!1},wrapStyle:{type:re([String,Object,Array]),default:""},wrapClass:{type:[String,Array],default:""},viewClass:{type:[String,Array],default:""},viewStyle:{type:[String,Array,Object],default:""},noresize:Boolean,tag:{type:String,default:"div"},always:Boolean,minSize:{type:Number,default:20},id:String,role:String,ariaLabel:String,ariaOrientation:{type:String,values:["horizontal","vertical"]}}),Te={scroll:({scrollTop:d,scrollLeft:u})=>[d,u].every(E)},Ne="ElScrollbar",Ce=L({name:Ne}),Re=L({...Ce,props:He,emits:Te,setup(d,{expose:u,emit:r}){const l=d,i=F("scrollbar");let v,s;const f=c(),a=c(),h=c(),b=c("0"),t=c("0"),S=c(),z=c(1),T=c(1),P=w(()=>{const e={};return l.height&&(e.height=U(l.height)),l.maxHeight&&(e.maxHeight=U(l.maxHeight)),[l.wrapStyle,e]}),O=w(()=>[l.wrapClass,i.e("wrap"),{[i.em("wrap","hidden-default")]:!l.native}]),N=w(()=>[i.e("view"),l.viewClass]),y=()=>{var e;a.value&&((e=S.value)==null||e.handleScroll(a.value),r("scroll",{scrollTop:a.value.scrollTop,scrollLeft:a.value.scrollLeft}))};function M(e,n){pe(e)?a.value.scrollTo(e):E(e)&&E(n)&&a.value.scrollTo(e,n)}const x=e=>{E(e)&&(a.value.scrollTop=e)},C=e=>{E(e)&&(a.value.scrollLeft=e)},o=()=>{if(!a.value)return;const e=a.value.offsetHeight-g,n=a.value.offsetWidth-g,p=e**2/a.value.scrollHeight,_=n**2/a.value.scrollWidth,R=Math.max(p,l.minSize),B=Math.max(_,l.minSize);z.value=p/(e-p)/(R/(e-R)),T.value=_/(n-_)/(B/(n-B)),t.value=R+g<e?`${R}px`:"",b.value=B+g<n?`${B}px`:""};return V(()=>l.noresize,e=>{e?(v==null||v(),s==null||s()):({stop:v}=ne(h,o),s=A("resize",o))},{immediate:!0}),V(()=>[l.maxHeight,l.height],()=>{l.native||q(()=>{var e;o(),a.value&&((e=S.value)==null||e.handleScroll(a.value))})}),ie(Q,ce({scrollbarElement:f,wrapElement:a})),ue(()=>{l.native||q(()=>{o()})}),ve(()=>o()),u({wrapRef:a,update:o,scrollTo:M,setScrollTop:x,setScrollLeft:C,handleScroll:y}),(e,n)=>(k(),J("div",{ref_key:"scrollbarRef",ref:f,class:H(m(i).b())},[$("div",{ref_key:"wrapRef",ref:a,class:H(m(O)),style:D(m(P)),onScroll:y},[(k(),W(me(e.tag),{id:e.id,ref_key:"resizeRef",ref:h,class:H(m(N)),style:D(e.viewStyle),role:e.role,"aria-label":e.ariaLabel,"aria-orientation":e.ariaOrientation},{default:G(()=>[fe(e.$slots,"default")]),_:3},8,["id","class","style","role","aria-label","aria-orientation"]))],38),e.native?de("v-if",!0):(k(),W(ke,{key:0,ref_key:"barRef",ref:S,height:t.value,width:b.value,always:e.always,"ratio-x":T.value,"ratio-y":z.value},null,8,["height","width","always","ratio-x","ratio-y"]))],2))}});var Be=Y(Re,[["__file","/home/runner/work/element-plus/element-plus/packages/components/scrollbar/src/scrollbar.vue"]]);const Oe=he(Be);export{Oe as E};
import{p as l,P as c,r as f,bY as v,n,g as a}from"./index-6239247c.js";function y(s){return s==null}const u=s=>{const o=c();return l(()=>{var e,t;return(t=(e=o==null?void 0:o.proxy)==null?void 0:e.$props)==null?void 0:t[s]})},m=Symbol("formContextKey"),b=Symbol("formItemContextKey"),z=(s,o={})=>{const e=f(void 0),t=o.prop?e:u("size"),d=o.global?e:v(),r=o.form?{size:void 0}:n(m,void 0),i=o.formItem?{size:void 0}:n(b,void 0);return l(()=>t.value||a(s)||(i==null?void 0:i.size)||(r==null?void 0:r.size)||d.value||"")},x=s=>{const o=u("disabled"),e=n(m,void 0);return l(()=>o.value||a(s)||(e==null?void 0:e.disabled)||!1)};export{x as a,m as b,b as f,y as i,z as u};

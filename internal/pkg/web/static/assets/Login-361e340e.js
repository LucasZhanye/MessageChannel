import{d as w,r as x,a as m,u as E,b as V,c as h,e as n,f as o,w as r,o as y,g as c,h as F,l as b,i as k,E as B}from"./index-6239247c.js";import{E as C,a as I}from"./el-form-item-f4a42828.js";import{E as L}from"./el-button-8c65751f.js";import{E as N}from"./el-input-16e5a8d1.js";import"./use-form-common-props-2644540a.js";import"./castArray-a30e3067.js";import"./use-form-item-5c7ee7ac.js";const R="/assets/logo-6653ca4f.png",U={class:"login-container"},q={class:"content"},H=n("div",{class:"logo"},[n("img",{src:R})],-1),M={class:"login"},S=n("div",{class:"title"}," Welcome ",-1),J=w({__name:"Login",setup(T){const u=x(),s=m({username:"",password:""}),d=m({username:[{required:!0,trigger:"blur"}],password:[{required:!0,trigger:"blur"}]}),p=E(),_=V(),f=async l=>{l&&await l.validate((e,i)=>{e?p.login(s).then(()=>{_.push({name:"Home"})}).catch(t=>{console.log("login fail:",t)}):B.error("Incorrect information entered")})};return(l,e)=>{const i=N,t=I,g=L,v=C;return y(),h("div",U,[n("div",q,[H,n("div",M,[S,o(v,{class:"loginForm",model:s,rules:d,ref_key:"ruleFormRef",ref:u},{default:r(()=>[o(t,{prop:"username"},{default:r(()=>[o(i,{modelValue:s.username,"onUpdate:modelValue":e[0]||(e[0]=a=>s.username=a),placeholder:"username","prefix-icon":c(F)},null,8,["modelValue","prefix-icon"])]),_:1}),o(t,{prop:"password"},{default:r(()=>[o(i,{modelValue:s.password,"onUpdate:modelValue":e[1]||(e[1]=a=>s.password=a),type:"password",placeholder:"password","prefix-icon":c(b)},null,8,["modelValue","prefix-icon"])]),_:1}),o(t,null,{default:r(()=>[o(g,{type:"primary",class:"submit",onClick:e[2]||(e[2]=a=>f(u.value))},{default:r(()=>[k(" Login ")]),_:1})]),_:1})]),_:1},8,["model","rules"])])])])}}});export{J as default};

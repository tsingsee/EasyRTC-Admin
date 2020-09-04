<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>
<script>
import "@/assets/style/global.scss";
import { userInfo } from "../src/request/modules/userInfo";
export default {
  mounted(){
    userInfo().then(res=>{
          this.$store.dispatch('user/getUserInfo', res)
    })
  }
}
</script>

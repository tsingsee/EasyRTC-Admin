<template>
  <div class="container_header">
    <div class="minNav">
      <i class="iconfont iconmenu" @click="showAside"></i>
    </div>
    <div class="container_logo">
      <span class="product_Name">EasyRTC</span>
      <span class="version_Number">V1.20</span>
    </div>
    <div class="container_body">
      <div class="container_user">
        <span class="user_photo">
          <img src="../../assets/image/head.png" alt />
        </span>
        <span
          class="user_name"
        >{{userInfo.displayName==''?userInfo.name:userInfo.displayName}}</span>
        <span class="user_operation">
          <i class="iconfont iconxia el-dropdown-link" @click="showDropdown"></i>
          <div class="dropdown" v-show="dropdown">
            <ul>
              <li class="iconfont iconzhuxiao" @click="logout">
                <span>注销</span>
              </li>
            </ul>
          </div>
        </span>
      </div>
    </div>
  </div>
</template>
<script>
import Vue from "vue";
import { mapGetters } from "vuex";
import { logout } from "../../request/modules/login";
export default {
  data() {
    return {
      dropdown: false,
    };
  },
  mounted() {
  },
  computed: {
    ...mapGetters(["userInfo"]),
  },
  methods: {
    // 注销
    logout() {
      logout().then((res) => {
        this.$message({
          message: "注销成功",
          type: "success",
        });
        this.$router.push("/login");
      });
    },
    showDropdown() {
      if (this.dropdown) {
        this.dropdown = false;
      } else {
        this.dropdown = true;
      }
    },
    // 展示菜单栏
    showAside() {
      this.$emit("showAside");
    },
  },
};
</script>
<style lang="scss">
</style>
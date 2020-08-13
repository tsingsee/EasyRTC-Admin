<template>
  <div class="container_asider">
    <div class="menu_box">
      <router-link
        :to="{path:nav.path}"
        v-for="(nav,index) in NavData"
        :key="index"
        :class="[{'action':PathIndex==nav.index},'menu']"
        @click.native="menuAction(nav.index)"
      >
        <i :class="[nav.icon]">
          <span>{{nav.lable}}</span>
        </i>
      </router-link>
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      PathIndex: "",
      NavData: [
        // {
        //   lable: "数据统计",
        //   icon: "iconfont iconbiaopan",
        //   path: "",
        // },
        {
          lable: "会议室",
          index: 1,
          icon: "iconfont iconhuiyi",
          path: "/MeetIndex",
        },
        // {
        //   lable: "视频直播",
        //   icon: "iconfont iconleft-shipinzhibo1",
        //   path: "",
        // },
        {
          index: 3,
          lable: "会议记录",
          icon: "iconfont iconjilu",
          path: "/MeetProc",
        },
        {
          index: 4,
          lable: "会议回看",
          icon: "iconfont iconleft-luxianghuikan",
          path: "/MeetBlack",
        },
        // {
        //   lable: "用户管理",
        //   icon: "iconfont iconjituan1",
        //   path: "",
        // },
        {
          index: 6,
          lable: "账户信息",
          icon: "iconfont iconcanhuirenyuan",
          path: "/UserInfo",
        },
      ],
    };
  },
  mounted() {
    this.RouteRedirection();
  },
  methods: {
    // 菜单点击
    menuAction(e) {
      this.PathIndex = e;
    },
    RouteRedirection() {
      var path = this.$route.path;
      let paths = [
        { path: "/MeetIndex", type: 1 },
        { path: "/editMeet", type: 1 },
        { path: "/MeetProc", type: 3 },
        { path: "/MeetBlack", type: 4 },
        { path: "/UserInfo", type: 6 },
      ];
      for (const i in paths) {
        if (path.indexOf(paths[i].path) != -1) {
          this.PathIndex = paths[i].type;
          return;
        }
      }
    },
  },
};
</script>
<style lang="scss">
</style>
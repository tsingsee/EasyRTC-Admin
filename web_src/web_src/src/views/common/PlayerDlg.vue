<template>
  <div class="container_PlayerDlg">
    <el-dialog :visible.sync="show" width="100%" @close="onHide">
      <div class="container_player" @mouseleave="mouseleave" @mouseover="mouseover">
        <video ref="Video" :src="videoUrl" controls autoplay="autoplay"></video>
        <span class="container_clone" v-show="Clone" @click="onClone">
          <i class="iconfont iconclose1"></i>
        </span>
      </div>
    </el-dialog>
  </div>
</template>

<script>
export default {
  props: ["onShow", "videoUrl"],
  watch: {
    onShow: function (val) {
      this.show = val;
    },
  },
  data() {
    return {
      show: false,
      Clone: false,
    };
  },
  mounted() {},
  methods: {
    onHide() {
      this.$refs.Video.pause();
      this.$emit("onHide");
    },
    onClone() {
      this.show = false;
    },
    mouseleave() {
      console.log("离开");
      this.oClone = false;
    },
    mouseover() {
      console.log("进来");
      this.Clone = true;
    },
  },
};
</script>
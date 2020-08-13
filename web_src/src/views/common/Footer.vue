<template>
  <div class="container_footer">
    <span style="color:#808080">
      Copyright &copy; 2014-{{ thisYear() }}
      <a
        href="http://www.tsingsee.com/"
        style="color:#2a88d7"
        target="_target"
      >
        <span
          style="width: 70px;height: 16px;position: relative;overflow: hidden;display: inline-block;"
        >
          <i
            class="iconfont iconqingxiLOGO"
            style="font-size: 70px;position: absolute;top: -15px;left: 0;color:#2a88d7"
          ></i>
        </span>
        .com
      </a>
			All rights reserved
    </span> 
  </div>
</template>
<script>
export default {
  data() {
    return {};
  },
  methods: {
    thisYear() {
      let date = new Date();
      return date.getFullYear();
    },
  },
};
</script>

<template>
  <div class="AdminTe">
    <Header class="Header" @showAside="showAside"></Header>
    <div class="container">
      <Aside class="Aside"></Aside>
      <div class="maxBox">
        <Main class="Main"></Main>
        <Footer class="Footer"></Footer>
      </div>
    </div>
  </div>
</template>

<script>
import Header from "./Header";
import Aside from "./Aside";
import Main from "./Main";
import Footer from "./Footer";
export default {
  data() {
    return {
      asideshow: false,
    };
  },
  components: {
    Header,
    Aside,
    Main,
    Footer,
  },
  methods: {
    showAside() {
      let show = $(".Aside").css("display");
      if (show == "none") {
        $(".Aside").css("display", "inline-block");
      } else {
        $(".Aside").css("display", "none");
      }
    },
  },
};
</script>

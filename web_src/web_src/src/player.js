import Vue from 'vue'
import store from './store'
import "@/assets/icon/iconfont.css"
import 'element-ui/lib/theme-chalk/index.css';
import Vuex from 'vuex'
import $ from 'jquery'
import VueClipboards from 'vue-clipboards';
import ElementUI from 'element-ui'
import Moment from 'moment';
import Play  from './views/MeetPLayer/player.vue'
Vue.filter('comverTime',function(data,format){
  return Moment(data).format(format);
});

Vue.config.productionTip = false
Vue.use(ElementUI)
Vue.use(Vuex)
Vue.use(VueClipboards);



new Vue({
  el: '#app',
  store,
  template: `<Play></Play>`,
  components: {
    Play
  }
})
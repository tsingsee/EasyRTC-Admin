<template>
  <div class="container_meetIndex">
    <!-- 头部 -->
    <!-- <div class="meetIndex_header">
      <div class="btn_suc" @click="meetAdd">
        <i class="iconfont iconadd"></i>
        创建
    </div>-->
    <!-- <div class="btn_err">
        <i class="iconfont icondelete"></i>
        删除
    </div>-->
    <!-- <div class="inp_search">
        <input type="text" placeholder="请输入搜索内容" />
        <i class="iconfont iconsearch"></i>
      </div>
    </div>-->
    <!-- 主体 -->
    <div class="meetIndex_body Table">
      <el-table :data="videoList" style="width: 100%" stripe>
        <el-table-column type="selection" width="55"></el-table-column>
        <el-table-column prop="roomName" label="会议室名称"></el-table-column>
        <el-table-column label="开始时间">
          <template slot-scope="scope">
            <div>{{scope.row.ctime | comverTime('YYYY-MM-DD HH:mm:ss')}}</div>
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="录制时长">
          <template slot-scope="scope">
            <div>{{duration(scope.row.duration)}}</div>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="文件大小">
          <template slot-scope="scope">
            <div>{{size(scope.row.size)}}</div>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="130">
          <template slot-scope="scope">
            <span class="suc_Col operation" @click="PlayerDlgShow(scope.row)">
              <i class="iconfont iconbofang1"></i>
            </span>

            <a
              :href="scope.row.downloadUrl"
              :download="downloadName(scope.row.downloadUrl)"
              class="suc_Col operation"
            >
              <i class="iconfont iconxiazai1"></i>
            </a>

            <span class="del_Col operation" @click="deleVideoL(scope.row)">
              <i class="iconfont icondelete"></i>
            </span>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="meetIndex_footer pagination">
      <el-pagination
        background
        :current-page.sync="page"
        :page-size="perPage"
        layout="total, prev, pager, next, jumper"
        :total="total"
      ></el-pagination>
    </div>

    <PlayerDlg :onShow="playerShow" @onHide="playerShow=false" :videoUrl="videoUrl"></PlayerDlg>
  </div>
</template>
<script>
import { getVideoList, deleVideoL } from "../../request/modules/meetback";
import PlayerDlg from "../common/PlayerDlg";
export default {
  data() {
    return {
      total: 0,
      page: 1,
      perPage: 10,
      videoList: [],
      playerShow: false,
      videoUrl: "",
    };
  },
  watch: {
    page: function () {
      this.getVideoList();
    },
  },
  components: {
    PlayerDlg,
  },
  computed: {
    size: function () {
      return function (e) {
        return `${(e / 1024 / 1024).toFixed(2)}M`;
      };
    },
    duration: function () {
      return  function(value) {
			let s = parseInt(value);
			let m = 0;
			let h = 0;
			let result = 0;
			if (s >= 60) {
				m = parseInt(s / 60);
				s = parseInt(s % 60);
				if (m >= 60) {
					h = parseInt(m / 60);
					m = parseInt(m % 60);
				}
			}
			result = (h <= 9 ? '0' + h : h).toString() + ':' + (m <= 9 ? '0' + m : m).toString() + ':' + (s <= 9 ? '0' + s : s).toString();
			return result;
		}

    },
    downloadName: function () {
      return function (e) {
        return e.substring(e.lastIndexOf("/") + 1);
      };
    },
  },
  mounted() {
    this.getVideoList();
  },
  methods: {
    PlayerDlgShow(row) {
      this.playerShow = true;
      this.videoUrl = row.downloadUrl;
    },
    // 获取录像列表
    getVideoList() {
      getVideoList({
        page: this.page - 1,
        perPage: this.perPage,
      }).then((res) => {
        console.log(res);
        this.total = res.count;
        this.videoList = res.items;
      });
    },
    // 删除录像
    deleVideoL(e) {
      this.$confirm(`确认删除吗？`, "提示")
        .then(() => {
          deleVideoL({ id: e.id }).then((res) => {
            this.$message({
              message: "删除成功",
              type: "success",
            });
            this.getVideoList();
          });
        })
        .catch(() => {});
    },
  },
};
</script>



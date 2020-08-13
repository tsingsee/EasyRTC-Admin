<template>
  <div class="container_meetIndex">
    <!-- 头部 -->
    <div class="meetIndex_header">
      <div class="btn_suc" @click="meetAdd">
        <i class="iconfont iconadd"></i>
        创建
      </div>
      <!-- <div class="btn_err">
        <i class="iconfont icondelete"></i>
        删除
      </div>-->
      <!-- <div class="inp_search">
        <input type="text" placeholder="请输入搜索内容" />
        <i class="iconfont iconsearch"></i>
      </div> -->
    </div>

    <!-- 主体 -->
    <div class="meetIndex_body Table">
      <el-table :data="meetLists" style="width: 100%" stripe>
        <el-table-column type="selection" width="55"></el-table-column>
        <el-table-column prop="roomName" label="会议室名称" width="100"></el-table-column>
        <el-table-column prop="roomConfig.subject" label="会议主题" width="100"></el-table-column>
        <el-table-column label="清晰度">
          <template slot-scope="scope">{{resolution(scope.row.roomConfig.resolution)}}</template>
        </el-table-column>
        <el-table-column prop="participantLimits" label="参会人数"></el-table-column>
        <el-table-column prop="roomConfig.lockPassword" label="密码"></el-table-column>
        <el-table-column label="参会链接">
          <template slot-scope="scope">
            <span class="meetLink" @click="meetLink(scope.row)">查看</span>
          </template>
        </el-table-column>
        <el-table-column  label="更新时间" width="220">
          <template slot-scope="scope">
            <div>
              {{scope.row.ctime | comverTime('YYYY-MM-DD HH:mm:ss')}}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="130">
          <template slot-scope="scope">
            <span class="suc_Col operation" @click="startMeet(scope.row)">
              <i class="iconfont iconshexiangtou"></i>
            </span>

            <router-link :to="{path:'/editMeet',query:{id:scope.row.id}}" class="suc_Col operation">
              <i class="iconfont iconedit"></i>
            </router-link>

            <span class="del_Col operation" @click="delMeet(scope.row)">
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
    <el-dialog
      title="会议室链接"
      :visible.sync="dialogVisible"
      width="100%"
      class="meetLink_dialog"
      center
    >
      <el-form ref="linkData" :model="linkData">
        <el-form-item label="主持人链接:" label-width="90px" class="linkInput">
          <el-input v-model="linkData.compereUrl" :readonly="true"></el-input>
          <i
            class="iconfont iconfuzhi copy"
            v-clipboard="linkData.compereUrl"
            @success="handleSuccess"
          ></i>
        </el-form-item>
        <el-form-item label="参会者链接:" label-width="90px" class="linkInput">
          <el-input v-model="linkData.commonUrl" :readonly="true"></el-input>
          <i
            class="iconfont iconfuzhi copy"
            v-clipboard="linkData.commonUrl"
            @success="handleSuccess"
          ></i>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="dialogVisible = false">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
import { getMeetList, delMeet, getMeetToken } from "../../request/modules/meet";
import { userInfo } from "../../request/modules/userInfo";
export default {
  data() {
    return {
      dialogVisible: false,
      total: 0,
      page: 1,
      perPage: 10,
      compereToken: "",
      meetLists: [],
      linkData: {
        compereUrl: "",
        commonUrl: "",
      },
    };
  },
  computed: {
    resolution() {
      return function (e) {
        if (e == 360) {
          return "流畅";
        } else if (e == 480) {
          return "标清";
        } else {
          return "高清";
        }
      };
    },
  },
  watch: {
    page: function () {
      this.getMeetList();
    },
  },
  mounted() {
    this.getUserInfo();
    this.getMeetList();
  },
  methods: {
    // 复制成功提示事件
    handleSuccess() {
      this.$message({
        message: "复制成功",
        type: "success",
      });
    },
    meetLink(row) {
      getMeetToken({
        roomName: row.roomName,
      }).then((res) => {
        this.compereToken = res.token;
        this.linkData = {
          compereUrl: `${location.origin}/admin/meetRoom.html/#/?rn=${row.roomName}&jwt=${res.token}`,
          commonUrl: `${location.origin}/admin/meetRoom.html/#/?rn=${row.roomName}`,
        };
        this.dialogVisible = true;
      });
    },
    // 开始会议
    startMeet(row) {
      getMeetToken({
        roomName: row.roomName,
      }).then((res) => {
        this.compereToken = res.token;
        window.open( `${location.origin}/admin/meetRoom.html/#/?rn=${row.roomName}&jwt=${res.token}`);
      });
    },
    // 更新用户信息
    getUserInfo() {
      userInfo().then((res) => {
        this.$store.dispatch("user/getUserInfo", res);
      });
    },
    // 获取房间列表
    getMeetList() {
      getMeetList({
        page: this.page - 1,
        perPage: this.perPage,
      }).then((res) => {
        this.total = res.count;
        this.meetLists = res.items;
      });
    },
    // 删除房间
    delMeet(e) {
      this.$confirm(`确认删除吗？`, "提示")
        .then(() => {
          delMeet({ id: e.id }).then((res) => {
            this.$message({
              message: "删除成功",
              type: "success",
            });
            this.getMeetList();
          });
        })
        .catch(() => {});
    },
    meetAdd() {
      this.$router.push("/editMeet");
    },
    // 获取会议室主持人token
    getMeetToken(e) {},
  },
};
</script>



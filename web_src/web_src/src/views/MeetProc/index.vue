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
      <el-table :data="recordList" style="width: 100%" stripe>
        <el-table-column type="selection" width="55"></el-table-column>
        <el-table-column prop="roomName" label="会议室名称"></el-table-column>
        <el-table-column prop="maxParticipants" label="最大参会人数"></el-table-column>
        <el-table-column label="开始时间">
          <template slot-scope="scope">
            <div>{{scope.row.ctime | comverTime('YYYY-MM-DD HH:mm:ss')}}</div>
          </template>
        </el-table-column>
        <el-table-column  label="结束时间">
          <template slot-scope="scope">
            <div>{{scope.row.etime | comverTime('YYYY-MM-DD HH:mm:ss')}}</div>
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
  </div>
</template>
<script>
import { getrecordList } from "../../request/modules/meetproc";
export default {
  data() {
    return {
      total: 0,
      page: 1,
      perPage: 10,
      recordList: [],
    };
  },
  watch: {
    page: function () {
      this.getrecordList();
    },
  },
  mounted() {
    this.getrecordList();
  },
  methods: {
    // 获取录像列表
    getrecordList() {
      getrecordList({
        page: this.page - 1,
        perPage: this.perPage,
      }).then((res) => {
        console.log(res);
        this.total = res.count;
        this.recordList = res.items;
      });
    },
  },
};
</script>



<template>
  <div class="container_meetEdit">
    <div class="Form">
      <el-form
        :model="ruleForm"
        status-icon
        :rules="rules"
        ref="ruleForm"
        label-width="100px"
        class="demo-ruleForm"
      >
        <div class="formTitle">
          <span>基本信息</span>
        </div>
        <el-form-item label="会议室名称:" prop="roomName">
          <el-input type="text" v-model="ruleForm.roomName" :readonly='readonly' autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="会议室主题:" prop="roomConfig.subject">
          <el-input type="text" v-model="ruleForm.roomConfig.subject" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="最多人数:" prop="participantLimits">
          <el-input type="text" v-model="ruleForm.participantLimits" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="清晰度:" prop="resolution">
          <el-select v-model="ruleForm.roomConfig.resolution" placeholder="请选择">
            <el-option
              v-for="item in options"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="会议室密码:" prop="lockPassword">
          <el-input v-model="ruleForm.roomConfig.lockPassword"></el-input>
        </el-form-item>
        <div class="formTitle">
          <span>其他设置</span>
        </div>
        <div class="Settings">
          <el-checkbox label="是否提示参会者输入名字" v-model="ruleForm.roomConfig.requireDisplayName"></el-checkbox>
          <el-checkbox label="是否加入会议室时不开启音频" v-model="ruleForm.roomConfig.startWithAudioMuted"></el-checkbox>
          <el-checkbox label="是否加入会议室时不开启视频" v-model="ruleForm.roomConfig.startWithVideoMuted"></el-checkbox>
          <el-checkbox label="是否允许服务器录制" v-model="ruleForm.roomConfig.fileRecordingsEnabled"></el-checkbox>
          <el-checkbox label="是否允许直播" v-model="ruleForm.roomConfig.liveStreamingEnabled"></el-checkbox>
          <el-checkbox label="是否允许匿名用户创建会议" v-model="ruleForm.allowAnonymous"></el-checkbox>
        </div>
        <el-form-item class="submit">
          <el-button @click="$router.go(-1)">取消</el-button>
          <el-button type="primary" @click="submitForm('ruleForm')">提交</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import { addMeet, getMeet, editMeet } from "../../request/modules/meet";
export default {
  data() {
    return {
      add: true,
      readonly:false,
      id: "",
      options: [
        {
          value: 360,
          label: "流畅",
        },
        {
          value: 480,
          label: "标清",
        },
        {
          value: 720,
          label: "高清",
        },
      ],
      ruleForm: {
        roomName: "",
        participantLimits: 10,
        allowAnonymous: false,
        roomConfig: {
          subject: "",
          resolution: 480,
          lockPassword: "",
          requireDisplayName: true,
          startWithAudioMuted: false,
          fileRecordingsEnabled: false,
          liveStreamingEnabled: false,
          startWithVideoMuted: false,
        },
      },
      rules: {
        roomName: [
          { required: true, message: "会议室名称不能为空" },

          {
            pattern: /^\d{6,}$/, //正则
            message: "最少6位纯数字",
          },
        ],
        participantLimits:[
      
          {
            pattern: /^((?!0)\d{1,3}|1000|1)$/, //正则
            message: "人数为1-1000",
          },
        ],
        subject: [{ required: true, message: "会议室主题不能为空" }],
      },
    };
  },
  mounted() {
    this.id = this.$route.query.id;
    if (this.id) {
      this.add = false;
      this.readonly=true
      this.getmeet();
    }
  },
  methods: {
    getmeet() {
      getMeet({ id: parseInt(this.id) }).then((res) => {
        this.ruleForm = res;
      });
    },
    // 表单提交
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          if (this.add) {
            // 将人数改为数字
            this.ruleForm.participantLimits = parseInt(
              this.ruleForm.participantLimits
            );
            addMeet(this.ruleForm).then((res) => {
              console.log(res);
              this.$message({
                message: "添加成功",
                type: "success",
              });
              this.$router.push("/MeetIndex");
            });
          } else {
            // 将人数改为数字
            this.ruleForm.participantLimits = parseInt(
              this.ruleForm.participantLimits
            );
            editMeet(this.ruleForm).then((res) => {
              console.log(res);
              this.$message({
                message: "修改成功",
                type: "success",
              });
              this.$router.push("/MeetIndex");
            });
          }
        } else {
          console.log("error submit!!");
          return false;
        }
      });
    },
  },
};
</script>
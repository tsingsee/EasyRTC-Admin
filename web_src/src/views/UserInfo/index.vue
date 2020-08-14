<template>
  <div class="container_userInfo">
    <div class="Form">
      <el-form
        :model="ruleForm"
        status-icon
        :rules="rules"
        ref="ruleForm"
        label-width="100px"
        class="demo-ruleForm"
        hide-required-asterisk
      >
        <div class="formTitle">
          <span>个人信息</span>
        </div>
        <el-form-item label="姓名:" prop="displayName">
          <el-input
            type="text"
            placeholder="请输入姓名"
            v-model="ruleForm.displayName"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="邮箱:" prop="email">
          <el-input type="text" placeholder="请输入邮箱" v-model="ruleForm.email" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="手机号码:" prop="phone">
          <el-input v-model="ruleForm.phone" placeholder="请输入手机号" autocomplete="off"></el-input>
        </el-form-item>

        <el-form-item label="公司名称:" prop="company">
          <el-input type="text" placeholder="请输入公司名称" v-model="ruleForm.company" autocomplete="off"></el-input>
        </el-form-item>
        <div class="formTitle">
          <span>修改密码</span>
        </div>
        <el-form-item label="旧密码:" prop="password">
          <el-input
            type="password"
            placeholder="请输入旧密码"
            v-model="ruleForm.password"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="新密码:" prop="newpass">
          <el-input
            type="password"
            placeholder="请输入新密码"
            v-model="ruleForm.newpass"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item type="password" label="确认密码:" prop="confirmpass">
          <el-input
            type="password"
            placeholder="请输入新密码"
            v-model="ruleForm.confirmpass"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item class="submit">
          <el-button @click="$router.go(-1)">取消</el-button>
          <el-button type="primary" @click="submitForm('ruleForm')">提交</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import { userInfo, editUserInfo } from "../../request/modules/userInfo";
export default {
  data() {
    var validatePass1 = (rule, value, callback) => {
      if (this.ruleForm.password != "") {
        if (value === "") {
          callback(new Error("密码不能为空"));
        } else {
          callback();
        }
      } else {
        callback();
      }
    };
    var validatePass2 = (rule, value, callback) => {
      if (this.ruleForm.password != "") {
        if (value === "") {
          callback(new Error("密码不能为空"));
        } else if (value !== this.ruleForm.newpass) {
          callback(new Error("两次输入密码不一致!"));
        } else {
          callback();
        }
      } else {
        callback();
      }
    };
    return {
      ruleForm: {
        displayName: "",
        email: "",
        phone: "",
        company: "",
        password: "",
        newpass: "",
        confirmpass: "",
      },
      rules: {
        newpass: [{ validator: validatePass1, trigger: "blur" }],
        confirmpass: [{ validator: validatePass2, trigger: "blur" }],
      },
    };
  },
  mounted() {
    this.getUserInfo();
  },
  methods: {
    getUserInfo() {
      userInfo().then((res) => {
        console.log(res);
        this.ruleForm.displayName = res.displayName;
        this.ruleForm.email = res.email;
        this.ruleForm.phone = res.phone;
        this.ruleForm.company = res.company;
      });
    },
    // 表单提交
    submitForm(formName) {
      console.log("jin来了");
      this.$refs[formName].validate((valid) => {
        if (valid) {
          // 将人数改为数字
          this.ruleForm.participantLimits = parseInt(
            this.ruleForm.participantLimits
          );
          editUserInfo(this.ruleForm).then((res) => {
            this.$message({
              message: "修改成功",
              type: "success",
            });
            this.$store.dispatch("user/getUserInfo", this.ruleForm);
            this.$router.go(0);
          });
        } else {
          console.log("error submit!!");
          return false;
        }
      });
    },
  },
};
</script>
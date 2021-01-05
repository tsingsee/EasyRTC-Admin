<template>
  <div class="container_login">
    <div class="container_body">
      <el-row>
        <el-col :xs="12" :sm="12" :md="12" :lg="12" class="slideshow">
          <div class="carousel">
            <el-carousel trigger="click" height="580px">
              <el-carousel-item v-for="(item,index) in carousel" :key="index">
                <a >
                  <img :src="item.img" />
                </a>
              </el-carousel-item>
            </el-carousel>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :md="12" :lg="12">
          <div class="loginForm">
            <div class="form">
              <div class="title">EasyRTC-SFU登录</div>
              <div class="body">
                <el-form :model="loginForm" ref="loginForm" :rules="rules" :show-message="false">
                  <el-form-item prop="name">
                    <div class="username">
                      <i class="iconfont iconadmin icon"></i>
                      <input
                        type="text"
                        v-model="loginForm.name"
                        :placeholder="UNplace"
                        class="formInput UNplace"
                      />
                    </div>
                  </el-form-item>
                  <el-form-item prop="password">
                    <div class="password">
                      <i class="iconfont iconpassword icon"></i>
                      <input
                        type="password"
                        v-model="loginForm.password"
                        :placeholder="PWplace"
                        class="formInput PWplace"
                      />
                    </div>
                  </el-form-item>
                  <el-form-item prop="captcha_code">
                    <div class="verification">
                      <input
                        type="text"
                        v-model="loginForm.captcha_code"
                        :placeholder="verPlace"
                        class="formInput verPlace"
                      />
                      <span class="verification_img" @click="getCaptchaId" title="点击刷新">
                        <img :src="CaptchaUrl" />
                      </span>
                    </div>
                  </el-form-item>
                  <el-form-item>
                    <div class="setPassword">
                      <span class="rememberPaw">
                        <el-checkbox v-model="single">记住密码</el-checkbox>
                      </span>
                      <span class="forgetPaw" @click="forgetPaw">忘记密码</span>
                    </div>
                  </el-form-item>
                  <el-form-item>
                    <div class="submit" @click="submit('loginForm')">登录</div>
                    <div class="loginTo">
                      <span>没有账号?</span>
                      <router-link to="/signin" class="loginLink">请注册</router-link>
                    </div>
                  </el-form-item>
                </el-form>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
    <div class="container_footer">
      <span style="color:#333333">
        Copyright &copy; {{ thisYear()}}
        <a href="http://www.tsingsee.com/" style="color:#2a88d7" target="_target">
          <span
            style="width: 78px;height: 16px;position: relative;overflow: hidden;display: inline-block;margin-left: -2px;"
          >
            <i
              class="iconfont iconqingxiLOGO"
              style="font-size: 78px;position: absolute;top: -15px;left: 0;color:#2a88d7"
            ></i>
          </span>
        </a>.com All Rights Reserved.
      </span>
    </div>
  </div>
</template>

<script>
import { getCaptchaId, login } from "../../request/modules/login";
export default {
  data() {
    return {
      CaptchaUrl: "",
      rememberPaw: "",
      single: false,
      UNplace: "请输入用户名",
      PWplace: "请输入密码",
      verPlace: "请输入验证码",
      loginForm: {
        name: "",
        password: "",
        captcha_id: "",
        captcha_code: "",
      },
      rules: {
        name: [{ required: true, message: "账号不能为空" }],
        password: [{ required: true, message: "密码不能为空" }],
        captcha_code: [{ required: true, message: "验证码不能为空" }],
      },
      carousel: [
        {
          img: require("../../assets/image/login1.png"),
          path: "",
        },
      ],
    };
  },
  mounted() {
    this.getCaptchaId();
    this.getCookie();
  },
  methods: {
    // 设置cookie方法
    setCookie(c_name, c_pwd, single, exdays) {
      var exdate = new Date();
      exdate.setTime(exdate.getTime() + 24 * 60 * 60 * 1000 * exdays);
      window.document.cookie =
        "userName" + "=" + c_name + ";path=/;expires=" + exdate.toGMTString();
      window.document.cookie =
        "userPwd" + "=" + c_pwd + ";path=/;expires=" + exdate.toGMTString();
      window.document.cookie =
        "single" + "=" + single + ";path=/;expires=" + exdate.toGMTString();
    },
    // 清楚cookie
    clearCookie: function () {
      this.setCookie("", "", -1);
    },
    // 获取cookie
    getCookie() {
      if (document.cookie.length > 0) {
        var arr = document.cookie.split("; ");
        for (var i = 0; i < arr.length; i++) {
          var arr2 = arr[i].split("=");
          if (arr2[0] == "userName") {
            this.loginForm.name = arr2[1];
          } else if (arr2[0] == "userPwd") {
            this.loginForm.password = arr2[1];
          } else if (arr2[0] == "single") {
            if (arr2[1] == "true") {
              this.single = true;
            }
          }
        }
      }
    },
    // 获取验证码照片
    getCaptchaId() {
      getCaptchaId().then((res) => {
        this.loginForm.captcha_id = res.id;
        this.CaptchaUrl = `${location.origin}/admin/captcha/${res.id}.png`;
        console.log(this.CaptchaUrl);
      });
    },
    // 获取本年年份
    thisYear() {
      let date = new Date();
      return date.getFullYear();
    },
    //提交
    submit(formName) {
      let that = this;
      this.$refs[formName].validate((valid, obj) => {
        console.log(obj);
        if (obj.name) {
          let usererr = obj.name[0].message;
          that.UNplace = usererr;
          $(".UNplace")[0].classList.add("err");
        }
        if (obj.password) {
          let pawerr = obj.password[0].message;
          that.PWplace = pawerr;
          $(".PWplace")[0].classList.add("err");
        }
        if (obj.captcha_code) {
          let vererr = obj.captcha_code[0].message;
          that.verPlace = vererr;
          $(".verPlace")[0].classList.add("err");
        }
        if (valid) {
          login(this.loginForm)
            .then((res) => {
              this.$message({
                message: "登录成功",
                type: "success",
              });
              this.clearCookie();
              if (this.single == true) {
                this.setCookie(
                  this.loginForm.name,
                  this.loginForm.password,
                  this.single,
                  7
                );
              } else {
                this.clearCookie();
              }
              this.$router.push("/MeetIndex");
            })
            .catch((res) => {
              this.getCaptchaId();
            });
        }
      });
    },
    // 忘记密码
    forgetPaw() {
      this.$alert("请联系工作人员！", "提示", {
        confirmButtonText: "确定",
        callback: (action) => {},
      });
    },
  },
};
</script>
<template>
  <div class="container_login">
    <div class="container_body">
      <el-row>
        <el-col :xs="12" :sm="12" :md="12" :lg="12" class="slideshow">
          <div class="carousel">
            <el-carousel trigger="click" height="580px">
              <el-carousel-item v-for="(item,index) in carousel" :key="index">
                <a :href="item.path" target="_blank">
                  <img :src="item.img" />
                </a>
              </el-carousel-item>
            </el-carousel>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :md="12" :lg="12">
          <div class="loginForm">
            <div class="form">
              <div class="title">EasyRTC-SFU注册</div>
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
                  <el-form-item prop="verifyPassword">
                    <div class="password">
                      <i class="iconfont iconpassword icon"></i>
                      <input
                        type="password"
                        v-model="loginForm.verifyPassword"
                        :placeholder="verifyPWplace"
                        class="formInput VPplace"
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
                    <div class="submit" @click="submit('loginForm')">注册</div>
                    <div class="loginTo">
                      <span>已有账号?</span>
                      <router-link to="/login" class="loginLink">请登录</router-link>
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
        </a>
        .com
        All rights reserved
      </span>
    </div>
  </div>
</template>

<script>
import { getCaptchaId, sigin } from "../../request/modules/login";
export default {
  data() {
    return {
      rememberPaw: "",
      CaptchaUrl: "",
      UNplace: "请输入用户名",
      PWplace: "请输入密码",
      verifyPWplace: "请输入确认密码",
      verPlace: "请输入验证码",
      loginForm: {
        name: "",
        password: "",
        captcha_id: "",
        verifyPassword: "",
        captcha_code: "",
      },
      rules: {
        name: [{ required: true, message: "账号不能为空" }],
        password: [{ required: true, message: "密码不能为空" }],
        verifyPassword: [{ required: true, message: "密码不能为空" }],
        captcha_code: [{ required: true, message: "验证码不能为空" }],
      },
      carousel: [
        {
          img: require("../../assets/image/login1.png"),
          path: "#",
        },
        {
          img: require("../../assets/image/login1.png"),
          path: "#",
        },
        {
          img: require("../../assets/image/login1.png"),
          path: "#",
        },
      ],
    };
  },
  mounted() {
    this.getCaptchaId();
  },
  methods: {
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
    // 提交
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
        if (obj.verifyPassword) {
          let verperr = obj.verifyPassword[0].message;
          that.verifyPWplace = verperr;
          $(".VPplace")[0].classList.add("err");
        } else {
          if (that.loginForm.verifyPassword != that.loginForm.password) {
            that.loginForm.verifyPassword = "";
            that.verifyPWplace = "两次密码输入不一致，请重新输入";
            $(".VPplace")[0].classList.add("err");
            return;
          }
        }
        if (valid) {
          sigin(this.loginForm)
            .then((res) => {
              this.$message({
                message: "注册成功",
                type: "success",
              });
              this.$router.push("/MeetIndex");
            })
            .catch((res) => {
              this.getCaptchaId();
            });
        } else {
          console.log("没通过");
        }
      });
    },
  },
};
</script>
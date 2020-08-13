package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"tsingsee.com/adminserver/util"

	"github.com/gin-gonic/gin"
	"tsingsee.com/adminserver/app"
)

type PassportServer struct {
	*app.App
}

func NewPassportServer(app *app.App) *PassportServer {
	return &PassportServer{
		App: app,
	}
}

func (s PassportServer) Signup(c *gin.Context) {
	var param struct {
		app.User
		CaptchaId   string `json:"captcha_id,omitempty"`
		CaptchaCode string `json:"captcha_code,omitempty"`
	}
	if c.BindJSON(&param) != nil {
		return
	}

	if !captcha.VerifyString(param.CaptchaId, param.CaptchaCode) {
		c.AbortWithError(http.StatusBadRequest, errors.New("图片验证码错误"))
		return
	}

	var err error
	param.Password, err = util.HashPassword(param.Password)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user := app.User{}
	err = s.DB().Select(app.SqlStar).From(app.UserTableName).
		Where(app.WhereUserName, param.Name).LoadOneContext(c, &user)
	if user.Name == param.Name {
		c.AbortWithError(http.StatusBadRequest, errors.New("用户名已存在！"))
		return
	}

	param.Ctime = time.Now()
	ctx := c.Request.Context()

	_, err = s.DB().InsertInto(app.UserTableName).
		Columns(app.UserNameCol, app.UserPasswordCol, app.CommonCtimeCol).
		Record(&param.User).ExecContext(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	token := s.CreateToken(param.Id)
	param.Password = ""
	c.SetCookie(app.CookieName, token, 0, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{
		"id": param.Id,
	})
}

func (s PassportServer) Login(c *gin.Context) {
	var param struct {
		app.User
		CaptchaId   string `json:"captcha_id,omitempty"`
		CaptchaCode string `json:"captcha_code,omitempty"`
	}
	if c.BindJSON(&param) != nil {
		return
	}

	if !captcha.VerifyString(param.CaptchaId, param.CaptchaCode) {
		c.AbortWithError(http.StatusBadRequest, errors.New("图片验证码错误"))
		return
	}

	errMsg := "用户名或密码错误"
	user := app.User{}
	err := s.DB().Select(app.SqlStar).From(app.UserTableName).
		Where(app.WhereUserName, param.Name).LoadOneContext(c, &user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New(errMsg))
		return
	}

	pass := util.CheckPasswordHash(param.Password, user.Password)

	if !pass {
		c.AbortWithError(http.StatusBadRequest, errors.New(errMsg))
		return
	}
	token := s.CreateToken(user.Id)
	// secure 为 true 则仅允许 ssl 和 https 协议传输 Cookie
	c.SetCookie(app.CookieName, token, 0, "/", "", false, true)
}

func (s PassportServer) Logout(c *gin.Context) {
	c.SetCookie(app.CookieName, "", -1, "/", "", false, true)
}

// 获取账户信息
func (s PassportServer) Info(c *gin.Context) {
	// 获取正在登陆的信息
	uid := c.GetInt64(app.UserID)
	user := app.User{}
	err := s.DB().Select(app.SqlStar).From(app.UserTableName).
		Where(app.WhereCommonId, uid).LoadOneContext(c, &user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// 修改账户信息
func (s PassportServer) Modify(c *gin.Context) {
	var param struct {
		app.User
		NewPass string `json:"newpass,omitempty"`
	}
	if c.BindJSON(&param) != nil {
		return
	}

	uid := c.GetInt64(app.UserID)
	// 检验原密码
	user := app.User{}
	err := s.DB().Select(app.SqlStar).From(app.UserTableName).
		Where(app.WhereCommonId, uid).LoadOneContext(c, &user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	pass := util.CheckPasswordHash(param.Password, user.Password)
	if !pass {
		c.AbortWithError(http.StatusBadRequest, errors.New("原密码错误"))
		return
	}

	// 通过，将新密码变为 hash
	newPassHash, err := util.HashPassword(param.NewPass)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 更新到数据库中
	_, err = s.DB().Update(app.UserTableName).
		Set(app.UserPasswordCol, newPassHash).
		Set(app.UserDisNameCol, param.DisplayName).
		Set(app.UserCompanyCol, param.Company).
		Set(app.UserPhoneCol, param.Phone).
		Set(app.UserEmailCol, param.Email).
		Where(app.WhereCommonId, uid).
		ExecContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

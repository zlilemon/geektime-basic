package web

import (
	"fmt"
	"geektime-basic/webook/internal/domain"
	"geektime-basic/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

const (
	userIdKey            = "userId"
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	svc              *service.UserService
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:              svc,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (c *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")

	ug.POST("/signup", c.SignUp)
	ug.POST("/login", c.Login)
	ug.POST("/loginJwt", c.LoginJWT)
	ug.POST("/edit", c.Edit)
	ug.POST("/profile", c.Profile)
	ug.POST("/profileJwt", c.ProfileJWT)

}

func (c *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq

	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := c.emailRegexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !isEmail {
		ctx.String(http.StatusOK, "邮箱不正确")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	isPassword, err := c.passwordRegexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含数字、特殊字符，并且长度不能小于8位")
		return
	}

	err = c.svc.SingUp(ctx.Request.Context(),
		domain.User{Email: req.Email, Password: req.ConfirmPassword})

	if err == service.ErrUserDuplicateEmail {

	}

	if err != nil {
		ctx.String(http.StatusOK, "服务异常，注册失败")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}

func (c *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := c.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}

	// 步骤2
	// 这里登录成功了，设置 session
	sess := sessions.Default(ctx)

	// 可以随便设置值，要放在session里面的值
	sess.Set("userId", user.Id)
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")

	return
}

func (c *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := c.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 在这里设置 JWT token
	token := jwt.New(jwt.SigningMethodHS512)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	ctx.Header("x-jwt-token", tokenStr)
	fmt.Println(token)
	fmt.Println(tokenStr)
	fmt.Println(user)

	ctx.String(http.StatusOK, "登录成功")
	return
}

func (c *UserHandler) Edit(ctx *gin.Context) {

}

func (c *UserHandler) Profile(ctx *gin.Context) {
	type Profile struct {
		Email string
	}

	sess := sessions.Default(ctx)
	id := sess.Get(userIdKey).(int64)
	u, err := c.svc.Profile(ctx, id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email: u.Email,
	})
}

func (c *UserHandler) ProfileJWT(ctx *gin.Context) {
	claimTmp, _ := ctx.Get("claims")

	claims, ok := claimTmp.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	println(claims.Uid)
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

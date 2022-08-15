package user

import (
	"net/http"
	"os"
	"shopping/config"
	"shopping/domain/user"
	"shopping/utils/api_helper"
	jwtHelper "shopping/utils/jwt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	userService *user.Service
	appConfig   *config.Configuration
}

// 实例化
func NewUserController(service *user.Service, appConfig *config.Configuration) *Controller {
	return &Controller{
		userService: service,
		appConfig:   appConfig,
	}
}

// CreateUser godoc
// @Summary 根据给定的用户名和密码创建用户
// @Tags Auth
// @Accept json
// @Produce json
// @Param CreateUserRequest body CreateUserRequest true "user information"
// @Success 201 {object} CreateUserResponse
// @Failure 400  {object} api_helper.ErrorResponse
// @Router /user [post]
func (c *Controller) CreateUser(g *gin.Context) {
	var req CreateUserRequest                  //createuserrequest:username,password,password2
	if err := g.ShouldBind(&req); err != nil { //错误处理
		api_helper.HandleError(g, api_helper.ErrInvalidBody)
		return
	}
	newUser := user.NewUser(req.Username, req.Password, req.Password2) //实例化user
	err := c.userService.Create(newUser)                               //调用userservicec层的create
	if err != nil {                                                    //错误处理
		api_helper.HandleError(g, err)
		return
	}

	g.JSON(
		http.StatusCreated, CreateUserResponse{
			Username: req.Username,
		})
}

// Login godoc
// @Summary 根据用户名和密码登录
// @Tags Auth
// @Accept json
// @Produce json
// @Param LoginRequest body LoginRequest true "user information"
// @Success 200 {object} LoginResponse
// @Failure 400  {object} api_helper.ErrorResponse
// @Router /user/login [post]
func (c *Controller) Login(g *gin.Context) {
	var req LoginRequest //登录请求
	if err := g.ShouldBind(&req); err != nil {
		api_helper.HandleError(g, api_helper.ErrInvalidBody)

	}
	currentUser, err := c.userService.GetUser(req.Username, req.Password)
	if err != nil { //错误处理
		api_helper.HandleError(g, err)
		return
	}
	decodedClaims := jwtHelper.VerifyToken(currentUser.Token, c.appConfig.SecretKey)
	if decodedClaims == nil {
		jwtClaims := jwt.NewWithClaims( //返回一个token
			jwt.SigningMethodHS256, jwt.MapClaims{
				"userId":   strconv.FormatInt(int64(currentUser.ID), 10),
				"username": currentUser.Username,
				"iat":      time.Now().Unix(),
				"iss":      os.Getenv("ENV"), //签发人
				"exp": time.Now().Add(
					24 *
						time.Hour).Unix(), //过期时间
				"isAdmin": currentUser.IsAdmin,
			})
		token := jwtHelper.GenerateToken(jwtClaims, c.appConfig.SecretKey) //得到token
		currentUser.Token = token                                          //赋值到user中
		err = c.userService.UpdateUser(&currentUser)                       //更新user
		if err != nil {
			api_helper.HandleError(g, err)
			return
		}
	}

	g.JSON(
		http.StatusOK, LoginResponse{Username: currentUser.Username, UserId: currentUser.ID, Token: currentUser.Token}) //登录响应
}

// 验证token
func (c *Controller) VerifyToken(g *gin.Context) {
	token := g.GetHeader("Authorization")
	decodedClaims := jwtHelper.VerifyToken(token, c.appConfig.SecretKey)

	g.JSON(http.StatusOK, decodedClaims)

}

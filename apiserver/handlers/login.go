package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gogo-cmdb/apiserver/forms"
	"gogo-cmdb/apiserver/models"
	"gogo-cmdb/apiserver/utils"
	"time"
)

func Login(c *gin.Context) {
	var ulf forms.LoginForm
	err := c.ShouldBindJSON(&ulf)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "非法登陆参数")
		return
	}
	userDetails, err := models.DefalutUserManager.GetUserDetailByName(ulf.UserName)
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "用户名或密码错误")
		return
	}

	if userDetails.UserStatus == utils.Deleted || userDetails.UserStatus == utils.Disabled {
		utils.Logger.Info("invalid user")
		utils.BadResponse(c, "无效用户")
		return
	}

	salt, _ := utils.SplitMd5SaltPass(userDetails.Password)

	if utils.Md5SaltPass(ulf.Password, salt) != userDetails.Password {
		utils.Logger.Info("username or password is incorrect")
		utils.BadResponse(c, "用户名或密码错误")
		return
	}

	utils.GlobalConfig.SetDefault("jwt.max_exp", 3600)
	expireAt := utils.GlobalConfig.GetInt("jwt.max_exp")
	standardClaims := jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: time.Now().Add(time.Second * time.Duration(expireAt)).Unix(),
	}

	tokenStr, err := utils.GenJwtAuthToken(
		&utils.AuthCustomClaims{
			UserAuthInfo: utils.UserAuthInfo{
				UserId:     userDetails.User.ID,
				UserName:   userDetails.Name,
				UserType:   userDetails.UserType,
				UserStatus: userDetails.UserStatus,
				Birthday:   userDetails.Birthday,
				Tel:        userDetails.Tel,
				Email:      userDetails.Email,
				Addr:       userDetails.Addr,
				Remark:     userDetails.Remark,
				Avatar:     userDetails.Avatar,
			},
			StandardClaims: standardClaims,
		})
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "获取token失败")
	}
	maxAge := utils.GlobalConfig.GetInt("server.token_exp")
	domain := utils.GlobalConfig.GetString("server.domain")

	c.SetCookie("authToken", tokenStr, maxAge, "/", domain, false, false)
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "登陆成功",
		"data": gin.H{
			"user":  userDetails,
			"token": tokenStr,
		},
	})
	utils.Logger.Info("login success")
}

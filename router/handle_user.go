package router

import (
	"gin-web-api-ws-mqtt-frame/db"
	"gin-web-api-ws-mqtt-frame/middleware"
	"gin-web-api-ws-mqtt-frame/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/utils"
)

var UserController = &userController{}

type userController struct {
}

// 用户登录
func (uc *userController) Login(c *gin.Context) {
	appG := appGin(c)

	//处理请求体中的json格式
	var ReqData model.LoginReq
	err := ParseJSON(c, &ReqData)
	if err != nil {

		return
	}

	//拿到请求体中的 phone 和 password
	phone := ReqData.Phone
	password := ReqData.Passwd

	//通过查找手机号返回该用户的 id ，因为后面的token加密需要 id
	data := db.User{}
	db.DB.Table("user").Where("phone = ? ", phone).First(&data)

	//token加密
	id := data.ID
	token, _ := middleware.JwtSign(id, phone)

	//获取头像和用户名
	name := data.Name
	url := data.Url
	phones := data.Phone

	//密码解密
	JPswd := middleware.ValidPassWord(password, "utek", data.Passwd)

	if data.Phone == "" || !JPswd {
		appG.RespErrData("用户名或密码不正确！！！")
	} else {
		resp := model.LoginRes{
			Name:  name,
			Url:   url,
			Phone: phones,
			Token: token,
		}
		appG.RespData(resp)
	}
}

// 更改密码
func (uc *userController) ChangePwd(c *gin.Context) {
	appG := appGin(c)
	data := db.User{}
	//从token中拿到用户id
	uid, _ := c.Get("uid")
	id := utils.ToString(uid)

	//处理json请求
	ReqData := model.ChangePwdReq{}
	err := ParseJSON(c, &ReqData)
	if err != nil {
		return
	}

	oldpasswd := ReqData.OldPasswd
	newpasswd := ReqData.NewPasswd

	//查询数据库中是否有该用户
	db.DB.Table("user").Where("id = ?", id).First(&data)

	//密码解密，判断密码是否正确
	Pd := middleware.ValidPassWord(oldpasswd, "utek", data.Passwd)

	//如果Pd为true，密码正确，把新密码加密后存进数据库
	if Pd {
		passwd := middleware.MakePassWord(newpasswd, "utek")
		db.DB.Table("user").Where("id = ?", id).Update("passwd", passwd)
		appG.RespOK()
	} else {
		appG.RespErrData("修改失败")
	}
}

// 拿所有用户，用作测试，可删
func (uc *userController) GetList(c *gin.Context) {
	appG := appGin(c)
	data := make([]*db.User, 10)
	db.DB.Table("user").Find(&data)
	appG.RespData(data)
}

// 创建用户，方便创建用户，可删
func (uc *userController) Create(c *gin.Context) {
	appG := appGin(c)
	data := db.User{}
	user := model.CreateUserReq{}

	user.Name = c.PostForm("name")
	user.Phone = c.PostForm("phone")
	user.Passwd = c.PostForm("password")
	repassword := c.PostForm("repassword")

	//根据手机号拿到用户的信息
	db.DB.Table("user").Where("phone = ?", user.Phone).First(&data)

	//密码加密
	password := middleware.MakePassWord(user.Passwd, "utek")

	if user.Name == "" {
		appG.RespErrData("用户名不能为空")
	} else if user.Passwd != repassword {
		appG.RespErrData("两次密码不一致")
	} else if data.Phone != "" {
		appG.RespErrData("手机号已存在")
	} else {
		//创建用户
		user.Passwd = password
		db.DB.Table("user").Create(&user)
		appG.RespOK()
	}
}

// 删除用户，用作测试，可删除
func (uc *userController) DeleteUser(c *gin.Context) {
	appG := appGin(c)
	id := c.PostForm("id")
	data := db.User{}
	db.DB.Model("user").Where("id=?", id).First(&data)

	//这里的data.ID==0，是因为数据库中找不到该用户的话，id的类型为unit，所以为0
	if data.ID == 0 {
		appG.RespErrData("找不到该用户")
	} else {
		user := model.DeleteUserReq{}
		db.DB.Table("user").Where("id=?", id).Delete(&user)
		appG.RespOKData("删除成功")
	}

}

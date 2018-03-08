package routers

import (
	"github.com/astaxie/beego"
	"lovehouse/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	//请求地域信息
	beego.Router("/api/v1.0/areas", &controllers.AreaController{}, "get:GetAreaInfo")
	//session

	beego.Router("/api/v1.0/session", &controllers.SessionController{}, "get:GetSessionName;delete:DelSessionName")

	beego.Router("/api/v1.0/sessions", &controllers.UserController{}, "post:Login")

	beego.Router("/api/v1.0/user/avatar", &controllers.UserController{}, "post:UploadAvatar")
	//house/index  房屋首页列表
	beego.Router("/api/v1.0/houses/index", &controllers.HousesIndexController{}, "get:GetHousesIndex")
	// api/v1.0/users [post]
	beego.Router("/api/v1.0/users", &controllers.UserController{}, "post:Reg")

	//更改用户名
	beego.Router("api/v1.0/user/name", &controllers.ModifyUserNameController{}, "put:ModfileUserName")

	//请求用户信息
	beego.Router("api/v1.0/user", &controllers.ShowUserInfoController{}, "get:GetUserInfo")

	//请求实名认证信息
	beego.Router("api/v1.0/user/auth", &controllers.AuthController{}, "get:GetAuthInfo")

	//更新实名认证信息
	beego.Router("api/v1.0/user/auth", &controllers.AuthController{}, "post:SetAuthInfo")
}

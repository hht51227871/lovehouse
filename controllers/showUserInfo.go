package controllers

import (
	_ "encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"lovehouse/models"
)

type ShowUserInfoController struct {
	beego.Controller
}

func (this *ShowUserInfoController) ReturnData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *ShowUserInfoController) GetUserInfo() {
	//通过seesion得到用户user_id
	user_id := this.GetSession("user_id")

	beego.Info("GetUserInfo get id = ", user_id)
	//查询数据库得到用户信息
	var user models.User
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	qs.Filter("Id", user_id.(int)).One(&user)

	//用户信息发给前端
	resp := make(map[string]interface{})

	defer this.ReturnData(resp)

	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	resp["data"] = user

	return

}

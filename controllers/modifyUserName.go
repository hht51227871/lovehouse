package controllers

import (
	"encoding/json"
	//      "fmt"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"lovehouse/models"
)

type ModifyUserNameController struct {
	beego.Controller
}

func (this *ModifyUserNameController) ReturnData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *ModifyUserNameController) ModfileUserName() {
	resp := make(map[string]interface{})

	defer this.ReturnData(resp)

	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	//得到前端发来的数据
	var request = make(map[string]interface{})
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)

	beego.Info("name = ", request["name"])

	//判断数据的合法性
	if request["name"] == "" {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}

	//通过session得到user.Mobile
	user_id := this.GetSession("user_id")
	user := models.User{Id: user_id.(int), Name: request["name"].(string)}

	//更新数据库
	o := orm.NewOrm()
	if _, err := o.Update(&user, "name"); err != nil {
		beego.Info("modify name err", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	//给前端返回成功
	this.SetSession("name", user.Name)
	this.SetSession("user_id", user.Id)
	this.SetSession("mobile", user.Mobile)

	return

}

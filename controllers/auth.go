package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"lovehouse/models"
)

type AuthController struct {
	beego.Controller
}

func (this *AuthController) ReturnData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *AuthController) GetAuthInfo() {

	beego.Info("auth get  suss")
	resp := make(map[string]interface{})

	defer this.ReturnData(resp)

	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	//从sessio中获取user_id
	user_id := this.GetSession("user_id")
	var user models.User

	//根据user_id查询数据库
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	err := qs.Filter("Id", user_id.(int)).One(&user)
	if err != nil {
		//查询数据库失败
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	//查询成功
	beego.Info("=====auth get name , ", user.Name)
	//返回信息给前端
	resp["data"] = user
	return
}

func (this *AuthController) SetAuthInfo() {
	beego.Info("auth get  suss")
	resp := make(map[string]interface{})

	defer this.ReturnData(resp)

	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	//获取前端发来的信息
	var request = make(map[string]interface{})
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)

	beego.Info("id_card = ", request["id_card"])
	beego.Info("real_name = ", request["real_name"])

	//从sessio中获取user_id
	user_id := this.GetSession("user_id")

	//根据id更新数据库
	o := orm.NewOrm()
	user := models.User{Id: user_id.(int), Real_name: request["real_name"].(string), Id_card: request["id_card"].(string)}

	if _, err := o.Update(&user, "real_name", "id_card"); err != nil {
		beego.Info("user updata err, ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return

	}

	this.SetSession("user_id", user.Id)
	this.SetSession("name", user.Name)

	return
	//返回成功信息给前端
}

package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"lovehouse/models"
	"path"
)

type UserController struct {
	beego.Controller
}

//将封装好的返回结构 变成json返回给前段
func (this *UserController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

//处理上传头像的业务
func (this *UserController) UploadAvatar() {

	//返回给前端的map结构体
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	//得到文件二进制数据
	file, header, err := this.GetFile("avatar")
	if err != nil {
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}

	fileBuffer := make([]byte, header.Size)
	if _, err := file.Read(fileBuffer); err != nil {
		resp["errno"] = models.RECODE_IOERR
		resp["errmsg"] = models.RecodeText(models.RECODE_IOERR)
		return
	}
	//beego.Info("fileBuffer get succ----------------", fileBuffer)

	suffix := path.Ext(header.Filename) // home.jpg.rmvb--->  .rmvb

	//将文件的二进制数据上传到fastdfs中 ---> fileid
	//fileBuffer--->fastdfs  ====>fileid
	groupName, fileId, err := models.FDFSUploadByBuffer(fileBuffer, suffix[1:]) //"rmvb"
	if err != nil {
		resp["errno"] = models.RECODE_IOERR
		resp["errmsg"] = models.RecodeText(models.RECODE_IOERR)
		beego.Info("upload file to fastdfs error err = ", err)
		return
	}

	beego.Info("fdfs upload succ groupname = ", groupName, "  fileid = ", fileId)

	//fileid ---> user 表里avatar_ur字段中
	//可以从seession中获得user.Id
	user_id := this.GetSession("user_id")
	user := models.User{Id: user_id.(int), Avatar_url: fileId}

	//数据库的操作，
	o := orm.NewOrm()
	if _, err := o.Update(&user, "avatar_url"); err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	//将fileid拼接成一个完整的url路径
	avatar_url := "http://127.0.0.1:6379/" + fileId

	//安装协议做出json返回给前端

	url_map := make(map[string]interface{})
	url_map["avatar_url"] = avatar_url
	resp["data"] = url_map

	return

}

//  /api/v1.0/users [post]
/*

{
        mobile: "123",
        password: "123",
        sms_code: "123"
}
*/

func (this *UserController) Reg() {
	beego.Info("==========/api/v1.0/users post succ!!!=========")

	//返回给前端的map结构体
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	var regRequestMap = make(map[string]interface{})

	//1 得到客户端请求的json数据 post数据
	json.Unmarshal(this.Ctx.Input.RequestBody, &regRequestMap)

	beego.Info("mobile = ", regRequestMap["mobile"])
	beego.Info("password = ", regRequestMap["password"])
	beego.Info("sms_code = ", regRequestMap["sms_code"])

	//2 判断数据的合法性
	if regRequestMap["mobile"] == "" || regRequestMap["password"] == "" || regRequestMap["sms_code"] == "" {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}

	//3 将数据存入到mysql数据库 user
	user := models.User{}
	user.Mobile = regRequestMap["mobile"].(string)
	//应该将password进行md5，SHA246,SHA1
	user.Password_hash = regRequestMap["password"].(string)
	user.Name = regRequestMap["mobile"].(string)

	o := orm.NewOrm()

	id, err := o.Insert(&user)
	if err != nil {
		beego.Info("insert error = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	beego.Info("reg succ !!! user id = ", id)

	//4 将当前的用户的信息存储到session中
	this.SetSession("name", user.Mobile)
	this.SetSession("user_id", id)
	this.SetSession("mobile", user.Mobile)

	return
}

//login

func (this *UserController) Login() {
	beego.Info("!!!!!!!!!/api/v1.0/session login succ!!!!!!!!!!!")

	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	var loginRequestMap = make(map[string]interface{})

	json.Unmarshal(this.Ctx.Input.RequestBody, &loginRequestMap)
	beego.Info("mobile = ", loginRequestMap["mobile"])
	beego.Info("password = ", loginRequestMap["password"])

	if loginRequestMap["mobile"] == "" || loginRequestMap["password"] == "" {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}

	var user models.User

	o := orm.NewOrm()

	qs := o.QueryTable("user")
	if err := qs.Filter("mobile", loginRequestMap["mobile"]).One(&user); err != nil {
		resp["errno"] = models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
		return
	}

	if user.Password_hash != loginRequestMap["password"].(string) {
		resp["errno"] = models.RECODE_PARAMERR
		resp["errmsg"] = models.RecodeText(models.RECODE_PARAMERR)
		return
	}

	beego.Info("login succ!!!")

	this.SetSession("name", user.Mobile)
	this.SetSession("user_id", user.Id)
	this.SetSession("mobile", user.Mobile)

	return
}

func (this *UserController) Get() {

}

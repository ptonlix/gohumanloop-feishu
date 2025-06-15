package sqlite

import (
	"github.com/beego/beego/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/ptonlix/gohumanloop-wework/models"
)

var Mydb orm.Ormer

func init() {

	datapath, err := beego.AppConfig.String("datapath")
	if err != nil {
		logs.Error("数据库路径配置加载失败: datapath")
		panic("数据库路径配置加载失败: datapath 不能为空")
	}

	orm.RegisterDataBase("default", "sqlite3", datapath)
	orm.RegisterModel(new(models.HumanLoop))
	// orm.RunSyncdb("default", true, true) # 第二个参数是控制强制建库

	Mydb = orm.NewOrm()

}

package DBBase

//
//import (
//	_ "database/sql"
//	_ "database/sql/driver"
//	_ "fmt"
//	//_ "github.com/go-sql-driver/mysql"
//	//"github.com/jmoiron/sqlx"
//)
//
//type SqlManagerInterface interface {
//}
//
////对应数据库表
//type im_group_members_relation struct { //用户群关系表
//	User_id             int64 `db:user_id`             //用户ID	主键
//	Group_id            int64 `db:group_id`            //群ID		外键
//	Member_status       int32 `db:member_status`       //默认正常1；-1禁言；踢出-1；待审核-2；-3踢出或退出；-4拒绝加入
//	Group_manage_type   int32 `db:group_mamange_type`  //默认群员0,群管理员1,群主2
//	Group_remind_status int32 `db:group_remind_status` //群提示状态：默认提示1；不提示0
//	Create_time         int64 `db:create_time`         //创建时间
//	Update_time         int64 `db:update_time`         //更新时间
//}
//
//type MysqlManager struct {
//	username     string
//	passwd       string
//	addr         string
//	port         string
//	maxContent   int32
//	timeOutValue int32
//	DBName       string
//	contents     *sqlx.DB
//}
//
////初始化数据库配置
//func NewMysqlManager() *MysqlManager {
//	sqlmanager := MysqlManager{
//		"root",
//		"siqirui",
//		"192.168.50.40",
//		"3306",
//		23,
//		3,
//		"im-group",
//		nil,
//	}
//	err := sqlmanager.connetsql()
//	if err == nil {
//		return nil
//	}
//	sqlmanager.contents = err
//	return &sqlmanager
//}
//
////连接数据库
//func (mysqlmanager *MysqlManager) connetsql() *sqlx.DB {
//	dataSourceName := mysqlmanager.username + ":" + mysqlmanager.passwd + "@" + "(" + mysqlmanager.addr + ":" + mysqlmanager.port + ")" + "/" + mysqlmanager.DBName
//	db, error := sqlx.Open("mysql", dataSourceName)
//
//	if error != nil {
//		db.Close()
//		return nil
//	}
//	return db
//}

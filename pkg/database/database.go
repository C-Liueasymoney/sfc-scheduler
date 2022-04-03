package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
	"strings"
)

const (
	userName = "root"
	password = "7856915061"
	ip       = "192.168.27.2"
	port     = "3306"
	dbName   = "sfc"
)

var DB *sql.DB

type Test struct {
	id     int
	sfc_id int
	src    string
	dst    string
}

func InitDB() {
	klog.V(3).Infof("*************** initDB ***************")
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail")
		klog.V(3).Infof("open database fail")
		return
	}
	fmt.Println("connect success")
	klog.V(3).Infof("database connect success")
}

func Query() {
	klog.V(3).Infof("*************** query DB ***************")
	var test Test
	rows, e := DB.Query("select * from test")
	if e != nil {
		klog.V(3).Infof("query error")
	}
	for rows.Next() {
		rows.Scan(&test.src, &test.dst, &test.sfc_id, &test.id)
		//if e == nil {
		klog.V(3).Infoln(json.Marshal(test))
		fmt.Println(test)
		//}
	}
	rows.Close()
}

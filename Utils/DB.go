package Utils

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB
var rdb *redis.Client

//init 初始化数据库
func init() {
	var err error
	dsn := `root:123456@tcp(localhost:3306)/ShippingTraceability`
	//dsn := `BackEnd:20010307ck.@tcp(101.42.99.73:3306)/ShippingTraceability`
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Println("[Mysql]连接失败")
		panic(err)
	}
	fmt.Println("[Mysql]连接成功")
	//连接Redis数据库
	rdb = redis.NewClient(&redis.Options{
		//Addr:     "101.42.99.73:6379",
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		log.Println("[Redis]连接失败")
		panic(err)
	}
	fmt.Println("[Redis]链接成功")
}

//RDB Redis缓存
func RDB() *redis.Client {
	return rdb
}

//DB Mysql数据库
func DB() *sql.DB {
	return db
}

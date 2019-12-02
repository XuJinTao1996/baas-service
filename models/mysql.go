package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/prometheus/common/log"
	"rds-front/pkg/setting"
)

var db *gorm.DB

type Mysql struct {
	Model
	DBInstanceName string `json:"db_instance_name"`
	DBUser         string `gorm:"default:'root'" json:"db_user"`
	DBPasswd       string `json:"db_passwd"`
	DBPort         int    `gorm:"default:3306"json:"db_port"`
	DBServiceUrl   string `json:"db_service_url"`
	DeploymentMode string `json:"deployment_mode"`
}

// 初始化创建数据表
func init() {
	var err error
	db_file := setting.DBFile
	db_type := setting.DBType
	table_prefix := setting.TablePrefix
	db, err = gorm.Open(db_type, db_file)
	if err != nil {
		log.Fatal("Failed to open the db file: %v", err)

		gorm.DefaultTableNameHandler = func(baas *gorm.DB, defaultTableName string) string {
			return table_prefix + defaultTableName
		}
	}
	db.AutoMigrate(&Mysql{})
}

// 获取所有 mysql 实例的函数
func GetAllMysqlInstance() ([]Mysql, int) {
	var mysqls []Mysql
	var count int
	db.Find(&mysqls)
	db.Model(&Mysql{}).Count(&count)
	return mysqls, count
}

// 判断是否存在该 mysql 实例
func ExistMysqlInstance(db_instance_name string) bool {
	var mysql Mysql
	db.Select("id").Where("db_instance_name=?", db_instance_name).First(&mysql)
	if mysql.ID > 0 {
		return true
	}

	return false
}

// 创建数据库实例
func AddTag(db_instance_name string, db_user string, db_passwd string, db_port int, deployment_mode string) bool {
	db.Create(&Mysql{
		DBInstanceName: db_instance_name,
		DBUser:         db_user,
		DBPasswd:       db_passwd,
		DBPort:         db_port,
		DeploymentMode: deployment_mode,
	})
	return true
}

func CloseDB() {
	defer db.Close()
}

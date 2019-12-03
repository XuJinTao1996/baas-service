package models

import (
	"baas-service/pkg/setting"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

var db *gorm.DB

// base model
type Model struct {
	ID        int        `gorm:"primary_key" json:"id"`
	UID       int        `json:"uid"`
	CreatedAt time.Time  `gorm:"column:created_at;"`
	UpdatedAt time.Time  `gorm:"column:updated_at;"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`
}

type MysqlCluster struct {
	Model
	Namespace   string `gorm:"default:'mysql-operator'" json:"namespace"`
	ClusterName string `json:"db_instance_name"`
	Member      int    `json:"member"`
	User        string `gorm:"default:'root'" json:"db_user"`
	Passwd      string `json:"db_passwd"`
	Port        int    `gorm:"default:3306"json:"db_port"`
	ServiceUrl  string `json:"db_service_url"`
	MultiMaster bool   `json:"multiMaster"`
	Version     string `json:"version"`
	StorageType string `json:"storage_type"`
	VolumeSize  int    `json:"volume_size"`
}

type MysqlConfig struct {
	Name string `json:"name"`
	Data string `json:"data"`
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
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return table_prefix + defaultTableName
	}
	db.AutoMigrate(&MysqlCluster{})
}

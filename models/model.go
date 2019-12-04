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
	Namespace   string `form:"namespace,default=default" binding:"required"`
	ClusterName string `form:"cluster_name" binding:"required"`
	Member      int    `form:"member,default=1"`
	User        string `form:"db_user,default=root"`
	Password    string `form:"password" binding:"required"`
	Port        int    `form:"db_port,default=3306"`
	ServiceUrl  string `form:"db_service_url"`
	MultiMaster bool   `form:"multi_master,default=true"`
	Version     string `form:"version,default=8.0.12"`
	StorageType string `form:"storage_type" binding:"required"`
	VolumeSize  string `form:"volume_size,default=1Gi"`
	Status      string `form:"status,default=NotReady"`
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

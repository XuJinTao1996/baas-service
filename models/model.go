package models

import (
	"baas-service/pkg/setting"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

var db *gorm.DB

type Model struct {
	ID        int        `gorm:"primary_key" json:"id"`
	UID       int        `json:"uid"`
	CreatedAt time.Time  `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;" json:"deleted_at"`
}

type MysqlCluster struct {
	Model
	Namespace                   string `form:"namespace,default=default" json:"namespace" binding:"required"`
	ClusterName                 string `form:"cluster_name" json:"cluster_name" binding:"required"`
	Member                      int    `form:"member,default=1" json:"member"`
	User                        string `form:"db_user,default=root" json:"user"`
	Password                    string `form:"password" json:"password" binding:"required"`
	Port                        int    `form:"db_port,default=3306" json:"port"`
	Host                        string `form:"host" json:"host"`
	MultiMaster                 bool   `form:"multi_master,default=false" json:"multi_master"`
	Version                     string `form:"version,default=8.0.12" json:"version"`
	StorageType                 string `form:"storage_type" binding:"required" json:"storage_type"`
	VolumeSize                  string `form:"volume_size,default=1Gi" json:"volume_size"`
	DefaultAuthenticationPlugin string `form:"default_authentication_plugin,default=mysql_native_password" json:"default_authentication_plugin"`
	CPU                         string `form:"cpu,default=500m" json:"cpu"`
	Memory                      string `form:"memory,default=1Gi" json:"memory"`
	MaxConnections              int    `form:"max_connections,default=10" json:"max_connections"`
	Status                      string `form:"status,default=NotReady" json:"status"`
}

type ZookeeperCluster struct {
	Model
	Namespace   string `form:"namespace,default=zookeeper-operator" json:"namespace" binding:"required"`
	ClusterName string `form:"cluster_name" json:"cluster_name" binding:"required"`
	Size        int    `form:"size,default=1" json:"size"`
	Status      string `form:"status,default=NotReady" json:"status"`
}

type KafkaCluster struct {
	Model
	Namespace          string `form:"namespace,default=kafka-operator" json:"namespace" binding:"required"`
	ClusterName        string `form:"cluster_name" json:"cluster_name" binding:"required"`
	ZookeeperClusterID int    `form:"zookeeper_cluster_id" json:"zookeeper_cluster_id"`
	OnBrokerPerNode    bool   `form:"on_broker_per_node" json:"on_broker_per_node"`
	ClusterImage       string `form:"cluster_image" json:"cluster_image"`
	BrokerNum          int    `form:"broker_num,default=3" json:"broker_num"`
	StorageType        string `form:"storage_type" json:"storage_type"`
	StorageSize        int    `form:"storage_size" json:"storage_size"`
	Status             string `form:"status,default=ClusterNotRunning" json:"status"`
}

type KafkaTopic struct {
	Model
	Name              string `form:"name" json:"name" binding:"required"`
	Namespace         string `form:"namespace,default=kafka-operator" json:"namespace" binding:"required"`
	ClusterRefName    string `form:"cluster_ref_name" json:"cluster_ref_name" binding:"required"`
	TopicName         string `form:"topic_name" json:"topic_name" binding:"required"`
	Partitions        int    `form:"partitions,default=1" json:"partitions"`
	ReplicationFactor int    `form:"replication_factor,default=1" json:"replication_factor"`
}

func init() {
	var err error
	dbFile := setting.DBFile
	dbType := setting.DBType
	tablePrefix := setting.TablePrefix
	db, err = gorm.Open(dbType, dbFile)
	if err != nil {
		log.Fatalf("Failed to open the db file: %v", err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}
	db.AutoMigrate(&MysqlCluster{})
	db.AutoMigrate(&ZookeeperCluster{})
	db.AutoMigrate(&KafkaCluster{})
	db.AutoMigrate(&KafkaTopic{})
}

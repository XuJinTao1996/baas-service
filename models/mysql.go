package models

import (
	"encoding/base64"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func (mysqlCluster MysqlCluster) Response() MysqlCluster {
	mysqlCluster.Password = base64.StdEncoding.EncodeToString([]byte(mysqlCluster.Password))
	return mysqlCluster
}

// 通过 ID 获取指定的 mysql 实例
func GetMysqlcluster(id int) (MysqlCluster, bool) {
	var mysqlCluster MysqlCluster
	db.First(&mysqlCluster, id).Assign("password", "******")
	if mysqlCluster.ID > 0 {
		return mysqlCluster.Response(), true
	}
	return mysqlCluster.Response(), false
}

// 通过 mysqlcluster name 获取指定的 mysql 实例
func GetMysqlclusterByName(name string) (MysqlCluster, bool) {
	var mysqlCluster MysqlCluster
	db.Where("cluster_name = ?", name).First(&mysqlCluster)
	if mysqlCluster.ID > 0 {
		return mysqlCluster, true
	}
	return mysqlCluster.Response(), false
}

// 获取所有 mysql 实例的函数
func GetAllMysqlClusters() ([]MysqlCluster, int) {
	var mysqlCluster []MysqlCluster
	var count int
	db.Find(&mysqlCluster)
	db.Model(&MysqlCluster{}).Count(&count)
	return mysqlCluster, count
}

// 删除指定的 mysql 实例
func DeleteMysqlCluster(id int) bool {
	db.Where("id = ?", id).Delete(&MysqlCluster{})
	return true
}

// 判断是否存在该 mysql 实例
func ExistMysqlCluster(clusterName string) bool {
	var mysql MysqlCluster
	db.Select("id").Where("cluster_name = ?", clusterName).First(&mysql)
	if mysql.ID > 0 {
		return true
	}
	return false
}

// 创建数据库实例
func AddMysqlCluster(mysqlCluster *MysqlCluster) (MysqlCluster, bool) {
	db.Create(&mysqlCluster)
	return mysqlCluster.Response(), true
}

// 更新数据库实例
func UpdateMysqlCluster(mysqlCluster MysqlCluster, status string) (MysqlCluster, bool) {
	db.Model(&mysqlCluster).Update("status", status)
	return mysqlCluster, true
}

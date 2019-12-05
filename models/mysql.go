package models

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// 通过 ID 获取指定的 mysql 实例
func GetMysqlcluster(id int) (MysqlCluster, bool) {
	var mysqlCluster MysqlCluster
	db.First(&mysqlCluster, id)
	if mysqlCluster.ID > 0 {
		return mysqlCluster, true
	}
	return mysqlCluster, false
}

// 通过 mysqlcluster name 获取指定的 mysql 实例
func GetMysqlclusterByName(name string) (MysqlCluster, bool) {
	var mysqlCluster MysqlCluster
	db.First(&mysqlCluster, name)
	if mysqlCluster.ID > 0 {
		return mysqlCluster, true
	}
	return mysqlCluster, false
}

// 获取所有 mysql 实例的函数
func GetAllMysqlClusters() ([]MysqlCluster, int) {
	var mysqls []MysqlCluster
	var count int
	db.Find(&mysqls)
	db.Model(&MysqlCluster{}).Count(&count)
	return mysqls, count
}

// 删除指定的 mysql 实例
func DeleteMysqlcluster(id int) bool {
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
func AddMysqlCluster(mysqlCluster *MysqlCluster) bool {
	db.Create(&mysqlCluster)
	return true
}

// 更新数据库实例
func UpdateMysqlCluster(mysqlCluster MysqlCluster, status string) (MysqlCluster, bool) {
	db.Model(&mysqlCluster).Update("status", status)
	return mysqlCluster, true
}

func CloseDB() {
	defer db.Close()
}

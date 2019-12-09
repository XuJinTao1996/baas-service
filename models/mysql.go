package models

import (
	"encoding/base64"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func (mysqlCluster *MysqlCluster) SetPassword() {
	mysqlCluster.Password = base64.StdEncoding.EncodeToString([]byte(mysqlCluster.Password))
}

func (mysqlCluster *MysqlCluster) SetHost() {
	mysqlCluster.Host = mysqlCluster.RouterDeploymentName()
}

// 通过 ID 获取指定的 mysql 实例
func GetMysqlcluster(id int) (*MysqlCluster, error) {
	var mysqlCluster MysqlCluster
	err := db.First(&mysqlCluster, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &mysqlCluster, err
}

// 通过 mysqlcluster name 获取指定的 mysql 实例
func GetMysqlclusterByName(name string) (*MysqlCluster, error) {
	var mysqlCluster MysqlCluster
	err := db.Where("cluster_name = ?", name).First(&mysqlCluster).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &mysqlCluster, nil
}

// 获取所有 mysql 实例的函数
func GetAllMysqlClusters() ([]*MysqlCluster, error) {
	var mysqlCluster []*MysqlCluster
	err := db.Find(&mysqlCluster).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return mysqlCluster, err
}

// 获取 mysql 实例的总数
func GetArticleTotal() (int, error) {
	var count int
	if err := db.Model(&MysqlCluster{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// 删除指定的 mysql 实例
func DeleteMysqlCluster(id int) error {
	if err := db.Where("id = ?", id).Delete(&MysqlCluster{}).Error; err != nil {
		return err
	}
	return nil
}

// 判断是否存在该 mysql 实例
func ExistMysqlClusterByName(clusterName string) (bool, error) {
	var mysql MysqlCluster
	err := db.Select("id").Where("cluster_name = ?", clusterName).First(&mysql).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, nil
	}
	if mysql.ID > 0 {
		return true, nil
	}
	return false, nil
}

func ExistMysqlClusterByID(id int) (bool, error) {
	var mysql MysqlCluster
	err := db.Select("id").Where("id = ?", id).First(&mysql).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, nil
	}

	if mysql.ID > 0 {
		return true, nil
	}
	return false, nil
}

// 创建数据库实例
func AddMysqlCluster(mysqlCluster *MysqlCluster) error {
	if err := db.Create(&mysqlCluster).Error; err != nil {
		return err
	}
	return nil
}

// 更新数据库实例
func UpdateMysqlCluster(mysqlCluster *MysqlCluster, status string) (*MysqlCluster, bool) {
	db.Model(&mysqlCluster).Update("status", status)
	return mysqlCluster, true
}

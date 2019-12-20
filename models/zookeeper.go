package models

import "github.com/jinzhu/gorm"

func GetAllZookeeperCluster() ([]*ZookeeperCluster, error) {
	var zookeeperCluster []*ZookeeperCluster
	err := db.Find(&zookeeperCluster).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return zookeeperCluster, nil
}

func GetZookeeperClusterByID(id int) (*ZookeeperCluster, error) {
	var zookeeperCluster ZookeeperCluster
	err := db.First(&zookeeperCluster, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &zookeeperCluster, err
}

func GetZookeeperClusterByName(name string) (*ZookeeperCluster, error) {
	var zookeeperCluster ZookeeperCluster
	err := db.Where("cluster_name = ?", name).First(&zookeeperCluster).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &zookeeperCluster, err
}

func ExistZookeeperClusterByID(id int) (bool, error) {
	var zookeeperCluster ZookeeperCluster
	err := db.Select("id").Where("id = ?", id).First(&zookeeperCluster).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if zookeeperCluster.ID > 0 {
		return true, nil
	}

	return false, nil
}

func ExistZookeeperClusterByName(name string) (bool, error) {
	var zookeeperCluster ZookeeperCluster
	err := db.Select("id").Where("cluster_name = ?", name).First(&zookeeperCluster).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if zookeeperCluster.ID > 0 {
		return true, nil
	}

	return false, nil
}

func GetZookeeperClusterTotal() (int, error) {
	var count int
	if err := db.Model(ZookeeperCluster{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func AddZookeeperCluster(zc *ZookeeperCluster) error {
	if err := db.Create(zc).Error; err != nil {
		return err
	}
	return nil
}

func CheckZookeeperClusterUsageStatus(zid int) (bool, error) {
	var kafkaCluster KafkaCluster
	err := db.Select("id").Where("zookeeper_cluster_id = ?", zid).First(&kafkaCluster).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if kafkaCluster.ID > 0 {
		return true, nil
	}
	return false, nil
}

func DeleteZookeeperCluster(id int) error {
	if err := db.Where("id = ?", id).Delete(&ZookeeperCluster{}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateZookeeperCluster(zc *ZookeeperCluster, status string) (*ZookeeperCluster, bool) {
	db.Model(&zc).Update("status", status)
	return zc, true
}

package models

import "github.com/jinzhu/gorm"

func GetAllKafkaCluster() ([]*KafkaCluster, error) {
	var kafkaClusters []*KafkaCluster
	err := db.Find(&kafkaClusters).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return kafkaClusters, nil
}

func GetKafkaClusterTotal() (int, error) {
	var count int
	if err := db.Model(KafkaCluster{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func AddKafkaCluster(kc *KafkaCluster) error {
	if err := db.Create(kc).Error; err != nil {
		return err
	}
	return nil
}

func GetKafkaClusterByID(id int) (*KafkaCluster, error) {
	var kafkaCluster KafkaCluster
	err := db.First(&kafkaCluster, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &kafkaCluster, err
}

func ExistKafkaClusterByID(id int) (bool, error) {
	var kafkaCluster KafkaCluster
	err := db.Select("id").Where("id = ?", id).First(&kafkaCluster).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, nil
	}

	if kafkaCluster.ID > 0 {
		return true, nil
	}

	return false, nil
}

func ExistKafkaClusterByName(clusterName string) (bool, error) {
	var kafkaCluster KafkaCluster
	err := db.Select("id").Where("cluster_name = ?", clusterName).First(&kafkaCluster).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, nil
	}

	if kafkaCluster.ID > 0 {
		return true, nil
	}

	return false, nil
}

func DeleteKafkaCluster(id int) error {
	if err := db.Where("id = ?", id).Delete(&KafkaCluster{}).Error; err != nil {
		return err
	}
	return nil
}

func GetKafkaClusterByName(name string) (*KafkaCluster, error) {
	var kafkaCluster KafkaCluster
	err := db.Where("cluster_name = ?", name).First(&kafkaCluster).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &kafkaCluster, err
}

func UpdateKafkaCluster(kc *KafkaCluster, status string) (*KafkaCluster, bool) {
	db.Model(&kc).Update("status", status)
	return kc, true
}

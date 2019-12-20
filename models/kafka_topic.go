package models

import "github.com/jinzhu/gorm"

func GetAllKafkaTopic() ([]*KafkaTopic, error) {
	var kafkaTopic []*KafkaTopic
	err := db.Find(&kafkaTopic).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return kafkaTopic, nil
}

func GetKafkaTopicTotal() (int, error) {
	var count int
	if err := db.Model(KafkaTopic{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func AddKafkaTopic(kt *KafkaTopic) error {
	if err := db.Create(kt).Error; err != nil {
		return err
	}
	return nil
}

func GetKafkaTopicByID(id int) (*KafkaTopic, error) {
	var kafkaTopic KafkaTopic
	err := db.First(&kafkaTopic, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &kafkaTopic, err
}

func ExistKafkaTopicByID(id int) (bool, error) {
	var kafkaTopic KafkaTopic
	err := db.Select("id").Where("id = ?", id).First(&kafkaTopic).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, nil
	}

	if kafkaTopic.ID > 0 {
		return true, nil
	}

	return false, nil
}

func DeleteKafkaTopic(id int) error {
	if err := db.Where("id = ?", id).Delete(&KafkaTopic{}).Error; err != nil {
		return err
	}
	return nil
}

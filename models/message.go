package models

import "github.com/sjlleo/traceSysBackend/database"

type Messages struct {
	ID     uint   `gorm:"primarykey; comment: '送信 ID'"`
	Type   uint   `gorm:"type:int; comment: '送信类型'"`
	Method uint   `gorm:"method:int; comment: '送信方式'"`
	Model  string `gorm:"model:string; comment: '送信模板'"`
}

func (m *Messages) TableName() string {
	return "messages"
}

func CreateMessages(m *Messages) error {
	db := database.GetDB()
	if err := db.Create(m).Error; err != nil {
		return err
	}
	return nil
}

func UpdateMessages(m *Messages) error {
	db := database.GetDB()
	if err := db.Model(m).Where("id =?", m.ID).Updates(&m).Error; err != nil {
		return err
	}
	return nil
}

func DeleteMessages(m *Messages) error {
    db := database.GetDB()
	if err := db.Where("id =?", m.ID).Delete(m).Error; err != nil {
		return err
	}
	return nil
}

func GetMessages(m *Messages) error {
	db := database.GetDB()
    if err := db.Model(m).Where("id =?", m.ID).First(&m).Error; err!= nil {
        return err
    }
	return nil
}
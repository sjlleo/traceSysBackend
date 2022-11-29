package models

import "github.com/sjlleo/traceSysBackend/database"

type Template struct {
	ID     uint   `gorm:"primarykey; comment: '送信 ID'"`
	Type   uint   `gorm:"type:int; comment: '送信类型'"`
	Method uint   `gorm:"method:int; comment: '送信方式'"`
	Model  string `gorm:"model:string; comment: '送信模板'"`
}

const (
	RTTExceed        = 1
	PacketLossExceed = 2
)

func (m *Template) TableName() string {
	return "Template"
}

func CreateTemplate(m *Template) error {
	db := database.GetDB()
	if err := db.Create(m).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTemplate(m *Template) error {
	db := database.GetDB()
	if err := db.Model(m).Where("id =?", m.ID).Updates(&m).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTemplate(m *Template) error {
	db := database.GetDB()
	if err := db.Where("id =?", m.ID).Delete(m).Error; err != nil {
		return err
	}
	return nil
}

func GetTemplate(exceedType int) (Template, error) {
	db := database.GetDB()
	m := Template{}
	if err := db.Model(m).Where("id =?", exceedType).First(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

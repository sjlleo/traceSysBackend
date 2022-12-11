package models

import (
	"errors"
	"log"

	"github.com/sjlleo/traceSysBackend/database"
)

type IPReviews struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	IP         string    `gorm:"type:varchar(65);" json:"ip"`
	Prefix     uint      `gorm:"type:int" json:"prefix"`
	ASN        uint      `gorm:"type:int" json:"asn"`
	Country    string    `json:"country"`
	Province   string    `json:"province"`
	City       string    `json:"city"`
	Domain     string    `json:"domain"`
	Authorized uint      `gorm:"type:int" json:"authStatus"`
	CreatedAt  LocalTime `json:"createdTime"`
}

func AddReview(r IPReviews) error {
	db := database.GetDB()
	log.Println(r)
	err := db.Model(&IPReviews{}).Create(&r).Error
	return err
}

func ModReview(r *IPReviews) error {
	db := database.GetDB()
	// 如果已经审核通过了，那么就不应该再允许被修改
	res := db.Model(&IPReviews{}).Where("id = ?", r.ID).Where("authorized = ?", 0).Updates(r)

	if res.RowsAffected == 0 {
		return errors.New("您提交的 IP 段已经审核通过，如需修改请重新提起一个新的 IP 修正工单")
	}

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func DeleteReview(id uint) error {
	db := database.GetDB()
	res := db.Where("id =?", id).Delete(&IPReviews{})
	return res.Error
}

func (a *Admin) PassReview(review_id uint) {
	db := database.GetDB()
	db.Model(&IPReviews{}).Where("id = ?", review_id).Update("authorized", 1)
}
func (n *Normal) PassReview(review_id uint) {

}

func (a *Admin) DeclineReview(review_id uint) {
	db := database.GetDB()
	db.Model(&IPReviews{}).Where("id = ?", review_id).Update("authorized", 2)
}

func (n *Normal) DeclineReview(review_id uint) {

}

func (a *Admin) SearchReview(p *PaginationQ) error {
	var ipReviews []IPReviews
	db := database.GetDB()
	tx := db.Model(&IPReviews{})
	// 查找条件
	if p.Parm != "" {
		tx = tx.Where("ip like ?", "%"+p.Parm+"%")
	}
	total, err := crudAll(p, tx, &ipReviews)
	if err != nil {
		return err
	} else {
		p.Total = uint(total)
		p.Data = ipReviews
		return nil
	}
}

func (n *Normal) SearchReview(p *PaginationQ) error {
	var ipReviews []IPReviews
	db := database.GetDB()
	tx := db.Model(&IPReviews{})
	// 查找条件
	if p.Parm != "" {
		tx = tx.Where("ip like ?", "%"+p.Parm+"%")
	}
	total, err := crudAll(p, tx, &ipReviews)
	if err != nil {
		return err
	} else {
		p.Total = uint(total)
		p.Data = ipReviews
		return nil
	}
}

package database

import (
	"GoProject/global"
	"gorm.io/gorm"
)

type ArticleTagItem struct {
	Label string `json:"label"`
	Total int    `json:"total"`
}

type MiniArticle struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
	Hits  int    `json:"hits"`
}

type Article struct {
	gorm.Model
	Title         string         `json:"title"`
	User          User           `gorm:"foreignKey:Author"`
	Author        int            `json:"author"`
	Content       string         `json:"content"`
	Private       bool           `json:"private"`
	Hits          int            `json:"hits"`
	Top           bool           `json:"top"`
	Recommend     bool           `json:"recommend"`
	ArticleTypes  []ArticleType  `json:"article_types" gorm:"many2many:article_types"`
	ArticleLabels []ArticleLabel `json:"article_labels" gorm:"many2many:article_labels"`
}

func (article *Article) AfterUpdate(tx *gorm.DB) error {
	var count int64
	temp := tx.Model(Article{}).Where("top=true").Order("updated_at desc").Offset(global.MaxTopArticle).Count(&count)
	if count > 0 {
		return temp.Update("top", false).Error
	}
	return nil
}

type ArticleType struct {
	gorm.Model
	Display string `json:"display"`
	Order   int    `json:"order"`
}

type ArticleLabel struct {
	gorm.Model
	Display string `json:"display"`
	Order   int    `json:"order"`
}

func PopularList() (data []MiniArticle, err error) {
	var result []MiniArticle
	err = DBM.Db.Model(&Article{}).Where("private=?", false).Order("hits").Limit(10).Find(&result).Error
	return
}

func ArchiveList() (data []ArticleTagItem, err error) {
	tempTable := DBM.Db.Model(&Article{}).Select("strftime('%Y年%m月', created_at) as label").Where("private=false")
	err = DBM.Db.Table("(?) as u", tempTable).Select("label, count(1) as total").Group("label").Find(&data).Error
	return
}

func TypeOverflow() (data []ArticleTagItem, err error) {
	return
}

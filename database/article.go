package database

import (
	"GoProject/global"
	"GoProject/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sync"
)

var localLock sync.Mutex

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
	Title     string `json:"title"`
	UserID    int    `json:"author"`
	User      User   `json:"user" gorm:"foreignKey:UserID"`
	Content   string `json:"content"`
	Private   bool   `json:"private"`
	Hits      int    `json:"hits"`
	Top       bool   `json:"top"`
	Recommend bool   `json:"recommend"`
}

func (article *Article) Hit(req *gin.Context, user *User) error {
	localLock.Lock()
	defer localLock.Unlock()

	key := fmt.Sprintf("%s_%v_article_lock", req.ClientIP(), article.ID)
	cacheKey := util.GetMd5(key)
	_, exists := util.GetCache(cacheKey)
	if exists {
		return nil
	}
	DBM.Db.Create(&UserHistory{User: *user, Article: *article})
	DBM.Db.Model(&article).Update("hits", article.Hits+1)
	util.SetCacheWithDefault(cacheKey, key)
	return nil
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

func (at *ArticleType) BeforeCreate(_ *gorm.DB) error {
	if at.Order == 0 {
		at.Order = int(at.ID * 10)
	}
	return nil
}

type ArticleLabel struct {
	gorm.Model
	Display string `json:"display"`
	Order   int    `json:"order"`
	Color   string `json:"color"`
}

func (al *ArticleLabel) BeforeCreate(_ *gorm.DB) error {
	if al.Order == 0 {
		al.Order = int(al.ID * 10)
	}
	if al.Color == "" {
		al.Color = util.RandColor()
	}
	return nil
}

type ArticleTypeRelation struct {
	gorm.Model
	ArticleID uint        `json:"article_id"`
	Article   Article     `json:"article"`
	TypeID    uint        `json:"type_id"`
	Type      ArticleType `json:"type"`
}

type ArticleLabelRelation struct {
	gorm.Model
	ArticleID uint         `json:"article_id"`
	Article   Article      `json:"article"`
	LabelID   uint         `json:"label_id"`
	Label     ArticleLabel `json:"label"`
}

type Comment struct {
	gorm.Model
	//Target Comment `json:"target"`
	UserID uint `json:"user_id"`
	User   User `json:"user"`
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
	tempTable := DBM.Db.Model(&ArticleTypeRelation{}).Select("type_id, count(1) as total").Group("type_id")
	err = DBM.Db.Table("(?) as u", tempTable).Joins("left join article_types at on u.type_id = at.id").Select("at.display as label, u.total").Find(&data).Error
	return
}

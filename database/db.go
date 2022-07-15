package database

import (
	"GoProject/global"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strconv"
)

var (
	DBM *DbManager
)

type DbManager struct {
	Db                 *gorm.DB
	defaultSelectCount int
	maxSelectCount     int
}

func (dm *DbManager) Migrate() {
	if err := dm.Db.AutoMigrate(
		&User{},
		&Article{},
		&ArticleType{},
		&ArticleLabel{},
	); err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("migrate database panic")
	}
	global.Logger.Info("migrate database success")
}

func (dm *DbManager) Transaction(fn func(tx *gorm.DB) error) {
	tx := dm.Db.Begin()
	if err := fn(tx); err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func (dm *DbManager) First(target interface{}, query interface{}, args ...interface{}) error {
	return dm.Db.Where(query, args).Order("updated_at desc").First(target).Error
}

func (dm *DbManager) List(target interface{}) error {
	return dm.Db.Order("updated_at desc").Find(target).Error
}

func (dm *DbManager) Create(target interface{}) error {
	return dm.Db.Create(target).Error
}

func (dm *DbManager) Paginate(target interface{}, req *gin.Context) error {
	size, _ := strconv.Atoi(req.DefaultQuery("page_size", strconv.Itoa(dm.defaultSelectCount)))
	page, _ := strconv.Atoi(req.DefaultQuery("page_number", "1"))

	return dm.PaginateList(target, size, page)
}

func (dm *DbManager) PaginateList(target interface{}, size int, page int) error {
	return dm.Db.Scopes(PaginateScope(size, page)).Find(target).Error
}

func InitDb() {
	db, err := gorm.Open(sqlite.Open(global.DbPath))
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("init Db panic")
	}

	DBM = &DbManager{db, 10, 100}
	DBM.Migrate()
}

func PaginateScope(pageSize int, pageNumber int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNumber < 0 || pageSize < 0 {
			db.AddError(errors.New("参数错误"))
			return db
		}
		if pageSize > DBM.maxSelectCount {
			pageSize = DBM.maxSelectCount
		}
		offset := (pageNumber - 1) * pageSize
		return db.Order("updated_at desc").Offset(offset).Limit(pageSize)
	}
}

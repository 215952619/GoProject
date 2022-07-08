package database

import (
	"GoProject/global"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	DBM *DbManager
)

type DbManager struct {
	db *gorm.DB
}

func (dm *DbManager) Migrate() {
	if err := dm.db.AutoMigrate(
		&User{},
		&Article{},
	); err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("migrate database panic")
	}
	global.Logger.Info("migrate database success")
}

func (dm *DbManager) First(target interface{}, query interface{}, args ...interface{}) error {
	return dm.db.Where(query, args).First(target).Error
}

func (dm *DbManager) List(target interface{}) error {
	return dm.db.Find(target).Error
}

func (dm *DbManager) Create(target interface{}) error {
	return dm.db.Create(target).Error
}

func InitDb() {
	db, err := gorm.Open(sqlite.Open(global.DbPath))
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("init db panic")
	}

	DBM = &DbManager{db}
	DBM.Migrate()
	DB = db
}

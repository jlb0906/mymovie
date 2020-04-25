// file: common/aria2.go

package db

import (
	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jlb0906/mymovie/common"
	"github.com/jlb0906/mymovie/model"
	log "github.com/kataras/golog"
)

var db *gorm.DB

type Config struct {
	Mysql struct{ Url string }
}

func init() {
	var conf Config
	var err error
	if _, err = toml.DecodeFile(common.ConfigFile, &conf); err != nil {
		log.Fatal(err)
	}
	db, err = gorm.Open("mysql", conf.Mysql.Url)
	if err != nil {
		log.Fatal(err)
	}
	db.SingularTable(true)
	AutoMigrate(db)
}

func Get() *gorm.DB {
	return db
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.Movie{})
}

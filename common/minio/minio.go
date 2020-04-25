// file: common/minio.go

package minio

import (
	"github.com/BurntSushi/toml"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jlb0906/mymovie/common"
	log "github.com/kataras/golog"
	"github.com/minio/minio-go/v6"
)

var (
	cli *minio.Client
	c   *MinioConf
)

type Config struct {
	Minio MinioConf
}

type MinioConf struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	Location        string
	WorkerCount     int
}

func init() {
	var conf Config
	var err error
	if _, err = toml.DecodeFile(common.ConfigFile, &conf); err != nil {
		panic(err)
	}
	c = &conf.Minio
	cli, err = minio.New(c.Endpoint, c.AccessKeyID, c.SecretAccessKey, c.UseSSL)
	if err != nil {
		log.Fatal(err)
	}
	err = cli.MakeBucket(c.BucketName, c.Location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := cli.BucketExists(c.BucketName)
		if errBucketExists == nil && exists {
			log.Infof("We already own %s", c.BucketName)
		} else {
			log.Error(err)
		}
	} else {
		log.Infof("Successfully created %s", c.BucketName)
	}
}

func Get() *minio.Client {
	return cli
}

func GetConf() *MinioConf {
	return c
}

// file: common/aria2.go

package aria2

import (
	"context"
	"github.com/BurntSushi/toml"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jlb0906/mymovie/common"
	log "github.com/kataras/golog"
	"github.com/zyxar/argo/rpc"
	"sync"
	"time"
)

var (
	cli    rpc.Client
	n      rpc.Notifier
	c      *Aria2
	mu     sync.RWMutex
	inited bool
)

type Config struct {
	Aria2 Aria2
}

type Aria2 struct {
	Uri     string
	Token   string
	Timeout int
	Prefix  string
}

func Get() rpc.Client {
	mu.Lock()
	defer mu.Unlock()
	if inited {
		return cli
	}
	if n == nil {
		log.Warn("please set notifier firstly")
	}
	var conf Config
	var err error
	if _, err = toml.DecodeFile(common.ConfigFile, &conf); err != nil {
		log.Fatal(err)
	}
	c = &conf.Aria2
	cli, err = rpc.New(context.TODO(), c.Uri, c.Token, time.Duration(c.Timeout)*time.Second, n)
	if err != nil {
		log.Fatal(err)
	}
	inited = true
	return cli
}

func Set(notifier rpc.Notifier) {
	n = notifier
}

func GetConf() *Aria2 {
	return c
}

package service

import (
	"github.com/jlb0906/mymovie/common/aria2"
	mconf "github.com/jlb0906/mymovie/common/minio"
	"github.com/jlb0906/mymovie/model"
	"github.com/jlb0906/mymovie/service/upload"
	log "github.com/kataras/golog"
	"github.com/minio/minio-go/v6"
	"github.com/zyxar/argo/rpc"
	"path/filepath"
	"strings"
	"time"
)

var Aria2Service = NewAria2Service()

type aria2Service struct {
	e upload.Engine
}

func NewAria2Service() *aria2Service {
	a := new(aria2Service)
	a.e = new(upload.AsyncEngine)
	a.e.Run()
	aria2.Set(a)
	return a
}

func (s *aria2Service) OnDownloadStart(events []rpc.Event) {
	log.Infof("%s started.", events)
	for _, e := range events {
		stat, err := aria2.Get().TellStatus(e.Gid)
		if err != nil {
			log.Error(err)
			continue
		}
		arr := strings.Split(stat.Files[0].URIs[0].URI, "/")
		MovieService.UpdateByGid(model.Movie{
			Gid:    e.Gid,
			Title:  arr[len(arr)-1],
			Uri:    stat.Files[0].URIs[0].URI,
			Status: stat.Status,
		})
	}
}

func (aria2Service) OnDownloadPause(events []rpc.Event) { log.Infof("%s paused.", events) }
func (aria2Service) OnDownloadStop(events []rpc.Event)  { log.Infof("%s stopped.", events) }

func (s *aria2Service) OnDownloadComplete(events []rpc.Event) {
	log.Infof("%s completed.", events)
	r := &upload.WorkerReq{
		Events: events,
		Action: func(i interface{}) {
			Action(i.([]rpc.Event))
		},
	}
	s.e.Submit(r)
}

func (aria2Service) OnDownloadError(events []rpc.Event) { log.Infof("%s error.", events) }

func (aria2Service) OnBtDownloadComplete(events []rpc.Event) {
	log.Infof("bt %s completed.", events)
}

func Action(events []rpc.Event) {
	for _, e := range events {
		stat, err := aria2.Get().TellStatus(e.Gid)
		if err != nil {
			log.Error(err)
			continue
		}
		// Upload the file
		filePath := stat.Files[0].Path
		arr := strings.Split(filePath, "/")
		objectName := arr[len(arr)-1]
		filePath = filepath.Join(aria2.GetConf().Prefix, filePath)
		// Upload the file with FPutObject
		n, err := mconf.Get().FPutObject(mconf.GetConf().BucketName, objectName, filePath, minio.PutObjectOptions{})
		if err != nil {
			log.Error(err)
			return
		}
		// 生成7天访问地址
		u, err := mconf.Get().PresignedGetObject(mconf.GetConf().BucketName, objectName,
			604800*time.Second, map[string][]string{})
		if err != nil {
			log.Error(err)
			return
		}
		MovieService.UpdateByGid(model.Movie{
			Gid:    e.Gid,
			Title:  objectName,
			Uri:    u.String(),
			Status: stat.Status,
		})
		log.Infof("Successfully uploaded %s of size %d", objectName, n)
	}
}

func (aria2Service) AddURI(uri string) (string, error) {
	gid, err := aria2.Get().AddURI(uri)
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Infof("添加了下载任务 uri: %v gid: %v", uri, gid)
	return gid, nil
}

func (aria2Service) Remove(gid string) error {
	gid, err := aria2.Get().Remove(gid)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("删除了下载任务：%v", gid)
	return nil
}

func (aria2Service) Pause(gid string) error {
	log.Info("Received handler.Pause request")

	gid, err := aria2.Get().Pause(gid)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("暂停了下载任务：%v", gid)

	return nil
}

func (aria2Service) TellStatus(gid string) (*rpc.StatusInfo, error) {
	log.Info("Received handler.TellStatus request")

	info, err := aria2.Get().TellStatus(gid)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Infof("下载任务的状态：%v", info)
	return &info, nil
}

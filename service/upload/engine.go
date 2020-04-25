package upload

import "github.com/jlb0906/mymovie/common/minio"

type Engine interface {
	Run()
	Submit(r *WorkerReq)
}

type AsyncEngine struct {
	s Scheduler
}

func (a *AsyncEngine) Submit(r *WorkerReq) {
	a.s.Submit(r)
}

func (a *AsyncEngine) Run() {
	a.s = new(AsyncScheduler)
	a.s.Run()
	for i := 0; i < minio.GetConf().WorkerCount; i++ {
		createWorker(a.s.WorkerChannel(), a.s)
	}
}

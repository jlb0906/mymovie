package upload

import (
	"github.com/jlb0906/mymovie/common"
	"github.com/zyxar/argo/rpc"
)

func createWorker(in chan *WorkerReq, s Scheduler) {
	go func() {
		for {
			s.WorkerReady(in)
			events := <-in
			worker(events)
		}
	}()
}

type WorkerReq struct {
	Events []rpc.Event
	Action common.WorkerAction
}

func worker(r *WorkerReq) {
	r.Action(r.Events)
}

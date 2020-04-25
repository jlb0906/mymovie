package upload

type Scheduler interface {
	WorkerReady(chan *WorkerReq)
	Submit(*WorkerReq)
	WorkerChannel() chan *WorkerReq
	Run()
}

// 异步调度器
type AsyncScheduler struct {
	eventChan  chan *WorkerReq
	workerChan chan chan *WorkerReq
}

func (s *AsyncScheduler) WorkerReady(c chan *WorkerReq) {
	s.workerChan <- c
}

func (s *AsyncScheduler) Submit(events *WorkerReq) {
	s.eventChan <- events
}

func (s *AsyncScheduler) WorkerChannel() chan *WorkerReq {
	return make(chan *WorkerReq)
}

func (s *AsyncScheduler) Run() {
	s.eventChan = make(chan *WorkerReq)
	s.workerChan = make(chan chan *WorkerReq)

	go func() {
		var eventQueue []*WorkerReq
		var workerQueue []chan *WorkerReq

		for {
			var activeEvent *WorkerReq
			var activeWorker chan *WorkerReq
			if len(eventQueue) > 0 && len(workerQueue) > 0 {
				activeEvent = eventQueue[0]
				activeWorker = workerQueue[0]
			}
			select {
			case events := <-s.eventChan:
				eventQueue = append(eventQueue, events)
			case worker := <-s.workerChan:
				workerQueue = append(workerQueue, worker)
			case activeWorker <- activeEvent:
				eventQueue = eventQueue[1:]
				workerQueue = workerQueue[1:]
			}
		}
	}()
}

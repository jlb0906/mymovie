package common

type WorkerAction func(interface{})

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var OK = Response{Code: 0, Msg: "ok"}

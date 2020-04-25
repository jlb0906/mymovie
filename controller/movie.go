// file: web/controller/movie.go

package controller

import (
	"github.com/jlb0906/mymovie/common"
	"github.com/jlb0906/mymovie/common/validate"
	"github.com/jlb0906/mymovie/model"
	"github.com/jlb0906/mymovie/service"
	"github.com/kataras/iris/v12"
)

// MovieController is our /movies controller.
type MovieController struct {
}

func (c *MovieController) GetBy(ctx iris.Context, id string) {
	if len(id) == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	movie, b := service.MovieService.GetByID(id)
	var rsp common.Response
	if b {
		rsp = common.Response{
			Code: 0,
			Msg:  "ok",
			Data: movie,
		}
	} else {
		rsp = common.Response{
			Code: -1,
			Msg:  "error",
			Data: "",
		}
	}
	ctx.JSON(rsp)
}

func (c *MovieController) GetListBy(ctx iris.Context, page int) {
	if page <= 0 {
		page = 1
	}
	movieList := service.MovieService.SelectMany(model.Movie{}, page)
	var rsp common.Response
	rsp = common.Response{
		Code: 0,
		Msg:  "ok",
		Data: movieList,
	}
	ctx.JSON(rsp)
}

func (c *MovieController) PostAdd(ctx iris.Context) {
	var m model.Movie
	ctx.ReadJSON(&m)
	err := validate.Get().Struct(m)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	service.MovieService.Insert(m)
	ctx.JSON(common.OK)
}

func (c *MovieController) PostDeleteBy(ctx iris.Context, id string) {
	if len(id) == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	service.MovieService.DeleteByID(id)
	ctx.JSON(common.OK)
}

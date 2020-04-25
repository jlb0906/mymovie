// file: service/movie.go

package service

import (
	"github.com/jlb0906/mymovie/common"
	"github.com/jlb0906/mymovie/dao"
	"github.com/jlb0906/mymovie/model"
)

var MovieService = NewMovieService()

// NewMovieService returns the default movie service.
func NewMovieService() *movieService {
	return new(movieService)
}

type movieService struct {
}

// GetByID returns a movie based on its id.
func (s *movieService) GetByID(id string) (model.Movie, bool) {
	return dao.MovieRepository.Select(model.Movie{Id: id})
}

// DeleteByID deletes a movie by its id.
//
// Returns true if deleted otherwise false.
func (s *movieService) DeleteByID(id string) {
	dao.MovieRepository.Delete(model.Movie{
		Id: id,
	})
}

func (s *movieService) UpdateById(movie model.Movie) {
	dao.MovieRepository.Update(movie)
}

func (s *movieService) UpdateByGid(movie model.Movie) {
	dao.MovieRepository.UpdateByGid(movie)
}

func (s *movieService) Insert(movie model.Movie) {
	gid, _ := Aria2Service.AddURI(movie.Uri)
	movie.Gid = gid
	movie.Status = common.StatusAdd
	dao.MovieRepository.Insert(movie)
}

func (s *movieService) SelectMany(movie model.Movie, page int) []model.Movie {
	return dao.MovieRepository.SelectMany(movie, (page-1)*common.PageSize, common.PageSize)
}

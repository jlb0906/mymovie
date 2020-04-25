// file: dao/movie.go

package dao

import (
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/jlb0906/mymovie/common/db"
	"github.com/jlb0906/mymovie/model"
)

var MovieRepository = NewMovieRepository()

// NewMovieRepository returns a new movie memory-based repository,
// the one and only repository type in our example.
func NewMovieRepository() *movieRepository {
	return new(movieRepository)
}

// movieRepository is a "MovieRepository"
// which manages the movies using the memory data source (map).
type movieRepository struct {
}

func (r *movieRepository) Select(m model.Movie) (movie model.Movie, found bool) {
	db.Get().Where(&m).First(&movie)
	empty := model.Movie{}
	if movie != empty {
		found = true
	}
	return
}

func (r *movieRepository) SelectMany(m model.Movie, offset, limit int) (results []model.Movie) {
	db.Get().Offset(offset).Limit(limit).Where(&m).Find(&results)
	return
}

func (r *movieRepository) Insert(m model.Movie) {
	id, _ := uuid.NewV1()
	m.Id = id.String()
	db.Get().Create(&m)
}

func (r *movieRepository) Update(m model.Movie) {
	db.Get().Save(&m)
}

func (r *movieRepository) UpdateByGid(m model.Movie) {
	db.Get().Model(&model.Movie{}).Where(&model.Movie{
		Gid: m.Gid,
	}).Updates(&m)
}

func (r *movieRepository) Delete(m model.Movie) {
	db.Get().Delete(&m)
}

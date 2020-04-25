// file: model/movie.go

package model

// Movie is our sample data structure.
// Keep note that the tags for public-use (for our web app)
// should be kept in other file like "web/viewmodels/movie.go"
// which could wrap by embedding the model.Movie or
// declare new fields instead butwe will use this datamodel
// as the only one Movie model in our application,
// for the shake of simplicty.

type Movie struct {
	Id     string `json:"id"`
	Title  string `json:"title" gorm:"size:65535" validate:"required"`
	Uri    string `json:"uri" gorm:"size:65535" validate:"required"`
	Status string `json:"status"`
	Gid    string
}

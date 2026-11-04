package movieticket

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrMovieNotFound      = errors.New("movie not found")
	ErrMovieAlreadyExists = errors.New("movie already exists")
)

type Movie struct {
	ID       string
	Name     string
	Categoty MovieCategory
}

type MovieRepo struct {
	movieRepository map[string]*Movie
	movieMu         sync.Mutex
}

func NewMovieRepo() *MovieRepo {
	return &MovieRepo{
		movieRepository: make(map[string]*Movie),
	}
}

func (c *MovieRepo) GetMovies() ([]*Movie, error) {
	c.movieMu.Lock()
	defer c.movieMu.Unlock()

	resp := make([]*Movie, 0, len(c.movieRepository))
	for _, movie := range c.movieRepository {
		resp = append(resp, movie)
	}
	return resp, nil
}

func (c *MovieRepo) CreateMovie(name string, category MovieCategory) (*Movie, error) {
	c.movieMu.Lock()
	defer c.movieMu.Unlock()

	for _, movie := range c.movieRepository {
		if movie.Name == name {
			return nil, errors.New("Movie already exists")
		}
	}

	movie := &Movie{
		ID:       uuid.NewString(),
		Name:     name,
		Categoty: category,
	}
	c.movieRepository[movie.ID] = movie
	return movie, nil
}

func (m *MovieRepo) GetMovieByID(id string) (*Movie, error) {
	m.movieMu.Lock()
	defer m.movieMu.Unlock()

	for _, movie := range m.movieRepository {
		if movie.ID == id {
			return movie, nil
		}
	}

	return nil, ErrMovieNotFound
}

// type MovieRepository interface {
// 	GetMovie(id string) (*Movie, error)
// 	GetMovies() ([]*Movie, error)
// 	CreateMovie(movie *Movie) error
// 	// UpdateMovie(movie *Movie) error
// 	// DeleteMovie(id string) error
// }

// type InMemoryMovieRepository struct {
// 	movies map[string]*Movie // map to allow mutation hold reference ot Movie
// 	mu     sync.Mutex
// }

// var _ MovieRepository = (*InMemoryMovieRepository)(nil)

// func NewInMemoryMovieRepository() MovieRepository {
// 	return &InMemoryMovieRepository{
// 		movies: make(map[string]*Movie),
// 	}
// }

// func (r *InMemoryMovieRepository) GetMovie(id string) (*Movie, error) {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	movie, ok := r.movies[id]
// 	if !ok {
// 		return nil, ErrMovieNotFound
// 	}
// 	return movie, nil
// }

// func (r *InMemoryMovieRepository) GetMovies() ([]*Movie, error) {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()

// 	movies := make([]*Movie, 0, len(r.movies))
// 	for _, movie := range r.movies {
// 		movies = append(movies, movie)
// 	}
// 	return movies, nil
// }

// func (r *InMemoryMovieRepository) CreateMovie(movie *Movie) error {

// }

package movieticket

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrShowAlreadyExists = errors.New("show already exists")
	ErrShowNotFound      = errors.New("show not found")
)

type Show struct {
	ID, ScreenID, MovieID string
	StartTime, EndTime    time.Time
}

func NewShow(screenID, movieID string, start, end time.Time) *Show {
	return &Show{
		ID:        uuid.NewString(),
		ScreenID:  screenID,
		MovieID:   movieID,
		StartTime: start,
		EndTime:   end,
	}
}

type ShowRepository interface {
	GetShows() ([]*Show, error)
	CreateShow(screenID, movieID string, start, end time.Time) (*Show, error)
	GetShowByID(id string) (*Show, error)
}

type InMemoryShowRepository struct {
	showRepository map[string]*Show
	movie          *MovieRepo
	screenRepo     *ScreenRepo
	showMu         sync.Mutex
}

var _ ShowRepository = (*InMemoryShowRepository)(nil)

func NewInMemoryShowRepository(movie *MovieRepo, screenRepo *ScreenRepo) ShowRepository {
	return &InMemoryShowRepository{
		showRepository: make(map[string]*Show),
		movie:          movie,
		screenRepo:     screenRepo,
	}
}

func (c *InMemoryShowRepository) GetShows() ([]*Show, error) {
	c.showMu.Lock()
	defer c.showMu.Unlock()

	resp := make([]*Show, 0, len(c.showRepository))
	for _, show := range c.showRepository {
		resp = append(resp, show)
	}

	return resp, nil
}

func (c *InMemoryShowRepository) GetShowByID(id string) (*Show, error) {
	c.showMu.Lock()
	defer c.showMu.Unlock()

	show, ok := c.showRepository[id]
	if !ok {
		return nil, ErrShowNotFound
	}
	return show, nil
}

func (c *InMemoryShowRepository) CreateShow(screenID, movieID string, start, end time.Time) (*Show, error) {
	c.showMu.Lock()
	defer c.showMu.Unlock()

	// for _, show := range c.showRepository {
	// 	fmt.Println(show)
	// 	if show.MovieID == movieID && show.StartTime == start {
	// 		return nil, ErrShowAlreadyExists
	// 	}
	// }
	movie, err := c.movie.GetMovieByID(movieID)
	if err != nil {
		return nil, err
	}

	screen, err := c.screenRepo.GetScreen(screenID)
	if err != nil {
		return nil, ErrScreenUnavailable
	}

	newShow := NewShow(screen.ID, movie.ID, start, end)
	c.showRepository[newShow.ID] = newShow
	return newShow, nil
}

// func (r *InMemoryShowRepository) GetShowByID(id string) (*Show, error) {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	show, ok := r.Shows[id]
// 	if !ok {
// 		return nil, ErrShowNotFound
// 	}
// 	return show, nil
// }

// func (r *InMemoryShowRepository) GetShows() ([]*Show, error) {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()

// 	shows := make([]*Show, 0, len(r.Shows))
// 	for _, show := range r.Shows {
// 		shows = append(shows, show)
// 	}
// 	return shows, nil
// }

// func (r *InMemoryShowRepository) CreateShow(show *Show) error {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	if _, ok := r.Shows[show.ID]; ok {
// 		return ErrShowAlreadyExists
// 	}
// 	r.Shows[show.ID] = show
// 	return nil
// }

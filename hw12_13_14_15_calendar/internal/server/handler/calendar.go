package handler

import (
	"fmt"
	"net/http"

	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
)

type Calendar struct {
	repo repository.EventsRepo
}

func NewCalendarHandler(repo repository.EventsRepo) *Calendar {
	c := new(Calendar)
	c.repo = repo

	return c
}

func (c *Calendar) HelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}

	fmt.Fprintf(w, "Hello World1")
}

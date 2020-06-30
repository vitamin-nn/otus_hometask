package calendar

import (
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
)

type Calendar struct {
	er repository.EventsRepo
}

func NewCalendar(eventRepo repository.EventsRepo) *Calendar {
	return &Calendar{
		er: eventRepo,
	}
}

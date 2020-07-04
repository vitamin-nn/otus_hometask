package handler

import (
	"fmt"
	"net/http"

	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/usecase"
)

type EventHandler struct {
	eUseCase *usecase.EventUseCase
}

func NewEventHandler(eUseCase *usecase.EventUseCase) *EventHandler {
	c := new(EventHandler)
	c.eUseCase = eUseCase

	return c
}

func (c *EventHandler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}

	fmt.Fprintf(w, "Hello World1")
}

package apihttp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
)

func (a *apiSuite) iShouldReceiveEventWithTitle(title string) error {
	time.Sleep(10 * time.Second)

	a.rabbit.messagesMutex.RLock()
	defer a.rabbit.messagesMutex.RUnlock()

	for _, msg := range a.rabbit.messages {
		n := new(repository.Notification)
		err := json.Unmarshal(msg, n)
		if err != nil {
			return fmt.Errorf("unmarshal error: %v", err)
			continue
		}

		if n.EventTitle == title {
			return nil
		}
	}

	return fmt.Errorf("event with text '%s' was not found in %s (%s)", title, a.rabbit.messages, time.Now().String())
}

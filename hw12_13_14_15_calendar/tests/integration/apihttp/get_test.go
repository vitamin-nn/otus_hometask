package apihttp

import (
	"fmt"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/tests/httpclient/client/calendar_service"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/tests/httpclient/models"
)

func (a *apiSuite) iSendGetEventsDayRequest(beginAt string) error {
	tBegin, err := strfmt.ParseDateTime(beginAt)
	if err != nil {
		return err
	}

	body := &models.GetEventsRequest{
		BeginAt: tBegin,
	}
	params := calendar_service.NewCalendarServiceGetEventsDayParams().WithBody(body)

	apiKeyHeaderAuth := httptransport.APIKeyAuth("Grpc-Metadata-User_id", "header", "1")
	clientHTTP := a.getClient()
	resp, err := clientHTTP.CalendarServiceGetEventsDay(params, apiKeyHeaderAuth)
	if err != nil {
		return err
	}
	a.getResp = resp.GetPayload()

	return nil
}

func (a *apiSuite) iSendGetEventsWeekRequest(beginAt string) error {
	tBegin, err := strfmt.ParseDateTime(beginAt)
	if err != nil {
		return err
	}

	body := &models.GetEventsRequest{
		BeginAt: tBegin,
	}
	params := calendar_service.NewCalendarServiceGetEventsWeekParams().WithBody(body)

	apiKeyHeaderAuth := httptransport.APIKeyAuth("Grpc-Metadata-User_id", "header", "1")
	clientHTTP := a.getClient()
	resp, err := clientHTTP.CalendarServiceGetEventsWeek(params, apiKeyHeaderAuth)
	if err != nil {
		return err
	}
	a.getResp = resp.GetPayload()

	return nil
}

func (a *apiSuite) iSendGetEventsMonthRequest(beginAt string) error {
	tBegin, err := strfmt.ParseDateTime(beginAt)
	if err != nil {
		return err
	}

	body := &models.GetEventsRequest{
		BeginAt: tBegin,
	}
	params := calendar_service.NewCalendarServiceGetEventsMonthParams().WithBody(body)

	apiKeyHeaderAuth := httptransport.APIKeyAuth("Grpc-Metadata-User_id", "header", "1")
	clientHTTP := a.getClient()
	resp, err := clientHTTP.CalendarServiceGetEventsMonth(params, apiKeyHeaderAuth)
	if err != nil {
		return err
	}
	a.getResp = resp.GetPayload()

	return nil
}

func (a *apiSuite) eventCountShouldBe(expectedCount int) error {
	if len(a.getResp.Events.Events) != expectedCount {
		return fmt.Errorf("Unexpected events count: %d, but expected: %d", len(a.getResp.Events.Events), expectedCount)
	}

	return nil
}

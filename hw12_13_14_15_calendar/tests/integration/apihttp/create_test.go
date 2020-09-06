package apihttp

import (
	"fmt"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/tests/httpclient/client/calendar_service"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/tests/httpclient/models"
)

func (a *apiSuite) iSendCreateRequest(title, description, startAt, endAt, notifyAt string) error {
	tStart, err := strfmt.ParseDateTime(startAt)
	if err != nil {
		return err
	}
	tEnd, err := strfmt.ParseDateTime(endAt)
	if err != nil {
		return err
	}

	req := &models.ModifyEventRequest{
		Title:       title,
		Description: description,
		StartAt:     tStart,
		EndAt:       tEnd,
	}

	if notifyAt != "" {
		tNotify, err := strfmt.ParseDateTime(notifyAt)
		if err != nil {
			return err
		}
		req.NotifyAt = tNotify
	}

	body := &models.CreateEventRequest{
		Event: req,
	}

	params := calendar_service.NewCalendarServiceCreateEventParams()
	params.SetBody(body)

	apiKeyHeaderAuth := httptransport.APIKeyAuth("Grpc-Metadata-User_id", "header", "1")
	clientHTTP := a.getClient()
	resp, err := clientHTTP.CalendarServiceCreateEvent(params, apiKeyHeaderAuth)
	if err != nil {
		return err
	}
	a.modifyResp = resp.GetPayload()
	return nil
}

func (a *apiSuite) theRespShouldNotHasError() error {
	if a.modifyResp.Error != "" {
		return fmt.Errorf("Unexpected error %s", a.modifyResp.Error)
	}
	return nil
}

func (a *apiSuite) theRespShouldHasErrorText(text string) error {
	if a.modifyResp.Error != text {
		return fmt.Errorf("Incorrect error text: %s", a.modifyResp.Error)
	}
	return nil
}

func (a *apiSuite) theRespShouldHasCorrectEventID() error {
	if a.modifyResp.Event.ID == 0 {
		return fmt.Errorf("Incorrect event id")
	}
	return nil
}

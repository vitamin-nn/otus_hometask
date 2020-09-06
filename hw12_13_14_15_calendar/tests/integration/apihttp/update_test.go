package apihttp

import (
	"fmt"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/tests/httpclient/client/calendar_service"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/tests/httpclient/models"
)

func (a *apiSuite) iSendUpdateRequest(id int32, title, description, startAt, endAt, notifyAt string) error {
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

	body := &models.UpdateEventRequest{
		Event:   req,
		EventID: id,
	}

	params := calendar_service.NewCalendarServiceUpdateEventParams()
	params.SetBody(body)

	apiKeyHeaderAuth := httptransport.APIKeyAuth("Grpc-Metadata-User_id", "header", "1")
	clientHTTP := a.getClient()
	resp, err := clientHTTP.CalendarServiceUpdateEvent(params, apiKeyHeaderAuth)
	if err != nil {
		return err
	}
	a.modifyResp = resp.GetPayload()
	return nil
}

func (a *apiSuite) theRespShouldHasTitle(title string) error {
	if a.modifyResp.Event.Title != title {
		return fmt.Errorf("Incorrect title: expected '%s' but actual '%s'", title, a.modifyResp.Event.Title)
	}
	return nil
}

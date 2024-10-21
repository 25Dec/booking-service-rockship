package openapis

import (
	"booking-service/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type larkCalendarAPI struct {
	SecretKey  string
	ClientId   string
	Timezone   string
	CalendarID string
	appToken   model.LarkAppToken
}

func NewLarkCalendarAPI(secretKey string, clientId string, timezone string, calendarID string) LarkCarlendarAPI {
	return &larkCalendarAPI{
		SecretKey:  secretKey,
		ClientId:   clientId,
		Timezone:   timezone,
		CalendarID: calendarID,
		appToken:   model.LarkAppToken{},
	}
}

func (l *larkCalendarAPI) GetTimeZone() (string, error) {
	return l.Timezone, nil
}

func (l *larkCalendarAPI) ObtainAppToken() (string, error) {
	if l.appToken.Expire < time.Now().Unix() || l.appToken.AppAccessToken == "" {
		httpClient := &http.Client{}
		req, err := http.NewRequest("POST", "https://open.larksuite.com/open-apis/auth/v3/app_access_token/internal", nil)
		if err != nil {
			return "", fmt.Errorf("LarkCalendarWebAPI - ObtainAppToken - http.NewRequest: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		reqBody := bytes.NewBuffer([]byte(fmt.Sprintf(`{"app_id":"%s","app_secret":"%s"}`, l.ClientId, l.SecretKey)))

		req.Body = io.NopCloser(reqBody)
		resp, err := httpClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("LarkCalendarWebAPI - ObtainAppToken - httpClient.Do: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("http status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("LarkCalendarWebAPI - ObtainAppToken - ioutil.ReadAll: %w", err)
		}
		var token model.LarkAppToken
		err = json.Unmarshal(body, &token)
		if err != nil {
			return "", fmt.Errorf("LarkCalendarWebAPI - ObtainAppToken - json.Unmarshal: %w", err)
		}

		token.Expire = time.Now().Unix() + token.Expire - 300
		l.appToken = token
	}

	return l.appToken.AppAccessToken, nil
}

func (l *larkCalendarAPI) CreateEvent(event model.LarkEventRequest) (model.LarkEvent, error) {
	httpClient := &http.Client{}
	reqBody, err := http.NewRequest("POST", fmt.Sprintf("https://open.larksuite.com/open-apis/calendar/v4/calendars/%s/events", l.CalendarID), nil)
	if err != nil {
		return model.LarkEvent{}, fmt.Errorf("LarkCalendarWebAPI - CreateEvent - http.NewRequest: %w", err)
	}

	token, err := l.ObtainAppToken()
	if err != nil {
		return model.LarkEvent{}, fmt.Errorf("LarkCalendarWebAPI - CreateEvent - l.ObtainAppToken: %w", err)
	}

	reqBody.Header.Set("Content-Type", "application/json")
	reqBody.Header.Set("Authorization", "Bearer "+token)

	body, err := json.Marshal(event)
	if err != nil {
		return model.LarkEvent{}, fmt.Errorf("LarkCalendarWebAPI - CreateEvent - json.Marshal: %w", err)
	}

	reqBody.Body = io.NopCloser(bytes.NewBuffer(body))

	resp, err := httpClient.Do(reqBody)
	if err != nil {
		return model.LarkEvent{}, fmt.Errorf("LarkCalendarWebAPI - CreateEvent - httpClient.Do: %w", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return model.LarkEvent{}, fmt.Errorf("LarkCalendarWebAPI - CreateEvent - ioutil.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return model.LarkEvent{}, fmt.Errorf("LarkCalendarWebAPI - CreateEvent - Status Code: %d - Body: %v", resp.StatusCode, string(body))
	}

	var eventData model.LarkResponse[model.LarkEventData]
	err = json.Unmarshal(body, &eventData)
	if err != nil {
		return model.LarkEvent{}, fmt.Errorf("LarkCalendarWebAPI - CreateEvent - json.Unmarshal: %w", err)
	}

	return eventData.Data.Event, nil
}

func (l *larkCalendarAPI) DeleteEvent(eventId string) error {
	httpClient := &http.Client{}
	reqBody, err := http.NewRequest("DELETE", fmt.Sprintf("https://open.larksuite.com/open-apis/calendar/v4/calendars/%s/events/%s", l.CalendarID, eventId), nil)
	if err != nil {
		return fmt.Errorf("LarkCalendarWebAPI - DeleteEvent - http.NewRequest: %w", err)
	}

	token, err := l.ObtainAppToken()
	if err != nil {
		return fmt.Errorf("LarkCalendarWebAPI - DeleteEvent - l.ObtainAppToken: %w", err)
	}

	reqBody.Header.Set("Content-Type", "application/json")
	reqBody.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(reqBody)
	if err != nil {
		return fmt.Errorf("LarkCalendarWebAPI - DeleteEvent - httpClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LarkCalendarWebAPI - DeleteEvent - Status Code: %d", resp.StatusCode)
	}

	return nil
}

func (l *larkCalendarAPI) CreateScheduleParticipant(eventId string, attendees []model.LarkEventAttendee) error {
	httpClient := &http.Client{}
	urlAttendees := fmt.Sprintf("https://open.larksuite.com/open-apis/calendar/v4/calendars/%s/events/%s/attendees", l.CalendarID, eventId)

	reqBody, err := http.NewRequest("POST", urlAttendees, nil)
	if err != nil {
		return fmt.Errorf("LarkCalendarWebAPI - CreateScheduleParticipant - http.NewRequest: %w", err)
	}

	token, err := l.ObtainAppToken()
	if err != nil {
		return fmt.Errorf("LarkCalendarWebAPI - CreateScheduleParticipant - l.ObtainAppToken: %w", err)
	}

	reqBody.Header.Set("Content-Type", "application/json")
	reqBody.Header.Set("Authorization", "Bearer "+token)

	attendeesRequest := struct {
		Attendees []model.LarkEventAttendee `json:"attendees"`
	}{
		Attendees: attendees,
	}

	body, err := json.Marshal(attendeesRequest)
	if err != nil {
		return fmt.Errorf("LarkCalendarWebAPI - CreateScheduleParticipant - json.Marshal: %w", err)
	}

	reqBody.Body = io.NopCloser(bytes.NewBuffer(body))

	resp, err := httpClient.Do(reqBody)
	if err != nil {
		return fmt.Errorf("LarkCalendarWebAPI - CreateScheduleParticipant - httpClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LarkCalendarWebAPI - CreateScheduleParticipant - http.StatusOK: %d", resp.StatusCode)
	}

	return nil
}

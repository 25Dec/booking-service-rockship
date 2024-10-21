package openapis

import "booking-service/internal/model"

type (

	// LarkCarlendarAPI -.
	LarkCarlendarAPI interface {
		CreateEvent(event model.LarkEventRequest) (model.LarkEvent, error)
		CreateScheduleParticipant(eventId string, attendees []model.LarkEventAttendee) error
		DeleteEvent(eventId string) error
		GetTimeZone() (string, error)
	}

	// LarkBaseAPI -.
	LarkBaseAPI interface {
		AppendRecord(record map[string]interface{}, tableName string) (string, error)
	}

	// EdtronautUserAPI -.
	EdtronautAPI interface {
		GetUserByID(id string) (model.EdtronautUser, error)
		GetUserByToken(token string) (model.EdtronautUser, error)
		GetUserCourses(userToken string) ([]model.EdtronautCourse, error)
	}
)

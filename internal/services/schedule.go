package services

import (
	"booking-service/internal/model"
	"booking-service/internal/openapis"
	"booking-service/internal/repositories"

	"booking-service/pkg/utils/errs"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	dateTimeFormat = "2006-01-02T15:04:05Z07:00"
)

// ScheduleService -.
type scheduleServiceImpl struct {
	AppointmentService AppointmentService
	MentorService      MentorService
	ScheduleRepo       repositories.ScheduleRepository
	Edtronaurt         openapis.EdtronautAPI
	LarkCarlendar      openapis.LarkCarlendarAPI
	LarkBase           openapis.LarkBaseAPI
	mutex              sync.Mutex
	appointmentCounter map[string]int
}

func NewScheduleService(sr repositories.ScheduleRepository,
	as AppointmentService,
	ms MentorService,
	edtronaut openapis.EdtronautAPI,
	larkCarlendar openapis.LarkCarlendarAPI,
	larkBase openapis.LarkBaseAPI) ScheduleService {
	return &scheduleServiceImpl{
		AppointmentService: as,
		ScheduleRepo:       sr,
		Edtronaurt:         edtronaut,
		MentorService:      ms,
		LarkCarlendar:      larkCarlendar,
		LarkBase:           larkBase,
		appointmentCounter: make(map[string]int),
	}
}

func fillDateOfSchedule(schedule *model.Schedule) error {

	timestamp, err := time.Parse(dateTimeFormat, schedule.StartAt)

	if err != nil {
		return err
	}
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return fmt.Errorf("ScheduleServiceImpl - fillDateOfSchedule - time.LoadLocation: %w", err)
	}
	localizedTime := timestamp.In(loc)

	schedule.Wday = localizedTime.Weekday().String()

	schedule.StartAt = localizedTime.Format(dateTimeFormat)
	hours := localizedTime.Hour()
	minutes := localizedTime.Minute()
	schedule.Interval.From = fmt.Sprintf("%02d:%02d", hours, minutes)
	schedule.Interval.To = fmt.Sprintf("%02d:%02d", hours+1, minutes)
	return nil
}

func (s *scheduleServiceImpl) createLarkEvent(
	schedule model.Schedule,
	mentor model.Mentor,
	learner model.EdtronautUser) (model.LarkEvent, error) {
	appointmentTime, err := time.Parse(dateTimeFormat, schedule.StartAt)
	if err != nil {
		return model.LarkEvent{}, fmt.Errorf("ScheduleServiceImpl- createLarkEvent - time.Parse: %w", err)
	}
	larkTimeZone, err := s.LarkCarlendar.GetTimeZone()
	if err != nil {
		return model.LarkEvent{}, fmt.Errorf("ScheduleServiceImpl - createLarkEvent - s.LarkCarlendar.GetTimeZone: %w", err)
	}

	parameters := map[string]any{
		"summary":     fmt.Sprintf("[Edtronaut] Mentoring session w %s", learner.FirstName),
		"description": fmt.Sprintf("[Edtronaut] Mentoring session w %s", learner.FirstName),
		"start_time": model.LarkEventTime{
			Timestamp: strconv.Itoa(int(appointmentTime.Unix())),
			TimeZone:  larkTimeZone,
		},
		"end_time": model.LarkEventTime{
			Timestamp: strconv.Itoa(int(appointmentTime.Unix() + 3600)),
			TimeZone:  larkTimeZone,
		},
	}
	larkRequest := model.NewLarkEventRequest(parameters)
	larkEvent, err := s.LarkCarlendar.CreateEvent(larkRequest)
	if err != nil {
		return model.LarkEvent{}, fmt.Errorf("ScheduleServiceImpl- createLarkEvent - s.LarkCarlendar.CreateEvent: %w", err)
	}

	// add event attendees
	attendees := []model.LarkEventAttendee{
		{
			Type:            "third_party",
			ThirdPartyEmail: learner.Email,
		},
		{
			Type:            "third_party",
			ThirdPartyEmail: mentor.Email,
		},
		{
			Type:			"third_party",
			ThirdPartyEmail: "educators@edtronaut.ai",
		},
	}

	err = s.LarkCarlendar.CreateScheduleParticipant(larkEvent.EventID, attendees)

	if err != nil {
		_ = s.LarkCarlendar.DeleteEvent(larkEvent.EventID)
		return model.LarkEvent{}, fmt.Errorf("ScheduleServiceImpl - createLarkEvent - s.LarkCarlendar.CreateScheduleParticipant: %w", err)
	}

	return larkEvent, nil
}

func (s *scheduleServiceImpl) createLarkBaseData(schedule model.Schedule, mentor model.Mentor, learner model.EdtronautUser, content string) map[string]interface{} {
	timestamp, err := time.Parse(dateTimeFormat, schedule.StartAt)
	if err != nil {
		return map[string]interface{}{
			"ScheduledID": schedule.ID,
		}
	}

	return map[string]interface{}{
		"ScheduleID":           schedule.ID,
		"Mentor Email":         mentor.Email,
		"Learner Email":        learner.Email,
		"Date":                 timestamp.UnixMilli(),
		"Mentor Confirmation":  "",
		"Learner Confirmation": "",
		"Content":              content,
		"Status":               "Incoming",
	}
}

func (s *scheduleServiceImpl) GetSchedulesByLearnerID(ctx context.Context, learnerID string) ([]model.Schedule, error) {
	schedules, err := s.ScheduleRepo.GetManyByLearnerID(ctx, learnerID)
	if err != nil {
		return nil, fmt.Errorf("ScheduleServiceImpl - GetSchedulesByLearnerID - s.ScheduleRepo.GetManyByLearnerID: %w", err)
	}

	for i := range schedules {
		fillDateOfSchedule(&schedules[i])
	}

	return schedules, nil

}

func (s *scheduleServiceImpl) GetScheduledAppointments(ctx context.Context, from string, to string, limit int, offset int) ([]model.Appointment, int64, error) {
	appointments, count, err := s.AppointmentService.GetAppointments(ctx, from, to, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("ScheduleServiceImpl - GetScheduledAppointments - s.AppointmentService.GetAppointments: %w", err)
	}

	for i := range appointments {
		schedule, err := s.GetScheduleByID(ctx, appointments[i].ScheduleID)
		if err != nil {
			return nil, 0, fmt.Errorf("ScheduleServiceImpl - GetScheduledAppointments - s.GetScheduleByID: %w", err)
		}
		appointments[i].Schedule = schedule
	}

	return appointments, count, nil
}

func (s *scheduleServiceImpl) GetAppointmentOfSchedule(ctx context.Context, scheduleID string) (model.Appointment, error) {
	appointment, err := s.AppointmentService.GetAppointmentByScheduleID(ctx, scheduleID)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("ScheduleServiceImpl - GetAppointmentOfSchedule - s.AppointmentService.GetAppointmentByScheduleID: %w", err)
	}
	schedule, err := s.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("ScheduleServiceImpl - GetAppointmentOfSchedule - s.GetScheduleByID: %w", err)
	}
	appointment.Schedule = schedule

	return appointment, nil
}

func (s *scheduleServiceImpl) ScheduleAppointment(ctx context.Context, learnerId string, appointmentAt string, content string) (model.Appointment, error) {
	timeStamp, err := time.Parse(dateTimeFormat, appointmentAt)
	if err != nil {
		return model.Appointment{}, errs.BadRequestError{Message: "Invalid time format"}
	}
	//check if the learner has an appointment in the same week
	weekNumber, year := timeStamp.ISOWeek()
	appointments, err := s.AppointmentService.GetLearnerAppointmentsInWeek(ctx, year, weekNumber, learnerId)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("ScheduleServiceImpl - ScheduleAppointment - s.AppointmentService.GetLearnerAppointmentsInWeek: %w", err)
	}

	if len(appointments) > 0 {
		return model.Appointment{}, errs.BadRequestError{Message: "You already have an appointment in this week"}
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	errChan := make(chan error)
	learnerChan := make(chan model.EdtronautUser)
	var mentor model.Mentor
	var appropriateSchedule model.Schedule
	go func() {
		learner, err := s.Edtronaurt.GetUserByID(learnerId)
		if err != nil {
			errChan <- err
			return
		}
		learnerChan <- learner
	}()

	schedules, err := s.GetAvailableSchedulesAtTime(ctx, appointmentAt)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("ScheduleServiceImpl - ScheduleAppointment - s.GetAvailableSchedulesAtTime: %w", err)
	}

	if len(schedules) == 0 {
		return model.Appointment{}, errs.NotFoundError{Message: "No available schedule"}
	}
	appropriateSchedule = schedules[0]
	for _, schedule := range schedules {
		if s.appointmentCounter[schedule.MentorID] < s.appointmentCounter[appropriateSchedule.MentorID] {
			appropriateSchedule = schedule
		}
	}
	mentor, err = s.MentorService.GetMentorByID(ctx, appropriateSchedule.MentorID)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("ScheduleServiceImpl - ScheduleAppointment - s.MentorService.GetMentorByID: %w", err)
	}

	var learner model.EdtronautUser

	select {
	case err := <-errChan:
		return model.Appointment{}, err
	case learner = <-learnerChan:
	}

	larkEvent, err := s.createLarkEvent(appropriateSchedule, mentor, learner)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("ScheduleServiceImpl - ScheduleAppointment - s.createLarkEvent: %w", err)
	}

	parsedTime, err := time.Parse(dateTimeFormat, appropriateSchedule.StartAt)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("ScheduleServiceImpl - ScheduleAppointment - time.Parse: %w", err)
	}
	weekNumberm, year := parsedTime.ISOWeek()
	appointment := model.Appointment{
		LearnerID:  learnerId,
		ScheduleID: appropriateSchedule.ID,
		Schedule:   appropriateSchedule,
		WeekNumber: weekNumberm,
		Year:       year,
		Detail:     larkEvent,
		Content:    content,
	}

	//send data to lark base
	go func() {
		larkBaseData := s.createLarkBaseData(appropriateSchedule, mentor, learner, content)
		_, _ = s.LarkBase.AppendRecord(larkBaseData, "appointments")
	}()

	appointment, err = s.AppointmentService.CreateAppointment(ctx, appointment)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("ScheduleServiceImpl - ScheduleAppointment - s.AppointmentService.CreateAppointment: %w", err)
	}

	//increment the appointment counter of the mentor
	s.appointmentCounter[appropriateSchedule.MentorID]++

	return appointment, nil
}

func tryCasting(value interface{}, target interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("ScheduleServiceImpl - tryCasting - json.Marshal: %w", err)
	}

	dataAsStr := strings.ReplaceAll(string(data), "\"", "")
	decodedBytes, err := base64.StdEncoding.DecodeString(dataAsStr)
	if err != nil {
		return fmt.Errorf("ScheduleServiceImpl - tryCasting - base64.StdEncoding.DecodeString: %w", err)
	}

	err = json.Unmarshal([]byte(decodedBytes), target)
	if err != nil {
		return fmt.Errorf("ScheduleServiceImpl - tryCasting - json.Unmarshal: %w", err)
	}

	return nil
}

func (s *scheduleServiceImpl) UnscheduleAppointment(ctx context.Context, scheduleID string) error {

	appointment, err := s.AppointmentService.GetAppointmentByScheduleID(ctx, scheduleID)
	if err != nil {
		return fmt.Errorf("ScheduleServiceImpl - CancelAppointment - s.AppointmentService.GetAppointmentByScheduleID: %w", err)
	}

	var detail model.LarkEvent

	castingErr := tryCasting(appointment.Detail, &detail)

	err = s.AppointmentService.DeleteAppointmentByScheduleID(ctx, scheduleID)
	if err != nil {
		return fmt.Errorf("ScheduleServiceImpl - CancelAppointment - s.AppointmentService.DeleteAppointmentByScheduleID: %w", err)
	}

	if castingErr == nil {
		_ = s.LarkCarlendar.DeleteEvent(detail.EventID)
	}
	cancelLog := map[string]interface{}{
		"schedule_id": scheduleID,
		"cancel_at":   time.Now().UnixMilli(),
	}

	s.LarkBase.AppendRecord(cancelLog, "cancel_logs")

	//decrement the appointment counter of the mentor
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.appointmentCounter[appointment.Schedule.MentorID]--
	return nil
}

func (s *scheduleServiceImpl) GetAvailableSchedulesAtTime(ctx context.Context, appointmentAt string) ([]model.Schedule, error) {
	schedules, err := s.ScheduleRepo.GetAvailableSchedulesAtTime(ctx, appointmentAt)
	if err != nil {
		return nil, fmt.Errorf("ScheduleServiceImpl - GetAvailableSchedulesAtTime - s.toAvailebleSchedules: %w", err)
	}
	for i := range schedules {
		fillDateOfSchedule(&schedules[i])
	}

	return schedules, nil
}

func (s *scheduleServiceImpl) GetAvailableSchedules(ctx context.Context, from string, to string) ([]model.Schedule, error) {
	schedules, err := s.ScheduleRepo.GetAvailableSchedules(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("ScheduleServiceImpl - GetAvailableSchedules - s.ScheduleRepo.GetAvailableSchedules: %w", err)
	}

	for i := range schedules {
		fillDateOfSchedule(&schedules[i])
	}

	return schedules, nil
}

func (s *scheduleServiceImpl) GetScheduleByID(ctx context.Context, id string) (model.Schedule, error) {
	schedule, err := s.ScheduleRepo.GetByID(ctx, id)
	if err != nil {
		return model.Schedule{}, fmt.Errorf("ScheduleServiceImpl - GetScheduleByID - s.ScheduleRepo.GetByID: %w", err)
	}

	fillDateOfSchedule(&schedule)
	return schedule, nil
}

func (s *scheduleServiceImpl) GetSchedulesByTime(ctx context.Context, scheduleAt string) ([]model.Schedule, error) {
	schedules, err := s.ScheduleRepo.GetAllByTime(ctx, scheduleAt)
	if err != nil {
		return nil, fmt.Errorf("ScheduleServiceImpl - GetScheduleByTime - s.ScheduleRepo.GetByTime: %w", err)
	}

	for i := range schedules {
		fillDateOfSchedule(&schedules[i])
	}

	return schedules, nil
}

func (s *scheduleServiceImpl) GetSchedulesByMentorID(ctx context.Context, mentorID string) ([]model.Schedule, error) {
	schedules, err := s.ScheduleRepo.GetManyByMentorID(ctx, mentorID)
	if err != nil {
		return nil, fmt.Errorf("ScheduleServiceImpl - GetSchedulesByMentorID - s.ScheduleRepo.GetManyByMentorID: %w", err)
	}
	for i := range schedules {
		fillDateOfSchedule(&schedules[i])
	}

	return schedules, nil
}

func (s *scheduleServiceImpl) GetSchedules(ctx context.Context, from string, to string, limit int, offset int) ([]model.Schedule, int64, error) {
	schedules, count, err := s.ScheduleRepo.GetMany(ctx, from, to, limit, offset)
	for i := range schedules {
		fillDateOfSchedule(&schedules[i])
	}

	if err != nil {
		return nil, 0, fmt.Errorf("ScheduleServiceImpl - GetSchedules - s.ScheduleRepo.GetMany: %w", err)
	}

	return schedules, count, nil
}

func (s *scheduleServiceImpl) CreateSchedule(ctx context.Context, schedule model.Schedule) (model.Schedule, error) {
	// check mentor is valid
	_, err := s.MentorService.GetMentorByID(ctx, schedule.MentorID)
	if err != nil {
		return model.Schedule{}, err
	}

	schedule, err = s.ScheduleRepo.Create(ctx, schedule)
	if err != nil {
		return model.Schedule{}, fmt.Errorf("ScheduleServiceImpl - CreateSchedule - s.ScheduleRepo.Create: %w", err)
	}
	fillDateOfSchedule(&schedule)
	return schedule, nil
}

func (s *scheduleServiceImpl) UpdateSchedule(ctx context.Context, schedule model.Schedule) (model.Schedule, error) {
	schedule, err := s.ScheduleRepo.Update(ctx, schedule)
	if err != nil {
		return model.Schedule{}, fmt.Errorf("ScheduleServiceImpl - UpdateSchedule - s.ScheduleRepo.Update: %w", err)
	}
	fillDateOfSchedule(&schedule)

	return schedule, nil
}

func (s *scheduleServiceImpl) DeleteSchedule(ctx context.Context, id string) error {
	err := s.ScheduleRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("ScheduleServiceImpl - DeleteSchedule - s.ScheduleRepo.Delete: %w", err)
	}

	return nil
}

package controller

import (
	"booking-service/internal/middleware"
	"booking-service/internal/model"
	"booking-service/internal/services"
	"booking-service/pkg/logger"
	"booking-service/pkg/utils/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScheduleRoutes struct {
	scheduleService services.ScheduleService
	logger          logger.Interface
}

func NewScheduleRoutes(handler *gin.RouterGroup,
	l logger.Interface,
	scheduleService services.ScheduleService,
	learnerMiddleware *middleware.LearnerMiddleware) *ScheduleRoutes {
	r := &ScheduleRoutes{scheduleService: scheduleService, logger: l}
	h := handler.Group("/schedule")
	{
		h.GET("/:id", r.getScheduleByID)
		h.POST("", r.createSchedule)
		h.PUT("", r.updateSchedule)
		h.DELETE("/:id", r.deleteSchedule)
	}
	h = handler.Group("/schedules")
	{
		h.GET("", r.getSchedules)
		h.GET("/mentor/:mentor_id", r.getSchedulesByMentorID)
		h.GET("/available", r.getAvailableSchedules)
		h.GET("/appointments", r.getAppointments)
		h.GET("/learner/:learner_id", r.getSchedulesByLearnerID)
	}

	h = handler.Group("/schedule/appointment", learnerMiddleware.VerifyLearner())
	{
		h.GET("/:schedule_id", r.getAppointmentByScheduleID)
		h.POST("", r.scheduleAppointment)
		h.DELETE("/:schedule_id", r.unscheduleAppointment)

	}

	return r

}

// @Summary     Get schedule by ID
// @Description Get schedule by ID
// @ID          get-schedule-by-id
// @Tags  	    schedule
// @Accept      json
// @Produce     json
// @Param       id path string true "Schedule ID"
// @Success     200 {object} model.Response[model.Schedule]
// @Failure     500 {object} model.ErrorResponse
// @Router      /schedule/{id} [get]
func (r *ScheduleRoutes) getScheduleByID(c *gin.Context) {
	id := c.Param("id")
	schedule, err := r.scheduleService.GetScheduleByID(c.Request.Context(), id)
	if err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - getScheduleByID")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[model.Schedule]{Data: schedule, Message: "success"})

}

// all schedules paginated with limit defaullt is 10 and offset default is 0 but "from" and "to" is required miliseconds
// @Summary     Get all schedules
// @Description Get all schedules
// @ID          get-all-schedules
// @Tags  	    schedule
// @Accept      json
// @Produce     json
// @Param       from query string false "Start time" default(2019-01-01 00:00:00+07)
// @Param       to query string false "End time" default(2222-01-01 00:00:00+07)
// @Param       limit query int false "Limit" default(10)
// @Param       offset query int false "Offset" default(0)
// @Success     200 {object} model.PaginationResponse[model.Schedule]
// @Failure     500 {object} model.ErrorResponse
// @Router      /schedules [get]
func (r *ScheduleRoutes) getSchedules(c *gin.Context) {
	var paginationParams model.PaginationParams
	if err := c.Bind(&paginationParams); err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - getSchedules")
		model.NewErrorResponse(c, errs.BadRequestError{Message: "invalid pagination parameters"})
		return
	}

	var fromToParams model.FromToParams
	if err := c.Bind(&fromToParams); err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - getSchedules")
		model.NewErrorResponse(c, errs.BadRequestError{Message: "invalid from to parameters"})
		return
	}

	schedules, count, err := r.scheduleService.GetSchedules(c.Request.Context(), fromToParams.From, fromToParams.To, paginationParams.Limit, paginationParams.Offset)
	if err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - getSchedules")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.PaginationResponse[model.Schedule]{
		Data:    schedules,
		Message: "success",
		Paging: model.Paging{
			Total:   count,
			PerPage: paginationParams.Limit,
			Page:    paginationParams.Offset,
		},
	})

}

// getby many mentor id
// @Summary     Get all schedules by mentor ID
// @Description Get all schedules by mentor ID
// @ID          get-all-schedules-by-mentor-id
// @Tags  	    schedule
// @Accept      json
// @Produce     json
// @Param       mentor_id path string true "Mentor ID"
// @Success     200 {object} model.Response[[]model.Schedule]
// @Failure     500 {object} model.ErrorResponse
// @Router      /schedules/mentor/{mentor_id} [get]
func (r *ScheduleRoutes) getSchedulesByMentorID(c *gin.Context) {
	mentorID := c.Param("mentor_id")
	schedules, err := r.scheduleService.GetSchedulesByMentorID(c.Request.Context(), mentorID)
	if err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - getSchedulesByMentorID")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[[]model.Schedule]{Data: schedules, Message: "success"})
}

// @Summary     Create schedule
// @Description Create schedule
// @ID          create-schedule
// @Tags  	    schedule
// @Accept      json
// @Produce     json
// @Param       schedule body model.CreateScheduleRequest true "Schedule"
// @Success     200 {object} model.Response[model.Schedule]
// @Failure     500 {object} model.ErrorResponse
// @Router      /schedule [post]
func (r *ScheduleRoutes) createSchedule(c *gin.Context) {
	var schedule model.Schedule
	if err := c.BindJSON(&schedule); err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - createSchedule")
		model.NewErrorResponse(c, errs.BadRequestError{Message: "invalid schedule"})
		return
	}

	schedule, err := r.scheduleService.CreateSchedule(c.Request.Context(), schedule)
	if err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - createSchedule")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[model.Schedule]{Data: schedule, Message: "success"})

}

// @Summary     Update schedule
// @Description Update schedule
// @ID          update-schedule
// @Tags  	    schedule
// @Accept      json
// @Produce     json
// @Param       schedule body model.Schedule true "Schedule"
// @Success     200 {object} model.Response[model.Schedule]
// @Failure     500 {object} model.ErrorResponse
// @Router      /schedule [put]
func (r *ScheduleRoutes) updateSchedule(c *gin.Context) {
	var schedule model.Schedule
	if err := c.BindJSON(&schedule); err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - updateSchedule")
		model.NewErrorResponse(c, errs.BadRequestError{Message: "invalid schedule"})
		return
	}

	schedule, err := r.scheduleService.UpdateSchedule(c.Request.Context(), schedule)
	if err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - updateSchedule")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[model.Schedule]{Data: schedule, Message: "success"})

}

// @Summary     Delete schedule
// @Description Delete schedule
// @ID          delete-schedule
// @Tags  	    schedule
// @Accept      json
// @Produce     json
// @Param       id path string true "Schedule ID"
// @Success     200 {object} model.Response[string]
// @Failure     500 {object} model.ErrorResponse
// @Router      /schedule/{id} [delete]
func (r *ScheduleRoutes) deleteSchedule(c *gin.Context) {
	id := c.Param("id")
	err := r.scheduleService.DeleteSchedule(c.Request.Context(), id)
	if err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - deleteSchedule")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[string]{
		Message: "success",
		Data:    id,
	})

}

// get available schedule from time to time at second
// @Summary     Get available schedule
// @Description Get available schedule
// @ID          get-available-schedule
// @Tags  	    schedule
// @Accept      json
// @Produce     json
// @Param       from query string false "Start time" default(2019-01-01 00:00:00+07)
// @Param       to query string false "End time" default(2222-01-01 00:00:00+07)
// @Success     200 {object} model.Response[[]model.Schedule]
// @Failure     500 {object} model.ErrorResponse
// @Router      /schedules/available [get]
func (r *ScheduleRoutes) getAvailableSchedules(c *gin.Context) {
	var fromToParams model.FromToParams
	if err := c.Bind(&fromToParams); err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - getAvailableSchedules")
		model.NewErrorResponse(c, errs.BadRequestError{Message: "invalid from to parameters"})
		return
	}

	schedules, err := r.scheduleService.GetAvailableSchedules(c.Request.Context(), fromToParams.From, fromToParams.To)
	if err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - getAvailableSchedules")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[[]model.Schedule]{Data: schedules, Message: "success"})
}

// schedule a appointment by schedule_time and learner id, those date in post body
// @Summary     Schedule appointment
// @Description Schedule appointment
// @ID          schedule-appointment
// @Tags  	    Learner
// @Accept      json
// @Produce     json
// @Param       appointment body model.ScheduleAppointmentRequest true "Appointment"
// @Param      	Authorization header string  true "Bearer <token>"
// @Success     200 {object} model.Response[model.Appointment]
// @Failure     500 {object} model.ErrorResponse
// @Router      /schedule/appointment [post]
func (r *ScheduleRoutes) scheduleAppointment(c *gin.Context) {
	var appointmentRequest model.ScheduleAppointmentRequest
	if err := c.BindJSON(&appointmentRequest); err != nil {
		model.NewErrorResponse(c, errs.BadRequestError{Message: "invalid appointment request"})
		return
	}
	leaner := c.MustGet("learner").(model.EdtronautUser)

	appointment, err := r.scheduleService.ScheduleAppointment(c.Request.Context(), leaner.UserID, appointmentRequest.ScheduleAt, appointmentRequest.Content)
	if err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - scheduleAppointment")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[model.Appointment]{Data: appointment, Message: "success"})

}

// unschedule a appointment by schedule_id
// @Summary     Unschedule appointment
// @Description Unschedule appointment
// @ID          unschedule-appointment
// @Tags  	    Learner
// @Accept      json
// @Produce     json
// @Param       schedule_id path string true "Schedule ID"
// @Param      	Authorization header string  true "Bearer <token>"
// @Success     200 {object} model.Response[string]
// @Failure     500 {object} model.ErrorResponse
// @Router      /schedule/appointment/{schedule_id} [delete]
func (r *ScheduleRoutes) unscheduleAppointment(c *gin.Context) {
	id := c.Param("schedule_id")
	err := r.scheduleService.UnscheduleAppointment(c.Request.Context(), id)
	if err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - unscheduleAppointment")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[string]{
		Message: "success",
		Data:    id,
	})
}

// @Summary     Get appointment by Schedule ID
// @Description Get appointment by Schedule ID
// @ID          get-appointment-by-schedule-id
// @Tags  	    schedule
// @Accept      json
// @Produce     json
// @Param       schedule_id path string true "Schedule ID"
// @Success     200 {object} model.Response[model.Appointment]
// @Failure     500 {object} model.ErrorResponse
// @Router      /schedule/appointment/{schedule_id} [get]
func (r *ScheduleRoutes) getAppointmentByScheduleID(c *gin.Context) {
	id := c.Param("schedule_id")
	appointment, err := r.scheduleService.GetAppointmentOfSchedule(c.Request.Context(), id)
	if err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - getAppointmentByScheduleID")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[model.Appointment]{Data: appointment, Message: "success"})

}

// @Summary     Get many appointments
// @Description Get many appointments
// @ID          get-many-appointments
// @Tags  	    schedule
// @Accept      json
// @Produce     json
// @Param       from query string false "Start time" default(2019-01-01 00:00:00+07)
// @Param       to query string false "End time" default(2222-01-01 00:00:00+07)
// @Param       limit query int false "Limit" default(10)
// @Param       offset query int false "Offset" default(0)
// @Success     200 {object} model.PaginationResponse[model.Appointment]
// @Failure     500 {object} model.ErrorResponse
// @Router      /schedules/appointments [get]
func (r *ScheduleRoutes) getAppointments(c *gin.Context) {
	var paginationParams model.PaginationParams
	if err := c.Bind(&paginationParams); err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - getAppointments")
		model.NewErrorResponse(c, errs.BadRequestError{Message: "invalid pagination parameters"})
		return
	}

	var fromToParams model.FromToParams
	if err := c.Bind(&fromToParams); err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - getAppointments")
		model.NewErrorResponse(c, err)
		return
	}

	appointments, count, err := r.scheduleService.GetScheduledAppointments(c.Request.Context(), fromToParams.From, fromToParams.To, paginationParams.Limit, paginationParams.Offset)
	if err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - getAppointments")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.PaginationResponse[model.Appointment]{
		Data:    appointments,
		Message: "success",
		Paging: model.Paging{
			Total:   count,
			PerPage: paginationParams.Limit,
			Page:    paginationParams.Offset,
		},
	})

}

// @Summary     Get all schedules by learner ID
// @Description Get all schedules by learner ID
// @ID          get-all-schedules-by-learner-id
// @Tags  	    schedule
// @Accept      json
// @Produce     json
// @Param       learner_id path string true "Learner ID"
// @Success     200 {object} model.Response[[]model.Schedule]
// @Failure     500 {object} model.ErrorResponse
// @Router      /schedules/learner/{learner_id} [get]
func (r *ScheduleRoutes) getSchedulesByLearnerID(c *gin.Context) {
	learnerID := c.Param("learner_id")
	schedules, err := r.scheduleService.GetSchedulesByLearnerID(c.Request.Context(), learnerID)
	if err != nil {
		r.logger.Error(err, "http - ScheduleRoutes - getAppointmentsByLearnerID")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[[]model.Schedule]{Data: schedules, Message: "success"})
}

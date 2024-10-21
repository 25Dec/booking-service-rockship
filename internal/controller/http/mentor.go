package controller

import (
	"booking-service/internal/model"
	"booking-service/internal/services"
	"booking-service/pkg/logger"
	"booking-service/pkg/utils/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MentorRoutes struct {
	mentorService services.MentorService
	l             logger.Interface
}

func NewMentorRoutes(handler *gin.RouterGroup, l logger.Interface, ms services.MentorService) *MentorRoutes {
	r := &MentorRoutes{ms, l}
	h := handler.Group("/mentor")
	{
		h.GET("/:id", r.getMentorByID)
		h.POST("", r.createMentor)
		h.PUT("", r.updateMentor)
		h.DELETE("/:id", r.deleteMentor)
	}
	h = handler.Group("/mentors")
	{
		h.GET("", r.getMentors)
	}
	return r
}

// @Summary     Get mentor by ID
// @Description Get mentor by ID
// @ID          get-mentor-by-id
// @Tags  	    mentor
// @Accept      json
// @Produce     json
// @Param       id path string true "Mentor ID"
// @Success     200 {object} model.Response[model.Mentor]
// @Failure     500 {object} model.ErrorResponse
// @Router      /mentor/{id} [get]
func (r *MentorRoutes) getMentorByID(c *gin.Context) {
	id := c.Param("id")
	mentor, err := r.mentorService.GetMentorByID(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "http - MentorRoutes - getMentorByID")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[model.Mentor]{
		Message: "Mentor",
		Data:    mentor,
	})
}

// @Summary     Get all mentors
// @Description Get all mentors
// @ID          get-all-mentors
// @Tags  	    mentor
// @Accept      json
// @Produce     json
// @Success     200 {object} model.Response[model.Mentor]
// @Failure     500 {object} model.ErrorResponse
// @Router      /mentors [get]
func (r *MentorRoutes) getMentors(c *gin.Context) {
	mentors, err := r.mentorService.GetMentors(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - MentorRoun - getMentors")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[[]model.Mentor]{
		Message: "Mentors",
		Data:    mentors,
	})
}

// @Summary     Create mentor
// @Description Create mentor
// @ID          create-mentor
// @Tags  	    mentor
// @Accept      json
// @Produce     json
// @Param       mentor body model.Mentor true "Mentor"
// @Success     200 {object} model.Response[model.Mentor]
// @Failure     500 {object} model.ErrorResponse
// @Router      /mentor [post]
func (r *MentorRoutes) createMentor(c *gin.Context) {
	var mentor model.Mentor
	if err := c.ShouldBindJSON(&mentor); err != nil {
		r.l.Error(err, "http - MentorRoun - createMentor")
		model.NewErrorResponse(c, errs.BadRequestError{Message: "invalid request body"})
		return
	}

	mentor, err := r.mentorService.CreatedMentor(c.Request.Context(), mentor)
	if err != nil {
		r.l.Error(err, "http - MentorRoun - createMentor")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[model.Mentor]{
		Message: "Mentor created",
		Data:    mentor,
	})
}

// @Summary     Update mentor
// @Description Update mentor
// @ID          update-mentor
// @Tags  	    mentor
// @Accept      json
// @Produce     json
// @Param       mentor body model.Mentor true "Mentor"
// @Success     200 {object} model.Response[model.Mentor]
// @Failure     500 {object} model.ErrorResponse
// @Router      /mentor [put]
func (r *MentorRoutes) updateMentor(c *gin.Context) {
	var mentor model.Mentor
	if err := c.ShouldBindJSON(&mentor); err != nil {
		r.l.Error(err, "http - MentorRoun - updateMentor")
		model.NewErrorResponse(c, errs.BadRequestError{Message: "invalid request body"})
		return
	}

	mentor, err := r.mentorService.UpdateMentor(c.Request.Context(), mentor)
	if err != nil {
		r.l.Error(err, "http - MentorRoun - updateMentor")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[model.Mentor]{
		Message: "Mentor updated",
		Data:    mentor,
	})
}

// @Summary     Delete mentor
// @Description Delete mentor
// @ID          delete-mentor
// @Tags  	    mentor
// @Accept      json
// @Produce     json
// @Param       id path string true "Mentor ID"
// @Success     200 {object} model.Response[string]
// @Failure     500 {object} model.ErrorResponse
// @Router      /mentor/{id} [delete]
func (r *MentorRoutes) deleteMentor(c *gin.Context) {
	id := c.Param("id")
	err := r.mentorService.DeleteMentor(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "http - MentorRoutes - deleteMentor")
		model.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, model.Response[any]{
		Message: "Mentor deleted",
		Data:    id,
	})
}

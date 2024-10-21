// Package app configures and runs application.
package app

import (
	"booking-service/config"
	"io"

	controller "booking-service/internal/controller/http"
	"booking-service/internal/middleware"
	"booking-service/internal/openapis"
	"booking-service/internal/repositories"

	service "booking-service/internal/services"

	"booking-service/pkg/httpserver"
	"booking-service/pkg/logger"

	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	env := os.Getenv("ENV")
	if env == "PROD" {
		err := os.MkdirAll("./.log", os.ModePerm)
		if err != nil {
			panic(err)
		}

		file, err := os.OpenFile(
			"./.log/server.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		gin.DefaultWriter = io.MultiWriter(os.Stdout, file)

	}
	l := logger.New(cfg.Log.Level)

	//cors config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.Cors.AllowedOrigins
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}

	// set cors

	//connect postgres with gorm
	db, err := gorm.Open(postgres.Open(cfg.PG.URL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Repositories
	mentorRepo := repositories.NewMentorRepositoryImpl(db)
	appointmentRepo := repositories.NewAppointmentRepositoryImpl(db)
	scheduleRepo := repositories.NewScheduleRepoImpl(db)

	// OpenAPIs
	larkCalendarAPI := openapis.NewLarkCalendarAPI(cfg.LarkCalendar.AppSecret,
		cfg.LarkCalendar.AppID,
		cfg.LarkCalendar.Timezone,
		cfg.LarkCalendar.CalendarID)
	larBaseAPI := openapis.NewLarkBaseAPI(cfg.LarkBase.AppID,
		cfg.LarkBase.AppSecret,
		cfg.LarkBase.BaseToken,
		map[string]string{
			"appointments": cfg.LarkBase.AppointmentTable,
			"cancel_logs":  cfg.LarkBase.CancelTable,
		},
	)
	edtronautApi := openapis.NewEdtronautAPI(cfg.Edtronaut.Domain)

	// middleware
	learnerMiddleware := middleware.NewVerifyLearner(edtronautApi)

	// Services
	mentorService := service.NewMentorService(mentorRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo)
	scheduleService := service.NewScheduleService(scheduleRepo,
		appointmentService,
		mentorService,
		edtronautApi,
		larkCalendarAPI,
		larBaseAPI)

	// HTTP Server
	handler := gin.New()
	handler.Use(cors.New(corsConfig))
	controller.NewRouter(handler, l, scheduleService, mentorService, learnerMiddleware)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}

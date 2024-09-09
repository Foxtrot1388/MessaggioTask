package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Foxtrot1388/MessaggioTask/api"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type httpServer struct {
	log     *slog.Logger
	srv     *http.Server
	usecase usecases
	r       *gin.Engine
}

func New(log *slog.Logger, s usecases) *httpServer {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	result := &httpServer{log: log, srv: srv, usecase: s, r: router}
	router.POST("/message/create", result.create)
	router.GET("/message/statistics", result.statistics)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return result

}

func (s *httpServer) Listen() {

	s.log.Info("Start listen port", "addr", s.srv.Addr)

	// запустим http сервер
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("listen: %s\n", err)
			os.Exit(1)
		}
	}()

	// запустим горутину отправки в кафку
	ctxjobr, canceljobr := context.WithCancel(context.Background())
	defer canceljobr()
	go s.usecase.StartJobOutboxRead(ctxjobr)

	// запустим горутину чтения из кафки
	ctxjobw, canceljobw := context.WithCancel(context.Background())
	defer canceljobw()
	go s.usecase.StartJobOutboxWrite(ctxjobw)

	// graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	s.log.Info("shutdown", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Error("listen: %s\n", err)
		os.Exit(1)
	}

	select {
	case <-ctx.Done():
		s.log.Info("shutdown timeout of 5 seconds.")
	}
	s.log.Info("shutdown server exiting")

}

// @Summary Create
// @Tags message
// @Description Create a message
// @ID create
// @Accept json
// @Produce json
// @Param input body []string true "message fo create models"
// @Success 200 {object} []model.OutputMessage "if all message have been create"
// @Failure 500 {object} response
// @Failure 400 {object} response "validation error"
// @Router /message/create [POST]
func (con *httpServer) create(c *gin.Context) {

	var req []string
	if err := c.BindJSON(&req); err != nil {
		con.logResponseWithError(c, http.StatusBadRequest, notCreateMessage, err)
		return
	}

	for _, mes := range req {
		if mes == "" {
			con.logResponseWithError(c, http.StatusBadRequest, notCreateMessage, emptyMessage)
			return
		}
	}

	result, err := con.usecase.Create(context.Background(), req)
	if err != nil {
		con.logResponseWithError(c, http.StatusInternalServerError, notCreateMessage, err)
		return
	} else {
		con.logResponse(c, http.StatusOK, result)
	}

}

// @Summary statistics
// @Tags message
// @Description Get statistics a message
// @ID statistics
// @Produce json
// @Param dateAt query string false "Some date at (format 2006-01-02)"
// @Param dateTo query string false "Some date to (format 2006-01-02)"
// @Success 200 {object} []model.StatMessage
// @Failure 500 {object} response
// @Router /message/statistics [GET]
func (con *httpServer) statistics(c *gin.Context) {

	dateAt := getQueryParamDate(c, "dateAt", time.Time{})
	dateTo := getQueryParamDate(c, "dateTo", time.Unix(1<<63-1, 0))

	result, err := con.usecase.GetStatistic(context.Background(), dateAt, dateTo)
	if err != nil {
		con.logResponseWithError(c, http.StatusInternalServerError, notGetStatistic, err)
		return
	} else {
		con.logResponse(c, http.StatusOK, result)
	}

}

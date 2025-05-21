package rest

import "github.com/gin-gonic/gin"

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(
	userTaskHandler *UserTaskHandler,
	processStartHandler *ProcessStartHandler,
) *ServerHTTP {
	engine := gin.New()

	// Use logger from Gin
	engine.Use(gin.Logger())

	api := engine.Group("/api")

	api.POST("start/:id", processStartHandler.Start)

	api.GET("data", userTaskHandler.GetData)
	api.POST("complete", userTaskHandler.Complete)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Run() error {
	err := sh.engine.Run(":3000")
	if err != nil {
		return err
	}
	return nil
}

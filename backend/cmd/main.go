package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"metroid_bookmarks/misc"
	"metroid_bookmarks/pkg/handler"
	"metroid_bookmarks/pkg/repository/redis"
	"metroid_bookmarks/pkg/repository/sql"
	"metroid_bookmarks/pkg/service"
	"os"
	"os/signal"
	"syscall"
)

var logger = misc.GetLogger()

// @title METROID BOOKMARKS API
// @version 1.0
// @description API Server for metroid bookmarks
// @host localhost:3000
// @BasePath /api/v1
func main() {
	config := misc.GetConfig()
	dbPool, err := sql.NewDbPool(config.Db.Dsn)
	if err != nil {
		logger.Errorf("failed to create db dbPool: %s\n", err.Error())
		return
	}
	SQL := sql.NewSQL(dbPool)
	if err != nil {
		logger.Errorf("failed to initialize db: %s\n", err.Error())
		return
	}
	redisClient, err := redis.NewRedisPool(config.Redis.Dsn)
	if err != nil {
		logger.Errorf("failed to initialize redis: %s\n", err.Error())
		return
	}
	newRedis := redis.NewRedis(redisClient)
	newService := service.NewService(SQL, newRedis)

	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	srv := new(misc.Server)
	go func() {
		if err = srv.Run(handler.InitRoutes(newService, config)); err != nil {
			logger.Errorf("error occured while running http server: %s", err.Error())
		}
	}()

	logger.Info("METROID BOOKMARKS API started.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("METROID BOOKMARKS API shutting down.")

	if err = srv.Shutdown(context.Background()); err != nil {
		logger.Errorf("error occured on server shutting down: %s\n", err.Error())
	}
	dbPool.Close()
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("Recovered from panic: %s", r)
		}
	}()

	if err = redisClient.Close(); err != nil {
		logger.Errorf("error occured on redis connection close: %s\n", err.Error())
	}
}

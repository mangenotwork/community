package websvc

import (
	"community/internal/local/websvc/router"
	"community/pkg/conf"
	"community/pkg/logger"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Server() {

	conf.InitLocalConfig()

	router.Router = gin.Default()

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", conf.LocalConf.Port),
		Handler:        router.Routers(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.InfoF("http服务启动 0.0.0.0:%d", conf.LocalConf.Port)
	if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		logger.ErrorF("http服务出现异常:%s\n", err.Error())
	}

}

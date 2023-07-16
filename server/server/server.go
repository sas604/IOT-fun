package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewServer() *http.Server {
	router := gin.Default()
	router.GET("/api/switches", GetAllSwitches)
	router.POST("/api/switch", SetSwitch)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return srv
}

func KillServer(srv *http.Server) {
	log.Println("Shutdown Server ...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.

	log.Println("Server exiting")
}

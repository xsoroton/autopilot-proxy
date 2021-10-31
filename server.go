package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/xsoroton/autopilot-proxy/handlers"
)

func main() {
	h := handlers.NewHandlersFromEnv()
	r := gin.Default()
	r.GET("/v1/contact/*contact_id", h.GetContact)
	r.POST("/v1/contact", h.PostContact) // POST support idempotency

	err := r.Run(":8080")
	log.Fatal(err)
}

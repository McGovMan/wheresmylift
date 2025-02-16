package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RootGet					godoc
//
//	@Summary	Redirect to swagger docs
//	@Tags		Root
//	@Success	307
//	@Header		307	{string}	Location	"docs/index.html"
//	@Router		/ [get]
func (s *Server) RootGet(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "docs/index.html")
}

// V0HealthCheck			godoc
//
//	@Summary		Get health of API
//	@Description	If accessing this endpoint via Cloudflare it will only accessible using the BetterStack user-agent https://betterstack.com/docs/uptime/frequently-asked-questions/#what-user-agent-does-uptime-use
//	@Tags			V0
//	@Success		204
//	@Router			/v0/healthcheck [get]
func (s *Server) V0HealthCheckGet(c *gin.Context) {
	c.Status(204)
}

package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	h "github.com/mcgovman/wheresmylift/api/internal/helpers"
	"github.com/rs/zerolog/log"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.CustomRecovery(func(ctx *gin.Context, recovered interface{}) {
		log.Error().Any("error", recovered).Msg("recovery middleware")
		h.RespondWithError(ctx, errors.New("a server error was encountered"), http.StatusInternalServerError)
	}))

	r.Use(func(ctx *gin.Context) {
		start := time.Now()

		contextId, err := uuid.Parse(ctx.Request.Header.Get("context-id"))
		if err != nil || contextId == uuid.Nil {
			ctx.Writer.Header().Set("context-id", uuid.New().String())
		} else {
			ctx.Writer.Header().Set("context-id", contextId.String())
		}

		requestLog := log.
			With().Str("ip", ctx.ClientIP()).Logger().
			With().Str("method", ctx.Request.Method).Logger().
			With().Str("path", ctx.Request.URL.Path).Logger().
			With().Str("context_id", ctx.Writer.Header().Get("context-id")).Logger()

		ctx.Next()

		requestLog = requestLog.With().Int64("latency_ns", time.Since(start).Nanoseconds()).Logger().
			With().Int("status", ctx.Writer.Status()).Logger()
		requestLog.Info().Msg("request_info")

		ctx.Next()
	})

	return r
}

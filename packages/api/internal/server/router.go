package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	h "github.com/mcgovman/wheresmylift/packages/api/internal/helpers"
	"github.com/rs/zerolog/log"
)

func newOrExistingUUIDv7(uuidStr string) uuid.UUID {
	contextId, err := uuid.Parse(uuidStr)
	// Unlikely to error - seems like it only would when it runs out of unique uuids
	uuidV7, _ := uuid.NewV7()
	if err != nil || contextId.Version() != 7 {
		return uuidV7
	}

	sec, nsec := contextId.Time().UnixTime()
	timeSinceCreation := time.Since(time.Unix(sec, nsec).UTC())
	if timeSinceCreation > 5*time.Minute || timeSinceCreation < time.Since(time.Now()) {
		return uuidV7
	}

	return contextId
}

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.CustomRecovery(func(ctx *gin.Context, recovered interface{}) {
		log.Error().Any("error", recovered).Msg("recovery middleware")
		h.RespondWithError(ctx, errors.New("a server error was encountered"), http.StatusInternalServerError)
	}))

	r.Use(func(ctx *gin.Context) {
		start := time.Now()

		ctx.Writer.Header().Set(
			"context-id",
			newOrExistingUUIDv7(ctx.Request.Header.Get("context-id")).String(),
		)

		requestLog := log.
			With().Str("ip", ctx.ClientIP()).Logger().
			With().Str("method", ctx.Request.Method).Logger().
			With().Str("path", ctx.Request.URL.Path).Logger().
			With().Str("context_id", ctx.Writer.Header().Get("context-id")).Logger()

		ctx.Next()

		requestLog = requestLog.With().Int64("latency_ns", time.Since(start).Nanoseconds()).Logger().
			With().Int("status", ctx.Writer.Status()).Logger()
		requestLog.Info().Msg("request_info")
	})

	return r
}

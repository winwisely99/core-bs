package cmd

import (
	"context"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/logger"
)

func getLoggerFromContext(ctx context.Context) *logger.Logger {
	return ctx.Value("logger").(*logger.Logger)
}

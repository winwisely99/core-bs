package ctx

import (
	"context"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/logger"
)

func GetLogger(ctx context.Context) *logger.Logger {
	return ctx.Value("logger").(*logger.Logger)
}

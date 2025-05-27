package janitor

import (
	"context"
	"time"

	"github.com/inmore/gopaste/internal/storage"
	"go.uber.org/zap"
)

func Run(ctx context.Context, log *zap.Logger, st storage.Storage) {
	ticker := time.NewTicker((time.Second))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			n, err := st.DeleteExpired()
			if err != nil {
				log.Error("janitor error", zap.Error(err))
			} else if n > 0 {
				log.Debug("janitor removed", zap.Int("count", n))
			}
		case <-ctx.Done():
			return
		}
	}
}

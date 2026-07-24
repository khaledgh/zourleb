// Package worker holds scheduled background jobs.
package worker

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/zourleb/zourleb-api/internal/models"
)

// Runner owns the DB handle and dispatches due jobs on a ticker.
type Runner struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewRunner(db *gorm.DB, log *slog.Logger) *Runner {
	return &Runner{db: db, log: log}
}

// Start runs the scheduling loop until ctx is cancelled.
func (r *Runner) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	r.tick(ctx) // run once on boot
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.tick(ctx)
		}
	}
}

func (r *Runner) tick(ctx context.Context) {
	r.expireBoosts(ctx)
	r.cleanupOTP(ctx)
	r.reconcileDepartures(ctx)
}

// expireBoosts marks active boosts past their end date as expired.
func (r *Runner) expireBoosts(_ context.Context) {
	res := r.db.Model(&models.Boost{}).
		Where("status = ? AND ends_at < ?", models.BoostActive, time.Now()).
		Update("status", models.BoostExpired)
	if res.Error != nil {
		r.log.Error("expireBoosts failed", "err", res.Error)
		return
	}
	if res.RowsAffected > 0 {
		r.log.Info("boosts expired", "count", res.RowsAffected)
	}
}

// cleanupOTP prunes phone-verification audit rows older than 30 days.
func (r *Runner) cleanupOTP(_ context.Context) {
	cutoff := time.Now().AddDate(0, 0, -30)
	res := r.db.Where("created_at < ?", cutoff).Delete(&models.PhoneVerification{})
	if res.Error != nil {
		r.log.Error("cleanupOTP failed", "err", res.Error)
		return
	}
	if res.RowsAffected > 0 {
		r.log.Debug("otp audit pruned", "count", res.RowsAffected)
	}
}

// reconcileDepartures marks open departures full when seats run out, and closes
// departures whose start date has passed.
func (r *Runner) reconcileDepartures(_ context.Context) {
	r.db.Model(&models.TourDeparture{}).
		Where("status = ? AND seats_booked >= capacity", models.DepartureOpen).
		Update("status", models.DepartureFull)
}

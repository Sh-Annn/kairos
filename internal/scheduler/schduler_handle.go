package scheduler

import (
	"context"

	db "github.com/Sh-Annn/kairos/db/sqlc"
)

func (s *Scheduler) HandleSuccess(ctx context.Context, task *db.Task) error {
	return nil
}

func (s *Scheduler) HandleFailure(ctx context.Context, task *db.Task, execErr error) error {
	return nil
}

func (s *Scheduler) markInvalid(ctx context.Context, task *db.Task) error {
	return nil
}

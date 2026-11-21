package scheduler

import (
	"context"
	"database/sql"
	"fmt"

	"time"

	db "github.com/Sh-Annn/kairos/db/sqlc"
	"github.com/Sh-Annn/kairos/internal/domain"
)

func (s *Scheduler) HandleSuccess(ctx context.Context, task *db.Task) error {
	//	log.Println("HandleSuccess")
	// must change updated_at
	// schedult_type: once, interval, cron
	var updateTaskParams db.UpdateTaskParams
	updateTaskParams.ID = task.ID
	updateTaskParams.UpdatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	switch task.ScheduleType {
	case db.ScheduleTypeOnce:
		updateTaskParams.Status = db.NullTaskStatus{
			TaskStatus: "completed",
			Valid:      true,
		}
	case db.ScheduleTypeInterval:
		updateTaskParams.NextRunAt = sql.NullTime{
			Time:  time.Now().Add(time.Duration(task.IntervalSeconds.Int32) * time.Second),
			Valid: true,
		}
		// case db.ScheduleTypeCron:
	}
	_, err := s.queries.UpdateTask(ctx, updateTaskParams)
	if err != nil {
		return fmt.Errorf("ScheduleType: %s, Error: %w", task.ScheduleType, err)
	}
	return nil
}

func (s *Scheduler) HandleFailure(ctx context.Context, task *db.Task, execErr domain.HTTPError) error {
	//	log.Println("handleFailure")
	// change status or next_run_at
	// Error: URL, Method, StatusCode, Timeout
	var updateTaskParams db.UpdateTaskParams
	updateTaskParams.ID = task.ID
	updateTaskParams.UpdatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	updateTaskParams.NextRunAt = sql.NullTime{
		Time:  time.Now().Add(time.Duration(task.RetryBackoff) * time.Second),
		Valid: true,
	}

	switch execErr {
	case domain.URLError, domain.MethodError:
		updateTaskParams.Status = db.NullTaskStatus{
			TaskStatus: "paused",
			Valid:      true,
		}

		//TODO: http.StatusCode
	}

	_, err := s.queries.UpdateTask(ctx, updateTaskParams)
	if err != nil {
		return fmt.Errorf("HTTP Bad Request: %s, Error: %w", execErr, err)
	}

	return nil
}

func (s *Scheduler) markInvalid(ctx context.Context, task *db.Task) error {
	// deleted
	//	log.Println("markInvalid")
	_, err := s.queries.UpdateTask(ctx, db.UpdateTaskParams{
		ID: task.ID,
		Status: db.NullTaskStatus{
			TaskStatus: "paused",
			Valid:      true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("cannot change task's status be deleted: %w", err)
	}
	return nil
}

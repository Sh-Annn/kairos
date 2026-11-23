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
		updateTaskParams.RetryCount = sql.NullInt32{
			Int32: 0,
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

func (s *Scheduler) HandleFailure(ctx context.Context, task *db.Task, execErr error) error {
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
	case domain.URLError, domain.MethodError, domain.ConfigError:
		updateTaskParams.Status = db.NullTaskStatus{
			TaskStatus: "paused",
			Valid:      true,
		}

	case domain.NewRequestError, domain.HTTPClientError, domain.StatusCodeError,
		domain.BodyError:
		updateTaskParams.NextRunAt = sql.NullTime{
			Time:  time.Now().Add(time.Duration((task.RetryCount+1)*task.RetryBackoff) * time.Second),
			Valid: true,
		}
		updateTaskParams.RetryCount = sql.NullInt32{
			Int32: task.RetryCount + 1,
			Valid: true,
		}
	default:
	}

	updateTask, err := s.queries.UpdateTask(ctx, updateTaskParams)
	if err != nil {
		return fmt.Errorf("HTTP Bad Request: %s, Error: %w", execErr, err)
	}

	if updateTask.RetryCount > updateTask.MaxRetry {
		err = s.deleteTask(ctx, &updateTask)
		if err != nil {
			return fmt.Errorf("retry_count >= max_retry, but cannot update status")
		}
	}

	return nil
}

func (s *Scheduler) deleteTask(ctx context.Context, task *db.Task) error {
	_, err := s.queries.UpdateTask(ctx, db.UpdateTaskParams{
		ID: task.ID,
		Status: db.NullTaskStatus{
			TaskStatus: "deleted",
			Valid:      true,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

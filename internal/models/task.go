package models

import (
	"errors"
	"fmt"
	"time"
)

type Task struct {
	TskId uint      `json:"tsk_id,omitempty"`
	EmpId uint      `json:"emp_id,omitempty"`
	Open  int64 `json:"open"`
	Close int64 `json:"close"`
	Closed bool     `json:"closed,omitempty"`
	Meta  string    `json:"meta"`
}

type ListTskRequest struct {
	PageSize uint `json:"page_size"`
	PageNumber uint `json:"page_number"`
}

type TaskFactoryConfig struct {
	//Just basic constraints, can be added more
	MinTaskLifespan time.Duration
}

type TaskFactory struct {
	fc TaskFactoryConfig
}

type TaskLifespanToShort struct {
	MinTaskLifespan time.Duration
	ProvidedTaskLifespan time.Duration
}

func (t TaskLifespanToShort) Error() string {
	return fmt.Sprintf(
		"Provided task lifespan is too short, min lifespan: %s, provided lifespan: %s",
		t.MinTaskLifespan,
		t.ProvidedTaskLifespan,
	)
}

func (tf TaskFactory) NewTask(EmpId uint, Open int64, Close int64, Closed bool, Meta string) (Task, error) {

	if Open > Close {
		return Task{}, errors.New("task opening date starts after closing date")
	}

	dif := time.Unix(Close, 0).Sub(time.Unix(Open, 0))
	if dif < tf.fc.MinTaskLifespan {
		return Task{}, TaskLifespanToShort{
			MinTaskLifespan: tf.fc.MinTaskLifespan,
			ProvidedTaskLifespan: dif,
		}
	}
	return Task{
		EmpId:  EmpId,
		Open:   Open,
		Close:  Close,
		Closed: Closed,
		Meta:   Meta,
	}, nil
}

func NewTaskFactory(fc TaskFactoryConfig) TaskFactory{
	return TaskFactory{fc: fc}
}

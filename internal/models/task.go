package models

import (
	"errors"
	"fmt"
	"time"
)

type Task struct {
	TskId uint      `json:"tsk_id,omitempty"`
	Open  int64 `json:"open"`
	Close int64 `json:"close"`
	Closed bool     `json:"closed,omitempty"`
	Meta  string    `json:"meta"`
	EmpId uint      `json:"emp_id"`
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

type taskLifespanToShort struct {
	MinTaskLifespan time.Duration
	ProvidedTaskLifespan time.Duration
}

func (t taskLifespanToShort) Error() string {
	return fmt.Sprintf(
		"Provided task lifespan is too short, min lifespan: %s, provided lifespan: %s",
		t.MinTaskLifespan,
		t.ProvidedTaskLifespan,
	)
}

func (tf TaskFactory) validate(EmpId uint, Open int64, Close int64, Closed bool, Meta string) error {
	if Open > Close {
		return errors.New("task opening date starts after closing date")
	}

	dif := time.Unix(Close, 0).Sub(time.Unix(Open, 0))
	if dif < tf.fc.MinTaskLifespan {
		return taskLifespanToShort{
			MinTaskLifespan: tf.fc.MinTaskLifespan,
			ProvidedTaskLifespan: dif,
		}
	}
	return nil
}

func (tf TaskFactory) NewTask(EmpId uint, Open int64, Close int64, Closed bool, Meta string) (Task, error) {
	if err := tf.validate(EmpId, Open, Close, Closed, Meta); err != nil {
		return Task{}, err
	}

	return Task{
		EmpId:  EmpId,
		Open:   Open,
		Close:  Close,
		Closed: Closed,
		Meta:   Meta,
	}, nil
}

// Validate - validates struct, returns validated struct or error
// Can be added mechanism of hashing data before store or passing it to UC
// It differs from NewTask, so when you're creating new task you don't know its ID. But when you need to update it, or something else you got id
func (tf TaskFactory) Validate(tsk *Task) (*Task, error){

	if err := tf.validate(tsk.EmpId, tsk.Open, tsk.Close, tsk.Closed, tsk.Meta); err != nil {
		return nil, err
	}

	return tsk, nil
}

func NewTaskFactory(fc TaskFactoryConfig) TaskFactory{
	return TaskFactory{fc: fc}
}

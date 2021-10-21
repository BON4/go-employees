package utils

import (
	"bufio"
	"goland-hello/internal/models"
	"strconv"
)

const sep = ','
const end = '\n'

//TODO SHITS NEEDS GENERICS FOR CLEANER IMPLEMENTATION

// EmpToByte - writes employee struct to buffer in csv RFC 4180 format, returns number of written bytes and error
func EmpToByte(emp *models.Employee,  bw *bufio.Writer) (int, error) {
	n := 0
	if m, err := bw.WriteString(strconv.FormatUint(uint64(emp.EmpId), 10)); err != nil {
		return 0, err
	} else {
		n += m
	}
	if err := bw.WriteByte(sep); err != nil {
		return 0, err
	}
	n += 1

	if m, err := bw.WriteString(emp.Fname); err != nil {
		return 0, err
	} else {
		n += m
	}
	if err := bw.WriteByte(sep); err != nil {
		return 0, err
	}
	n += 1

	if m, err := bw.WriteString(emp.Lname); err != nil {
		return 0, err
	} else {
		n += m
	}
	if err := bw.WriteByte(sep); err != nil {
		return 0, err
	}
	n += 1

	if m, err := bw.WriteString(strconv.FormatFloat(emp.Sal, 'g', -1, 32)); err != nil {
		return 0, err
	} else {
		n += m
	}
	if err := bw.WriteByte(end); err != nil {
		return 0, err
	}
	n += 1

	return n, nil
}

// TaskToByte - writes task struct to buffer in csv RFC 4180 format, returns number of written bytes and error
func TaskToByte(tsk *models.Task, bw *bufio.Writer) (int, error) {
	n := 0
	if m, err := bw.WriteString(strconv.FormatUint(uint64(tsk.TskId), 10)); err != nil {
		return 0, err
	} else {
		n += m
	}
	if err := bw.WriteByte(sep); err != nil {
		return 0, err
	}
	n += 1

	if m, err := bw.WriteString(strconv.FormatInt(tsk.Open, 10)); err != nil {
		return 0, err
	} else {
		n += m
	}
	if err := bw.WriteByte(sep); err != nil {
		return 0, err
	}
	n += 1

	if m, err := bw.WriteString(strconv.FormatInt(tsk.Close, 10)); err != nil {
		return 0, err
	} else {
		n += m
	}
	if err := bw.WriteByte(sep); err != nil {
		return 0, err
	}
	n += 1

	if m, err := bw.WriteString(strconv.FormatBool(tsk.Closed)); err != nil {
		return 0, err
	} else {
		n += m
	}
	if err := bw.WriteByte(sep); err != nil {
		return 0, err
	}
	n += 1

	if m, err := bw.WriteString(`"`+tsk.Meta+`"`); err != nil {
		return 0, err
	} else {
		n += m
	}
	if err := bw.WriteByte(sep); err != nil {
		return 0, err
	}
	n += 1

	if m, err := bw.WriteString(strconv.FormatUint(uint64(tsk.EmpId), 10)); err != nil {
		return 0, err
	} else {
		n += m
	}
	if err := bw.WriteByte(end); err != nil {
		return 0, err
	}
	n += 1

	return n, nil
}

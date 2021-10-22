package models

import (
	"fmt"
)

type Employee struct {
	EmpId uint    `json:"emp_id,omitempty"`
	Fname string  `json:"fname,omitempty"`
	Lname string  `json:"lname,omitempty"`
	Sal   float64 `json:"sal,omitempty"`
}

type ListEmpRequest struct {
	PageSize uint `json:"page_size"`
	PageNumber uint `json:"page_number"`
}

type EmployeeFactoryConfig struct {
	//Just basic constraints, can be added more
	MinFirstNameLength int
	MinLastNameLength int
	MinSalary float64
}

type EmployeeFactory struct {
	fc EmployeeFactoryConfig
}

func NewEmployeeFactory(fc EmployeeFactoryConfig) EmployeeFactory{
	return EmployeeFactory{fc: fc}
}

type firstNameTooShortError struct {
	MinFirstNameLength int
	ProvidedFirstNameLength int
}

func (p firstNameTooShortError) Error() string {
	return fmt.Sprintf(
		"Provided first name is too short, min length: %d, provided password lentgh: %d",
		p.MinFirstNameLength,
		p.ProvidedFirstNameLength,
	)
}

type lastNameTooShortError struct {
	MinLastNameLength int
	ProvidedLastNameLength int
}

func (p lastNameTooShortError) Error() string {
	return fmt.Sprintf(
		"Provided last name is too short, min length: %d, provided password lentgh: %d",
		p.MinLastNameLength,
		p.ProvidedLastNameLength,
	)
}

type salaryOutOfBoundsError struct {
	MinSalary float64
	ProvidedSalary float64
}

func (p salaryOutOfBoundsError) Error() string {
	return fmt.Sprintf(
		"Salary out of bounds, min salary is: %d, provided salary is: %d",
		p.MinSalary,
		p.ProvidedSalary,
	)
}

func (efc EmployeeFactory) validate(fname string, lname string, sal float64) error {
	if len(fname) < efc.fc.MinFirstNameLength {
		return firstNameTooShortError{
			MinFirstNameLength:      efc.fc.MinFirstNameLength,
			ProvidedFirstNameLength: len(fname),
		}
	}

	if len(lname) < efc.fc.MinLastNameLength {
		return lastNameTooShortError{
			MinLastNameLength:      efc.fc.MinLastNameLength,
			ProvidedLastNameLength: len(lname),
		}
	}

	if sal < efc.fc.MinSalary {
		return salaryOutOfBoundsError{
			MinSalary:      efc.fc.MinSalary,
			ProvidedSalary: sal,
		}
	}
	return nil
}

func (efc EmployeeFactory) NewUser(fname string, lname string, sal float64) (Employee, error) {
	if err := efc.validate(fname, lname, sal); err != nil {
		return Employee{}, err
	}

	return Employee{
		Fname: fname,
		Lname: lname,
		Sal:   sal,
	}, nil
}

// Validate - validates struct, returns validated struct or error
// Can be added mechanism of hashing data before store or passing it to UC
// It differs from NewUser, so when you're creating new employee you don't know its ID. But when you need to update it, or something else you got id
func (efc EmployeeFactory) Validate(emp *Employee) (*Employee, error) {
	if err := efc.validate(emp.Fname, emp.Lname, emp.Sal); err != nil {
		return nil, err
	}
	return emp, nil
}
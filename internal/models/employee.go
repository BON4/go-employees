package models

import "fmt"

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

type FirstNameTooShortError struct {
	MinFirstNameLength int
	ProvidedFirstNameLength int
}

func (p FirstNameTooShortError) Error() string {
	return fmt.Sprintf(
		"Provided first name is too short, min length: %d, provided password lentgh: %d",
		p.MinFirstNameLength,
		p.ProvidedFirstNameLength,
	)
}

type LastNameTooShortError struct {
	MinLastNameLength int
	ProvidedLastNameLength int
}

func (p LastNameTooShortError) Error() string {
	return fmt.Sprintf(
		"Provided last name is too short, min length: %d, provided password lentgh: %d",
		p.MinLastNameLength,
		p.ProvidedLastNameLength,
	)
}

type SalaryOutOfBoundsError struct {
	MinSalary float64
	ProvidedSalary float64
}

func (p SalaryOutOfBoundsError) Error() string {
	return fmt.Sprintf(
		"Salary out of bounds, min salary is: %d, provided salary is: %d",
		p.MinSalary,
		p.ProvidedSalary,
	)
}

func (efc EmployeeFactory) NewUser(fname string, lname string, sal float64) (Employee, error) {
	if len(fname) < efc.fc.MinFirstNameLength {
		return Employee{}, FirstNameTooShortError{
			MinFirstNameLength:      efc.fc.MinFirstNameLength,
			ProvidedFirstNameLength: len(fname),
		}
	}

	if len(lname) < efc.fc.MinLastNameLength {
		return Employee{}, LastNameTooShortError{
			MinLastNameLength:      efc.fc.MinLastNameLength,
			ProvidedLastNameLength: len(lname),
		}
	}

	if sal < efc.fc.MinSalary {
		return Employee{}, SalaryOutOfBoundsError{
			MinSalary:      efc.fc.MinSalary,
			ProvidedSalary: sal,
		}
	}

	return Employee{
		Fname: fname,
		Lname: lname,
		Sal:   sal,
	}, nil
}
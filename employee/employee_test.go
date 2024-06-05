package employee

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmployeeStore_CreateEmployee(t *testing.T) {
	store := NewEmployeeStore()
	emp := Employee{
		ID:       1,
		Name:     "Shivangi",
		Position: "Software Engineer",
		Salary:   100000,
	}

	err := store.CreateEmployee(emp)
	assert.NoError(t, err, "Error creating employee")

	err = store.CreateEmployee(emp)
	assert.Error(t, err, "Expected error when creating duplicate employee")
}

func TestEmployeeStore_GetEmployeeByID(t *testing.T) {
	store := NewEmployeeStore()
	emp := Employee{
		ID:       1,
		Name:     "Shivangi",
		Position: "Software Engineer",
		Salary:   100000,
	}

	store.CreateEmployee(emp)

	
	expected := emp
	actual, ok := store.GetEmployeeByID(1)
	assert.True(t, ok, "Expected employee found")
	assert.Equal(t, expected, actual, "Expected employee")

	
	_, ok = store.GetEmployeeByID(2)
	assert.False(t, ok, "Expected no employee for ID 2")
}

func TestEmployeeStore_UpdateEmployee(t *testing.T) {
	store := NewEmployeeStore()
	emp := Employee{
		ID:       1,
		Name:     "Shivangi",
		Position: "Software Engineer",
		Salary:   100000,
	}

	store.CreateEmployee(emp)

	
	emp.Name = "Shivani"
	ok := store.UpdateEmployee(emp)
	assert.True(t, ok, "Expected to update existing employee")

	nonExistingEmp := Employee{
		ID:       2,
		Name:     "Non Existing",
		Position: "Tester",
		Salary:   80000,
	}
	ok = store.UpdateEmployee(nonExistingEmp)
	assert.False(t, ok, "Expected not to update non-existing employee")
}

func TestEmployeeStore_DeleteEmployee(t *testing.T) {
	store := NewEmployeeStore()
	emp := Employee{
		ID:       1,
		Name:     "Shivangi",
		Position: "Software Engineer",
		Salary:   100000,
	}

	store.CreateEmployee(emp)


	ok := store.DeleteEmployee(1)
	assert.True(t, ok, "Expected to delete existing employee")

	ok = store.DeleteEmployee(2)
	assert.False(t, ok, "Expected not to delete non-existing employee")
}

func TestEmployeeStore_ListEmployees(t *testing.T) {
	store := NewEmployeeStore()
	emp1 := Employee{
		ID:       1,
		Name:     "Shivangi",
		Position: "Software Engineer",
		Salary:   100000,
	}
	emp2 := Employee{
		ID:       2,
		Name:     "Rashmi",
		Position: "Data Analyst",
		Salary:   90000,
	}

	store.CreateEmployee(emp1)
	store.CreateEmployee(emp2)

	expected := []Employee{emp1, emp2}
	actual := store.ListEmployees()
	assert.ElementsMatch(t, expected, actual, "Expected list of employees")
}

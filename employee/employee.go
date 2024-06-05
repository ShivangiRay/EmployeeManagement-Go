package employee

import (
	"fmt"
	"sync" // Import the sync package for mutexes
)

// Employee represents an employee structure
type Employee struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}

// EmployeeStore manages the employees in memory
type EmployeeStore struct {
	employees map[int]Employee
	mu        sync.RWMutex // Mutex for read-write access to the employees map
}

// EmployeeStoreInterface defines methods for managing employees
type EmployeeStoreInterface interface {
	CreateEmployee(emp Employee) error
	GetEmployeeByID(id int) (Employee, bool)
	UpdateEmployee(emp Employee) bool
	DeleteEmployee(id int) bool
	ListEmployees() []Employee
	ListEmployeesWithPagination(offset, limit int) []Employee
}

// NewEmployeeStore creates a new instance of EmployeeStore
func NewEmployeeStore() *EmployeeStore {
	return &EmployeeStore{
		employees: make(map[int]Employee),
	}
}

// CreateEmployee adds a new employee to the in-memory store
func (store *EmployeeStore) CreateEmployee(emp Employee) error {
	store.mu.Lock()         // Lock for write access
	defer store.mu.Unlock() // Unlock when done

	if _, ok := store.employees[emp.ID]; ok {
		return fmt.Errorf("employee with ID %d already exists", emp.ID)
	}
	store.employees[emp.ID] = emp
	return nil
}

// GetEmployeeByID retrieves an employee from the in-memory store by their ID
func (store *EmployeeStore) GetEmployeeByID(id int) (Employee, bool) {
	store.mu.RLock()         // Lock for read access
	defer store.mu.RUnlock() // Unlock when done

	emp, ok := store.employees[id]
	return emp, ok
}

// UpdateEmployee updates the details of an existing employee in the in-memory store
func (store *EmployeeStore) UpdateEmployee(emp Employee) bool {
	store.mu.Lock()         // Lock for write access
	defer store.mu.Unlock() // Unlock when done

	_, ok := store.employees[emp.ID]
	if !ok {
		return false
	}
	store.employees[emp.ID] = emp
	return true
}

// DeleteEmployee deletes an employee from the in-memory store by ID
func (store *EmployeeStore) DeleteEmployee(id int) bool {
	store.mu.Lock()         // Lock for write access
	defer store.mu.Unlock() // Unlock when done

	_, ok := store.employees[id]
	if !ok {
		return false
	}
	delete(store.employees, id)
	return true
}

// ListEmployees returns a slice of all employees from the in-memory store
func (store *EmployeeStore) ListEmployees() []Employee {
	store.mu.RLock()         // Lock for read access
	defer store.mu.RUnlock() // Unlock when done

	var employees []Employee
	for _, emp := range store.employees {
		employees = append(employees, emp)
	}
	return employees
}

// ListEmployeesWithPagination retrieves employees with pagination support
func (store *EmployeeStore) ListEmployeesWithPagination(offset, limit int) []Employee {
	// Define a slice to store paginated employees
	var paginatedEmployees []Employee

	// Lock the mutex to ensure safe concurrent access
	store.mu.RLock()
	defer store.mu.RUnlock()

	// Iterate through all employees and append to the paginated slice if it falls within the requested range
	for _, emp := range store.employees {
		if offset > 0 {
			offset--
			continue
		}
		if limit == 0 {
			break
		}
		paginatedEmployees = append(paginatedEmployees, emp)
		limit--
	}

	return paginatedEmployees
}

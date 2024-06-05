package utils

import (
	"encoding/json"
	"net/http"
	"strconv"

	"employeedb/employee"
	"fmt"
	"sync"

	"github.com/go-resty/resty/v2"
)

// EmployeeHandler defines the interface for handling employee-related requests
type EmployeeHandler interface {
	HandleEmployeeRequest(http.ResponseWriter, *http.Request)
	HandleAllEmployeesRequest(http.ResponseWriter, *http.Request)
	HandleEmployeeCreation(w http.ResponseWriter, r *http.Request)
	HandleEmployeeUpdate(http.ResponseWriter, *http.Request)
	HandleEmployeeDeletion(http.ResponseWriter, *http.Request)
	HandleAllEmployeesRetrieval(http.ResponseWriter)
}

// EmployeeHandlerImpl implements the EmployeeHandler interface
type EmployeeHandlerImpl struct {
	Store *employee.EmployeeStore
	mu    sync.RWMutex // Mutex for concurrency safety
}

// NewEmployeeHandler creates a new instance of EmployeeHandlerImpl
func NewEmployeeHandler(store *employee.EmployeeStore) *EmployeeHandlerImpl {
	return &EmployeeHandlerImpl{
		Store: store,
	}
}

// HandleEmployeeRequest handles requests related to individual employees
func (eh *EmployeeHandlerImpl) HandleEmployeeRequest(store *employee.EmployeeStore) http.HandlerFunc {
	client := resty.New()
	eh.mu.Lock()         // Lock for write access
	defer eh.mu.Unlock() // Unlock when done

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			resp, err := client.R().Get("http://localhost:8080/employees/all")
			if err != nil {
				http.Error(w, "Failed to fetch employees", http.StatusInternalServerError)
				fmt.Println("error is : ", err)
				return
			}
			var employees []employee.Employee
			if err := json.Unmarshal(resp.Body(), &employees); err != nil {
				http.Error(w, "Failed to parse response body", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(employees)
		case http.MethodPost:
			eh.HandleEmployeeCreation(w, r)
		case http.MethodPut:
			eh.HandleEmployeeUpdate(store, w, r)
		case http.MethodDelete:
			eh.HandleEmployeeDeletion(store, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
// HandleAllEmployeesRequest handles requests related to all employees
func (eh *EmployeeHandlerImpl) HandleAllEmployeesRequest(store *employee.EmployeeStore) http.HandlerFunc {
	eh.mu.Lock()         // Lock for write access
	defer eh.mu.Unlock() // Unlock when done
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			eh.HandleAllEmployeesRetrieval(w,r) // Pass the http.ResponseWriter directly
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

// HandleEmployeeCreation handles employee creation requests
func (eh *EmployeeHandlerImpl) HandleEmployeeCreation(w http.ResponseWriter, r *http.Request) {
	var emp employee.Employee
	eh.mu.Lock()         // Lock for write access
	defer eh.mu.Unlock() // Unlock when done
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := eh.Store.CreateEmployee(emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

// HandleEmployeeUpdate handles employee update requests
func (eh *EmployeeHandlerImpl) HandleEmployeeUpdate(store *employee.EmployeeStore, w http.ResponseWriter, r *http.Request) {
	var emp employee.Employee
	eh.mu.Lock()         // Lock for write access
	defer eh.mu.Unlock() // Unlock when done
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if !store.UpdateEmployee(emp) {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emp)
}

// HandleEmployeeDeletion handles employee deletion requests
func (eh *EmployeeHandlerImpl) HandleEmployeeDeletion(store *employee.EmployeeStore, w http.ResponseWriter, r *http.Request) {
	var id int
	eh.mu.Lock()         // Lock for write access
	defer eh.mu.Unlock() // Unlock when done
	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		var err error
		id, err = strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid employee ID", http.StatusBadRequest)
			return
		}
	} else {
		var deleteRequest struct {
			ID int `json:"delete_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		id = deleteRequest.ID
	}

	if !store.DeleteEmployee(id) {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
// HandleAllEmployeesRetrieval handles requests to retrieve all employees with pagination
func (eh *EmployeeHandlerImpl) HandleAllEmployeesRetrieval(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters for pagination
    pageStr := r.URL.Query().Get("page")
    perPageStr := r.URL.Query().Get("perPage")

    // Convert page and perPage parameters to integers
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1 // Default to page 1 if invalid or not provided
    }
    perPage, err := strconv.Atoi(perPageStr)
    if err != nil || perPage < 1 {
        perPage = 10 // Default to 10 items per page if invalid or not provided
    }

    // Calculate offset based on pagination parameters
    offset := (page - 1) * perPage

    // Retrieve employees from the store with pagination
    employees := eh.Store.ListEmployeesWithPagination(offset, perPage)

    // Marshal the employees slice into a JSON object with custom formatting
    response, err := json.MarshalIndent(employees, "", "    ")
    if err != nil {
        http.Error(w, "Failed to format response", http.StatusInternalServerError)
        return
    }

    // Set the Content-Type header to application/json
    w.Header().Set("Content-Type", "application/json")

    // Write the response body to the client
    w.Write(response)
}

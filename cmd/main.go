package main

import (
	"fmt"
	"log"
	"net/http"

	"employeedb/employee"
	"employeedb/utils"
)

func main() {
	store := employee.NewEmployeeStore()
	employeeHandler := utils.NewEmployeeHandler(store)

	http.HandleFunc("/employees", employeeHandler.HandleEmployeeRequest(store))
	http.HandleFunc("/employees/all", employeeHandler.HandleAllEmployeesRequest(store))

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

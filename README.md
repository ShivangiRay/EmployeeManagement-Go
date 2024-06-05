# Employee Management System

This is a simple Employee Management System implemented in Go. It provides basic CRUD operations for managing employees using an in-memory data store.

## Installation

Clone the repository:

git clone path/to/repository

## Usage

1. Navigate to the repository 
2. go build cmd/main.go
3. ./main
4. The server will start running on http://localhost:8080.

## Endpoints

1. to get all employees - GET /employees/all

2. to create new employees - POST /employees
    Example request body: 
        {
            "id": 1,
            "name": "Shivangi Rai",
            "position": "Software Engineer",
            "salary": 100000
        }

3. to update an employee's information - PUT /employees
    Example request body:
        {
            "id": 1,
            "name": "Shivangi Rai",
            "position": "Software Engineer",
            "salary": 100000
        }

4. delete an employee - DELETE /employees?id=1

5. http://localhost:8080/employees/all?page=1&perPage=10



##  Dependencies 

1. github.com/go-resty/resty/v2: For making HTTP requests
2. github.com/stretchr/testify: For assertions in tests



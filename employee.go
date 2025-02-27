package main

type Employee struct {
    EmpID     int    `json:"empID"`
    DeptID    int    `json:"deptID"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}
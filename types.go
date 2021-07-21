package main

type Employee struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	Salary    int    `json:"salary"`
	Manager   int    `json:"manager"`
}

type Manager struct {
	Employee `json:"manager"`
	DirectReports map[int] Employee `json:"direct_reports"`
	EmployeeNames map[string]int
}

type Hierarchy struct {
	Manager *Manager
}

type Organization struct {
	totalSalary int
	employeeNames map[string]int
	employees map[int] Employee
	managers map[int] *Manager
	hierarchy Hierarchy
}

type OutputTree struct {
	Name string `json:"name"`
	Employees map[int]OutputTree `json:"employees,omitempty"`
}

type OutputStruct struct {
	OutputTree
	TotalSalary int `json:"total_salary"`
}

/**
Output

{
	"name": "Jeff",
	"employees": [
		{
			"name": "Dave",
			"employees": [
				{"name": "Andy"},

		    ]
		}
	]
}
 */
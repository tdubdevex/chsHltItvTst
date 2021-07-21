package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	// Read contents of the file provided by the file path
	contents, err := os.ReadFile("./input.json")
	if err != nil {
		panic(err)
	}

	employees := parseEmployees(contents)

	org := Organization{
		totalSalary:   0,
		employeeNames: make(map[string]int, len(employees)),
		employees:     make(map[int]Employee, len(employees)),
		managers:      map[int]*Manager{},
		hierarchy: Hierarchy{
			Manager: &Manager{
				Employee:      Employee{},
				DirectReports: map[int]Employee{},
				EmployeeNames: map[string]int{},
			},
		},
	}

	populateEmployees(employees, &org)

	sortEmployees(&org)

	printManagerTree(&org)
}

// printManagerTree Print out organization tree
func printManagerTree(organization *Organization) {
	employee := organization.hierarchy.Manager.Employee

	manager, ok := organization.managers[employee.Id]

	fmt.Println(employee.FirstName)
	if ok {
		fmt.Printf("Employees of: %s\n", employee.FirstName)
		printEmployees(
			manager.DirectReports,
			organization,
		)
	}

	fmt.Printf("Total Salary: %d\n", organization.totalSalary)
}

// printEmployees print to stdOut the employee and subordinates
func printEmployees(employees map[int]Employee, organization *Organization) {
	for idx, employee := range employees {
		fmt.Println(employees[idx].FirstName)

		if _, ok := organization.managers[employee.Id]; ok {
			fmt.Printf("Employees Of: %s\n", employee.FirstName)

			printEmployees(
				organization.managers[employee.Id].DirectReports,
				organization,
			)
		}
	}
}

// populateEmployees Add employees to a lookup map by employee.ID
func populateEmployees(employees []Employee, org *Organization) {
	// First loop to map employees for quick lookups
	for _, val := range employees {
		//fmt.Printf("Populating array with %s\n", val.FirstName)
		// Sum salary
		(*org).totalSalary += val.Salary
		(*org).employees[val.Id] = val

		// If this employee has no manager they are at the top of the hierarchy
		if val.Manager == 0 {
			org.managers[val.Id] = org.hierarchy.Manager
			org.hierarchy.Manager.Employee = val
		}
	}
}

// sortEmployees categorically sort employees by name and by manager
func sortEmployees(org *Organization) {
	// Loop for the second time to sort employees into direct
	// reports for hierarchy
	//fmt.Printf("%v", org.employees)
	for _, val := range (*org).employees {
		//fmt.Printf("org employee %d %s\n", idx, (*org).employees[idx].FirstName)

		populateNames(org, &val)

		if val.Manager != 0 {
			_, ok := org.managers[val.Manager]

			if !ok {
				org.managers[val.Manager] = &Manager{
					Employee:      org.employees[val.Manager],
					DirectReports: map[int]Employee{},
					EmployeeNames: map[string]int{},
				}
			}
			org.managers[val.Manager].DirectReports[val.Id] = val
			org.managers[val.Manager].EmployeeNames[val.FirstName] = val.Id
		}
	}
}

func populateNames(org *Organization, val *Employee) {
	// Go sorts string maps alphabetically.
	_, alreadyExists := (*org).employeeNames[(*val).FirstName]

	if alreadyExists {
		//fmt.Printf("Found name %s\n", (*val).FirstName)
		// Not a huge fan of recursion,
		// but this was the quickest working solution of dealing with duplicates
		populateNames(org, &Employee{
			Id:        val.Id,
			FirstName: fmt.Sprintf("%s ", (*val).FirstName), // add a space
			Salary:    val.Salary,
			Manager:   val.Manager,
		})
	} else {
		//fmt.Printf("Adding name to string map, %s.\n", (*val).FirstName)
		(*org).employeeNames[val.FirstName] = (*val).Id
	}
}

func parseEmployees(contents []byte) (employees []Employee) {
	// Force input json to adhere to the inputJson struct type declarations
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&employees)

	if err != nil {
		panic(err)
	}

	return
}

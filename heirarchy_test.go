package main

import (
	"reflect"
	"testing"
)

func Test_parseEmployees(t *testing.T) {
	type args struct {
		contents []byte
	}
	tests := []struct {
		name          string
		args          args
		wantEmployees []Employee
	}{
		{
			"happy path",
			args{
				contents: []byte(`[{"id": 1,"first_name": "Dave","manager": null,"salary": 100000}]`),
			},
			[]Employee{
				{
					Id:        1,
					FirstName: "Dave",
					Salary:    100000,
					Manager:   0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEmployees := parseEmployees(tt.args.contents); !reflect.DeepEqual(gotEmployees, tt.wantEmployees) {
				t.Errorf("parseEmployees() = %v, want %v", gotEmployees, tt.wantEmployees)
			}
		})
	}
}

func Test_populateEmployees(t *testing.T) {
	type args struct {
		employees []Employee
		org       *Organization
		wantSalary int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"happy path",
			args{
				employees: []Employee{
					{
						Id:        2,
						FirstName: "Jeff",
						Salary:    110000,
						Manager:   0,
					},
					{
						Id:        1,
						FirstName: "Jeff",
						Salary:    100000,
						Manager:   2,
					},
				},
				org:       &Organization{
					totalSalary:   0,
					employeeNames: make(map[string]int, 2),
					employees:     make(map[int]Employee, 2),
					managers:      make(map[int]*Manager, 1),
					hierarchy:     Hierarchy{
						Manager: &Manager{
							Employee:      Employee{},
							DirectReports: map[int]Employee{},
						},
					},
				},
				wantSalary: 210000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			populateEmployees(tt.args.employees, tt.args.org)

			if tt.args.wantSalary != tt.args.org.totalSalary {
				t.Errorf(
					"Salary does not equal want %d, got %d",
					tt.args.wantSalary,
					tt.args.org.totalSalary,
				)
			}
		})
	}
}

func Test_populateNames(t *testing.T) {
	type args struct {
		org *Organization
		val *Employee
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy path",
			args: args{
				org: &Organization{
					totalSalary:   0,
					employeeNames: map[string]int {
						"Dave": 0,
					},
					employees: 	   map[int]Employee{},
					managers:      nil,
					hierarchy:     Hierarchy{},
				},
				val: &Employee{
					Id:        1,
					FirstName: "Dave",
					Salary:    100000,
					Manager:   2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			populateNames(tt.args.org, tt.args.val)

			if _, ok := tt.args.org.employeeNames[tt.args.val.FirstName]; !ok {
				t.Errorf("employee name is not populated in map %v", tt.args.val.FirstName)
			}
		})
	}
}

func Test_sortEmployees(t *testing.T) {
	type args struct {
		org *Organization
	}

	jeff := Employee{
		Id:        2,
		FirstName: "Jeff",
		Salary:    110000,
		Manager:   0,
	}
	dave := Employee{
		Id:        1,
		FirstName: "Dave",
		Salary:    100000,
		Manager:   2,
	}
	expectedManager := Manager{
		Employee: jeff,
		DirectReports: map[int]Employee{
			1: dave,
		},
		EmployeeNames: map[string]int{
			"Dave": 1,
		},
	}

	tests := []struct {
		name string
		args args
		want args
	}{
		{
			name: "happy path",
			args: args{
				org: &Organization{
					totalSalary:   210000,
					employeeNames: map[string]int{},
					employees: map[int]Employee{
						dave.Id: dave,
						jeff.Id: jeff,
					},
					managers: map[int]*Manager{},
					hierarchy:     Hierarchy{Manager: &Manager{
						Employee:      Employee{},
						DirectReports: nil,
					}},
				},
			},
			want: args{
				org: &Organization{
					totalSalary:   210000,
					employeeNames: map[string]int{
						"Dave": 1,
						"Jeff": 2,
					},
					employees: map[int]Employee{
						1: {
							Id:        1,
							FirstName: "Dave",
							Salary:    100000,
							Manager:   2,
						},
						2: {
							Id:        2,
							FirstName: "Jeff",
							Salary:    110000,
							Manager:   0,
						},
					},
					managers: map[int]*Manager{
						2: &expectedManager,
					},
					hierarchy:     Hierarchy{
						&expectedManager,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortEmployees(tt.args.org)
			if !reflect.DeepEqual(tt.args.org.employeeNames, tt.want.org.employeeNames) {
				t.Errorf("employeeNames = %v, want %v", tt.args.org.employeeNames, tt.want.org.employeeNames)
			}
			if !reflect.DeepEqual(tt.args.org.managers, tt.want.org.managers) {
				t.Errorf("managers = %v, want %v", tt.args.org.managers, tt.want.org.managers)
			}
		})
	}
}

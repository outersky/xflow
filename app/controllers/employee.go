package controllers

import (
    "github.com/robfig/revel"
	"github.com/outersky/xflow/app/models"
	_ "github.com/outersky/xflow/app/routes"
)

type Dept struct {
    AuthedApp
}

type Employee struct {
    AuthedApp
}

type User struct {
    AuthedApp
}

/*
func loadCompany() []interface{} {
    companies,_ = c.Txn.Select(models.Company{}, `select * from Company `)
    return companies
}
*/

func (c Dept) List () revel.Result {
    depts, _ := c.Txn.Select(models.Dept{}, `select * from Dept where CompanyId = ? `, c.CompanyId)
	return c.RenderJson(List(depts))
}

func (c Dept) ListOfCompany (companyId string) revel.Result {
    depts, _ := c.Txn.Select(models.Dept{}, `select * from Dept where CompanyId = ? `, companyId)
	return c.RenderJson(List(depts))
}

func (c Employee) List (deptId string) revel.Result {
	revel.INFO.Println("Employee.List")
    employees,_ := c.Txn.Select(models.Employee{}, `select * from Employee where DeptId = ? `, deptId)
	return c.RenderJson(List(employees))
}

func (c User) Add(user models.User, verifyPassword string, company models.Company) revel.Result {
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).  Message("Password does not match")
	user.Validate(c.Validation)
	c.Validation.Required(company.Name)
	c.Validation.Required(company.Domain)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
        return c.RenderJson(Error(c.Validation))
	}

	user.HashedPassword = encode(passwd(user.Password))
    company.Id = models.GenId(models.TCompany)
    company.CompanyId = company.Id
    user.Id = models.GenId(models.TUser)
    user.CompanyId = company.Id
	err := c.Txn.Insert(&user, &company)
	if err != nil {
		panic(err)
	}

    user.Password = ""
    user.HashedPassword = ""
    return c.RenderJson(Single(user))
}

func (c Dept) Add(dept models.Dept, companyId string) revel.Result {
	c.Validation.Required(dept.Name)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		return c.RenderJson(Error(c.Validation))
	}

    dept.Id = models.GenId(models.TDept)
    dept.CompanyId = companyId
	err := c.Txn.Insert(&dept)
	if err != nil {
		panic(err)
	}
	return c.RenderJson(Single(dept))
}

func (c Employee) Add(employee models.Employee, deptId string) revel.Result {
	c.Validation.Required(employee.Name)
	c.Validation.Required(employee.Email)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		return c.RenderJson(Error(c.Validation))
	}

    employee.Id = models.GenId(models.TEmployee)
    employee.DeptId = deptId
    employee.CompanyId = c.CompanyId
	err := c.Txn.Insert(&employee)
	if err != nil {
		panic(err)
	}
	return c.RenderJson(Single(employee))
}

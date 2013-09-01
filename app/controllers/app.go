package controllers

import (
    "fmt"
	"code.google.com/p/go.crypto/bcrypt"
    "github.com/robfig/revel"
	"github.com/outersky/xflow/app/models"
	"github.com/outersky/xflow/app/routes"
)

type App struct {
    GorpController
}

func (c App) Index() revel.Result {
	revel.INFO.Println("App.Index")
	var companies []*models.Company
    companies = loadCompany(c.Txn.Select(models.Company{},
        `select * from Company `))
	return c.Render(companies)
}

func (c App) AddUser() revel.Result {
	fmt.Printf("... App.AddUser() .\n")

	if user := c.connected(); user != nil {
		c.RenderArgs["user"] = user
	}
	return nil
}

func (c App) connected() *models.User {
	if c.RenderArgs["user"] != nil {
		return c.RenderArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username)
	}
	return nil
}

func (c App) getUser(username string) *models.User {
	users, err := c.Txn.Select(models.User{}, `select * from User where Username = ?`, username)
	if err != nil {
		panic(err)
	}
	if len(users) == 0 {
		return nil
	}
	return users[0].(*models.User)
}

func (c App) Login(username, password string) revel.Result {
	user := c.getUser(username)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(decode(user.HashedPassword), []byte(password))
		if err == nil {
			c.Session["user"] = username
			c.Flash.Success("Welcome, " + username)
			return c.Redirect(routes.App.Index())
		}
	}

	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed")
	return c.Redirect(routes.App.Index())
}

func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Index())
}


func (c App) SaveUser(user models.User, verifyPassword string, company models.Company) revel.Result {
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).  Message("Password does not match")
	user.Validate(c.Validation)
	c.Validation.Required(company.Name)
	c.Validation.Required(company.Domain)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Index())
	}

    //pwd, _ := bcrypt.GenerateFromPassword( []byte(user.Password), bcrypt.DefaultCost)
	user.HashedPassword = encode(passwd(user.Password))
    company.Id = models.GenId(models.TCompany)
    company.CompanyId = company.Id
    user.Id = models.GenId(models.TUser)
    user.CompanyId = company.Id
	err := c.Txn.Insert(&user, &company)
	if err != nil {
		panic(err)
	}

	c.Session["user"] = user.Username
	c.Flash.Success("Welcome, " + user.Name)
	return c.Redirect(routes.App.Index())
}

func loadCompany(results []interface{}, err error) []*models.Company{
	if err != nil {
		panic(err)
	}
	var companies []*models.Company
	for _, r := range results {
		companies = append(companies, r.(*models.Company))
	}
	return companies
}

func (c App) SaveDept(dept models.Dept, companyId string) revel.Result {
	c.Validation.Required(dept.Name)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Index())
	}

    dept.Id = models.GenId(models.TDept)
    dept.CompanyId = companyId
	err := c.Txn.Insert(&dept)
	if err != nil {
		panic(err)
	}
	return c.Redirect(routes.App.Index())
}

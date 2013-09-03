package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/outersky/xflow/app/models"
	"github.com/robfig/revel"
)

type App struct {
	GorpController
	UserId    string
    UserName  string
	CompanyId string
	Roles     string
}

type AuthedApp struct {
	App
}

func (c *App) Index() revel.Result {
	revel.INFO.Println("App.Index")
	companies, _ := c.Txn.Select(models.Company{}, `select * from Company `)
	depts, _ := c.Txn.Select(models.Dept{}, `select * from Dept`)
	return c.Render(companies, depts)
}

func (c *App) Services() revel.Result {
    return c.Render()
}

func (c *App) AddUser() revel.Result {
	fmt.Printf("... App.AddUser() .\n")

	c.UserName = c.loadUserName()
	c.UserId = c.loadUserId()
	c.CompanyId = c.loadCompanyId()

	return nil
}

func (c *App) Current() revel.Result{
    if(c.UserId==""){
        return c.RenderJson(Error("NotLogged"))
    }else{
        m := map[string]string{
            "UserName":c.UserName,
            "UserId":c.UserId,
            "CompanyId":c.CompanyId,
        }
        return c.RenderJson(Single(m))
    }
}

func (c *App) loadUserName() string {
	if userName, ok := c.Session["UserName"]; ok {
		fmt.Printf(" UserName loaded : %s\n", userName)
		return userName
	}
	return ""
}

func (c *App) loadUserId() string {
	if userId, ok := c.Session["UserId"]; ok {
		fmt.Printf(" UserId loaded : %s\n", userId)
		return userId
	}
	return ""
}

func (c *App) loadCompanyId() string {
	if companyId, ok := c.Session["CompanyId"]; ok {
		fmt.Printf(" CompanyId loaded : %s\n", companyId)
		return companyId
	}
	return ""
}

func (c *AuthedApp) CheckAuth() revel.Result {
	fmt.Printf("... AuthedApp.CheckAuth() : %s .\n", c.UserId)
	if c.UserId == "" {
		return c.RenderJson(Error("Not Logged"))
	}
	return nil
}

/*
func (c App) connected() *models.User {
	if c.RenderArgs["user"] != nil {
		return c.RenderArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username)
	}
	return nil
}

*/

func (c *App) getUser(username string) *models.User {
	users, err := c.Txn.Select(models.User{}, `select * from User where Username = ?`, username)
	if err != nil {
		panic(err)
	}
	if len(users) == 0 {
		return nil
	}
	return users[0].(*models.User)
}

func (c *App) Login(username, password string) revel.Result {
	user := c.getUser(username)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(decode(user.HashedPassword), []byte(password))
		if err == nil {
			c.Session["UserName"] = username
			c.Session["UserId"] = user.Id
			c.Session["CompanyId"] = user.CompanyId
			return c.RenderJson(Single(user))
		}
	}

	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed")
	return c.RenderJson(Error("Login ERROR"))
}

func (c *App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.RenderJson(Single("OK"))
}

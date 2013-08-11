package models

import (
	"fmt"
	"strconv"
	"github.com/robfig/revel"
	"regexp"
)

const (
	EUsr = 100 + iota
	ECompany
)


func GenId(t int) string {
	return ""+strconv.Itoa(t)+"123" //+ random()
}

type Base struct {
	Id string
	CompanyId string
	Created  int64
	Updated  int64
	Version  int64
	Status   int64

	GenId func() string
}

type Usr struct {
	Id string
	CompanyId string
	Created  int64
	Updated  int64
	Version  int64
	Status   int64

	Name               string
	Username, Password string
	HashedPassword     string
}

type Company struct {
	Id string
	CompanyId string
	Created  int64
	Updated  int64
	Version  int64
	Status   int64

	Name        string
	Code 	    string
	HomePage    string
	Domain		string
}

func (c *Company) GenId() string {
	c.Id=GenId(ECompany)
	return c.Id
}

func (u *Usr) GenId() string {
	u.Id=GenId(EUsr)
	return u.Id
}

func (u *Usr) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}

var userRegex = regexp.MustCompile("^\\w*$")

func (user *Usr) Validate(v *revel.Validation) {
	v.Check(user.Username,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{4},
		revel.Match{userRegex},
	)

	ValidatePassword(v, user.Password).
		Key("user.Password")

	v.Check(user.Name,
		revel.Required{},
		revel.MaxSize{100},
	)
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{5},
	)
}

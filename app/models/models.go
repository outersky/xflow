package models

import (
	"fmt"
	"strconv"
	"math/rand"
	"github.com/robfig/revel"
	"regexp"
)

type EntityType int

const (
	TUser EntityType = 100 + iota
	TCompany
	TDept
)

var r *rand.Rand = rand.New(rand.NewSource(99))

func GenId(t EntityType) string {
	return ""+strconv.Itoa(int(t)) + strconv.Itoa(r.Int())
}

type Base struct {
	Id string
	CompanyId string
	Created  int64
	Updated  int64
	Version  int64
	Status   int64

//	GenId func() string
}

type Company struct {
	Base
	Name        string
	Code	    string
	HomePage    string
	Domain		string
}

type User struct {
	Base
	Name               string
	Username, Password string
	HashedPassword     string
	StaffId	string
}

type Dept struct {
	Base
	Name string
	Code string
	ParentDeptId string
}

type Staff struct {
	Base
	Name string
	Code string
	Email string
	DeptId string
}

type Application struct {
	Base
	FormId string
	UserId  string
}

type Applicant struct {

}

type Path struct {
	Name string
}

type Node struct {
	PathId string
	PrevNodeId string
}

type Step struct {
	NodeId string

}

type Form struct {
	Name string
	PathId string
}

type FormField struct {
	FormId string
	FieldType string
	Content string
	Col int
	Row int
}

type Category struct {

}

type Item struct {

}

type Model interface {
	GenId() string
	PreInsert()
	PreSave()
}

func (c *Company) GenId() string {
	c.Id=GenId(TCompany)
	return c.Id
}

func (u *User) GenId() string {
	u.Id=GenId(TUser)
	return u.Id
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}

var userRegex = regexp.MustCompile("^\\w*$")

func (user *User) Validate(v *revel.Validation) {
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

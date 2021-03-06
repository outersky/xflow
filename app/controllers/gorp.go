package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/ziutek/mymysql/godrv"
	"reflect"

	// _ "github.com/mattn/go-sqlite3"
	"github.com/outersky/xflow/app/models"
	r "github.com/robfig/revel"
	"github.com/robfig/revel/modules/db/app"
)

var (
	Dbm *gorp.DbMap
)

func Dbg() {
	t := reflect.TypeOf(models.User{})
	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)
		fmt.Printf("Field Name: %s, type: %s \n", f.Name, f.Type.Name())
	}
}
func encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func decode(src string) []byte {
	data, _ := base64.StdEncoding.DecodeString(src)
	return data
}

func passwd(src string) []byte {
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(src), bcrypt.DefaultCost)
	return bcryptPassword
}

func Init() {
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
		for col, size := range colSizes {
			t.ColMap(col).MaxSize = size
		}
	}

	//Dbg()

	if false {
		return
	}

	t := Dbm.AddTable(models.User{}).SetKeys(false, "Id")
	t.ColMap("Password").Transient = true
	setColumnSizes(t, map[string]int{
		"Username": 20,
		"Name":     100,
	})

	t = Dbm.AddTable(models.Dept{}).SetKeys(false, "Id")

	t = Dbm.AddTable(models.Employee{}).SetKeys(false, "Id")

	t = Dbm.AddTable(models.Company{}).SetKeys(false, "Id")
	setColumnSizes(t, map[string]int{
		"Name": 50,
	})

	Dbm.TraceOn("[gorp]", r.INFO)
	Dbm.DropTables()
	Dbm.CreateTables()

	//	bcryptPassword, _ := bcrypt.GenerateFromPassword(
	//		[]byte("demo"), bcrypt.DefaultCost)

	company := &models.Company{}
	company.Id = models.GenId(models.TCompany)
	company.CompanyId = company.Id
	company.Name = "Demo Company"
	company.Domain = "demo.com"

	company2 := &models.Company{}
	company2.Id = models.GenId(models.TCompany)
	company2.CompanyId = company2.Id
	company2.Name = "Demo2 Company"
	company2.Domain = "demo2.com"

	dept := &models.Dept{}
	dept.Id = models.GenId(models.TDept)
	dept.CompanyId = company.Id
	dept.Name = "HR"

	dept2 := &models.Dept{}
	dept2.Id = models.GenId(models.TDept)
	dept2.CompanyId = company2.Id
	dept2.Name = "IT"

	employee := &models.Employee{}
	employee.Id = models.GenId(models.TEmployee)
	employee.CompanyId = company.Id
	employee.DeptId = dept.Id
	employee.Name = "John"
	employee.Email = "john@demo.com"

	employee2 := &models.Employee{}
	employee2.Id = models.GenId(models.TEmployee)
	employee2.CompanyId = company2.Id
	employee2.DeptId = dept2.Id
	employee2.Name = "Rose"
	employee2.Email = "rose@demo.com"

	demoUser := &models.User{models.Base{Id: models.GenId(models.TUser)}, "Demo User", "demo", "demo", encode(passwd("demo")), employee.Id}
	demoUser.CompanyId = company.Id

	demoUser2 := &models.User{models.Base{Id: models.GenId(models.TUser)}, "Demo2 User", "demo2", "demo2", encode(passwd("demo2")), employee2.Id}
	demoUser2.CompanyId = company2.Id

	if err := Dbm.Insert(company, company2, dept, dept2, employee, employee2, demoUser, demoUser2); err != nil {
		panic(err)
	}
	/*
		hotels := []*models.Hotel{
			&models.Hotel{0, "Marriott Courtyard", "Tower Pl, Buckhead", "Atlanta", "GA", "30305", "USA", 120},
			&models.Hotel{0, "W Hotel", "Union Square, Manhattan", "New York", "NY", "10011", "USA", 450},
			&models.Hotel{0, "Hotel Rouge", "1315 16th St NW", "Washington", "DC", "20036", "USA", 250},
		}
		for _, hotel := range hotels {
			if err := Dbm.Insert(hotel); err != nil {
				panic(err)
			}
		}*/
	Dbm.TraceOff()
}

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

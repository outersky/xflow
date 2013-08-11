package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"github.com/coopernurse/gorp"
    _ "github.com/ziutek/mymysql/godrv"

	// _ "github.com/mattn/go-sqlite3"
	r "github.com/robfig/revel"
	"github.com/robfig/revel/modules/db/app"
	"github.com/outersky/xflow/app/models"
)

var (
	Dbm *gorp.DbMap
)

func Init() {
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
		for col, size := range colSizes {
			t.ColMap(col).MaxSize = size
		}
	}

	t := Dbm.AddTable(models.Usr{}).SetKeys(false, "Id")
	t.ColMap("Password").Transient = true
	setColumnSizes(t, map[string]int{
		"Username": 20,
		"Name":     100,
	})

	t = Dbm.AddTable(models.Company{}).SetKeys(false, "Id")
	setColumnSizes(t, map[string]int{
		"Name":    50,
	})

	Dbm.TraceOn("[gorp]", r.INFO)
	Dbm.DropTables()
	Dbm.CreateTables()

	bcryptPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("demo"), bcrypt.DefaultCost)

	demoUser := &models.Usr{Id:models.GenId(models.EUsr),Name:"Demo User", Username:"demo", Password:"demo", HashedPassword:string(bcryptPassword)}
	if err := Dbm.Insert(demoUser); err != nil {
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

package controllers
/*
import (
    "github.com/robfig/revel"
	"github.com/outersky/xflow/app/models"
	_ "github.com/outersky/xflow/app/routes"
)

type Dept struct {
    GorpController
}

func (c Dept) List (companyId string) revel.Result {
	revel.INFO.Println("Dept.List")
    depts, _ := c.Txn.Select(models.Dept{},
        `select * from Dept where CompanyId = ? `, companyId)
	return c.RenderJson(List(depts))
}

func loadDept(results []interface{}, err error) []*models.Dept{
	if err != nil {
		panic(err)
	}
	var depts []*models.Dept
	for _, r := range results {
		depts = append(depts, r.(*models.Dept))
	}
	return depts
}
*/

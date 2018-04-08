package accounts

import (
	"github.com/aosfather/bingo"
	"github.com/aosfather/bingo/utils"
	"github.com/aosfather/bingo/sql"
)

/*
  账户管理
   1、创建账户
   2、获取账户信息
   3、获取账户流水
   4、更新账户状态和名称
 */
type AccountService struct {
	logger utils.Log
	dao *bingo.BaseDao
}
func (this *AccountService)Init(c *bingo.ApplicationContext) {
	this.logger=c.GetLog("account_service")
	this.dao=c.CreateDao()
}

func (this *AccountService)CreateAccount(id,namespace,owner,label string)error  {
       account:=Account{}
       account.Type=namespace
       account.Id=id
       account.Owner=owner
       account.Label=label

      _,err:=this.dao.Insert(&account)
      return err
}

func (this *AccountService)FindAccounts(owner string) []*Account {
	if owner=="" {
		return nil
	}
	account:=Account{}
	account.Owner=owner
	result:=this.dao.QueryAll(&account,"Owner")
	return toAccounts(result)
}

func toAccounts(source []interface{})[]*Account {
	if source!=nil &&len(source)>0 {
		var datas []*Account
		for _,item:=range source {
			datas=append(datas,item.(*Account))
		}
		return datas

	}
	return nil
}

func (this *AccountService)FindAccountsByType(owner string,namespace string)[]*Account {
	if owner=="" ||namespace==""{
		return nil
	}
	account:=Account{}
	account.Owner=owner
	account.Type=namespace
	result:=this.dao.QueryAll(&account,"Owner","Type")
	return toAccounts(result)
}

func (this *AccountService)FindAccount(id string)*Account {
	if id=="" {
		return nil
	}
	account:=Account{}
	account.Id=id
	if this.dao.Find(&account,"Id"){
		return &account
	}
	return nil
}

func (this *AccountService)UpdateAccountLabel(id,label string)bool {
	if id==""||label==""{
		return false
	}
	account:=Account{}
	account.Id=id
	account.Label=label
	count,err:=this.dao.Update(&account,"Label")
	if count==1&&err==nil {
		return true
	}
	return false
}

func (this *AccountService)UpdateAccountStatus(id string,status AccountStatus)bool {
	if id=="" {
		return false
	}
	account:=Account{}
	account.Id=id
	account.Status=status
	count,err:=this.dao.Update(&account,"Label")
	if count==1&&err==nil {
		return true
	}
	return false
}

//获取账户流水
func(this *AccountService)GetFlowWater(id string,pageindex,pagesize int)[]*FlowWater {
	if id=="" {
		return nil
	}

	if pagesize<=0 {
		pagesize=20
	}

	flow:=FlowWater{}
	flow.Account=id
	result:=this.dao.Query(&flow,sql.Page{pagesize,pageindex,0},"Account")
	if result!=nil &&len(result)>0 {
		var datas []*FlowWater
		for _,item:=range result {
			datas=append(datas,item.(*FlowWater))
		}
		return datas

	}
	return nil
}


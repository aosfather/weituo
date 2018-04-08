package accounts

import (
 "github.com/aosfather/bingo/sql"
)

/**
交易
 */


func createTransaction(theType TransactionType,account,toaccount string,code,id string,amount int64,detail string)*Transaction {
 t:=Transaction{}
 if theType==T_IN {
  t.InAccount = account
 }else if theType==T_OUT {
  t.OutAccount=account
 }else if theType==T_TRANSFER {
  t.InAccount=toaccount
  t.OutAccount=account
 }
 t.Type=theType
 t.Code=code
 t.Request=id
 t.Amount=amount
 t.Detail=detail
 return &t
}

func createFlow(ftype FlowType,account string,code string,amount int64,detail string)*FlowWater{
 flow:=FlowWater{}
 flow.Account=account
 flow.Amount=amount
 flow.Type=ftype
 flow.SubType=code
 flow.Detail=detail
 return &flow
}

//入账交易
func (this *AccountService)DoInTransaction(account string,code,id string,amount int64,detail string)bool {
 if amount<=0 { //发生额必须大于零
  return false
 }
   t:=createTransaction(T_IN,account,"",code,id,amount,detail)
   count,err:=this.dao.Insert(&t)
   //插入失败
   if count!=1||err!=nil {

      return false
   }

   //插入流水并更新交易表的入账流水号
   flow:=createFlow(FLOW_IN,account,code,amount,detail)
   session:=this.dao.GetSession()
   session.Begin()
   flowId:=this.insertFlow(session,flow)
   if flowId!=0 {
       t.InFlow=flowId
       t.Version=t.Version+1
       _,count,err=session.ExeSql("update wt_acc_transactions set inflow=?,version=? where code=? and request=? and version=?",flowId,t.Version+1,t.Code,t.Request,t.Version)
       if count==1&&err!=nil {
        session.Commit()
        return true
       }


   }

   session.Rollback()
   return false

 }

 //出账交易
func (this *AccountService)DoOutTransaction(account string,code,id string,amount int64,detail string)bool {
 if amount<=0 { //发生额必须大于零
  return false
 }
 t:=createTransaction(T_OUT,account,"",code,id,amount,detail)
 count,err:=this.dao.Insert(&t)
 //插入失败
 if count!=1||err!=nil {

  return false
 }

 //插入流水并更新交易表的入账流水号
 flow:=createFlow(FLOW_OUT,account,code,amount,detail)
 session:=this.dao.GetSession()
 session.Begin()
 flowId:=this.insertFlow(session,flow)
 if flowId!=0 {
  t.InFlow=flowId
  t.Version=t.Version+1
  _,count,err=session.ExeSql("update wt_acc_transactions set outflow=?,version=? where code=? and request=? and version=?",flowId,t.Version+1,t.Code,t.Request,t.Version)
  if count==1&&err!=nil {
   session.Commit()
   return true
  }


 }

 session.Rollback()
 return  false
}

//转账交易
func (this *AccountService)DoTrasferTransaction(from,to string,code,id string,amount int64,detail string)bool {
 if amount<=0 { //发生额必须大于零
  return false
 }
 t:=createTransaction(T_TRANSFER,from,to,code,id,amount,detail)
 count,err:=this.dao.Insert(&t)
 //插入失败
 if count!=1||err!=nil {

  return false
 }

 //插入流水并更新交易表的入账流水号T
 outflow:=createFlow(FLOW_TOUT,from,code,amount,detail)
 inflow:=createFlow(FLOW_TIN,to,code,amount,detail)
 session:=this.dao.GetSession()
 session.Begin()
 flowId:=this.insertFlow(session,outflow)
 if flowId!=0 {
  t.OutFlow=flowId
  inFlowId:=this.insertFlow(session,inflow)
  if inFlowId!=0 {
   t.InFlow=inFlowId
   _,count,err=session.ExeSql("update wt_acc_transactions set outflow=?,inflow=?,version=? where code=? and request=? and version=?",flowId,inFlowId,t.Version+1,t.Code,t.Request,t.Version)
   if count==1&&err!=nil {
    session.Commit()
    return true
   }
  }
 }

 session.Rollback()
 return false
}

//插入流水
func (this *AccountService)insertFlow(session *sql.TxSession,flow *FlowWater) int64{
 if flow.Amount<=0 { //流水金额必须大于零
  return 0
 }
 //根据流水类型，进行账户余额判断
 acc:=Account{}
 acc.Id=flow.Account

 exist:=session.Find(&acc,"Id")
 if !exist {
     return 0
 }

 flow.Begin=acc.Amount
 var sum int64
 var sumField string
 //余额计算
 switch flow.Type {
 case FLOW_IN:
  sumField="in_sum"
  sum=acc.InSum+flow.Amount
  acc.Amount=acc.Amount+flow.Amount
 case FLOW_TIN:
  sumField="tin_sum"
  sum=acc.TinSum+flow.Amount
  acc.Amount=acc.Amount+flow.Amount
 case FLOW_TOUT:
    sumField="tout_sum"
    acc.Amount=acc.Amount-flow.Amount
    sum=acc.ToutSum+flow.Amount
 case FLOW_OUT:
     sumField="out_sum"
     acc.Amount=acc.Amount-flow.Amount
     sum=acc.OutSum+flow.Amount

 }

 //余额不足
 if acc.Amount<0 {
   return 0
 }

 flow.End=acc.Amount
 //更新余额
 _,count,err:=session.ExeSql("update wt_acc_accounts set amount=?,"+sumField+"=? ,version=? where id=? and version=?",acc.Amount,sum,acc.Version,acc.Id,acc.Version+1)
 if count==1&&err==nil {//插入流水，并返回流水号
  flowId,_,_:=session.Insert(flow)
   return flowId
  }

   return 0
}
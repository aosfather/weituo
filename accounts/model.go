package accounts

/**
基本模型
1、账户模型
2、账户流水模型
3、账户交易模型
交易产生流水，流水对账户可用余额进行操作

账户：唯一编码code、类型、可用余额、入账总额、出账总额、转出总额、转入总额、归属标识
     1）其中同一种类型的账户之间才能进行
     2）入账：从该类型账户体系外进入
     3）出账：从该类型账户体系进入另外的体系，如余额账户提现
     4）转出：同类型账户体系中的流出
     5）转入：同类型账户系统中的流入
     6）（入账+转入）-（出账+转出）=余额

账户流水：flow、账户、发生额、类型、子类型、期初、期末、摘要说明

交易：业务类型、交易id、交易类型、出账账户、入账账户、出账流水号、入账流水号、交易说明
 */

type FlowType int //流水类型
type AccountStatus byte //账户状态
type  TransactionType byte //交易类型

	const(
	//----交易类型----------//
	T_IN TransactionType=1       //入账
	T_TRANSFER TransactionType=2 //转账
	T_OUT TransactionType=3      //出账

	//----流水类型----------//
	FLOW_IN FlowType=11  //入账
	FLOW_TIN FlowType=12 //转入
	FLOW_OUT FlowType=21 //出账
	FLOW_TOUT FlowType=22//转出
	//---------------------//
	AS_ACTIVE AccountStatus=1 //有效状态
	AS_NO_OUT AccountStatus=2 //禁止出账
	AS_NO_TOUT AccountStatus=3//禁止转出
	AS_FROZEN_OUT AccountStatus=4 //冻结流出状态，禁止所有出
	AS_NO_IN AccountStatus=5 //禁止入账
	AS_NO_TIN AccountStatus=6 //禁止转入
	AS_FROZEN_IN AccountStatus=7 //冻结流入状态，禁止所有入
	AS_FROZEN AccountStatus=10 //冻结状态，禁止所有入和所有出
	AS_CLOSED AccountStatus=110 //账户处于被注销关闭状态，属于不在使用。

)

//账户
type Account struct {
	Id string  `json:"id" Table:"wt_acc_accounts"`//唯一编码
	Type string `json:"style"`//账户类型
	Label string `json:"label"`//账户名称
	Owner string `json:"owner"`//归属
	Amount int64  `json:"amount"`//余额
	InSum int64 `json:"in_sum"`//转入总额
	OutSum int64 `json:"out_sum"`//转出总额
	TinSum int64 `json:"tin_sum"`//转入总额
	ToutSum int64`json:"tout_sum"` //转出总额
	CreateTime string `json:"create_time"`//创建时间
	Status AccountStatus `json：“status”`//可用状态
	Version int64   //状态版本，乐观锁
}
//流水
type FlowWater struct {
	Flow int64  `json:"id" Table:"wt_acc_flows" Option:"auto"` //流水号
	Account string //账户
	Amount int64 //发生额
	Type FlowType //类型
	SubType string //子类型
	Begin int64  //期初
	End int64    //期末
	Detail string //摘要
	TimeStamp string //发生时间

}

//交易
type Transaction struct {
	Type TransactionType  `json:"type" Table:"wt_acc_transactions"`//交易类型
	Code string  //业务类型
	Request string //请求id，业务id
	Amount int64 //交易金额
	OutAccount string //出账账户
	InAccount string //入账账户
	OutFlow int64   //出账流水
	InFlow  int64   //入账流水
	Detail string //摘要
	Version int64  //版本
}

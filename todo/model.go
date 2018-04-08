package todo
/**
  待办
  1、获取待办
  2、待办完成
  3、生成待办

待办模型
 1、用户、标题、描述、类型、业务类型、业务id、待办时间、消息收到时间、最后更新时间、状态
 2、类型：通知、消息、任务
 */

 type TodoStyle byte

 const (
 	TS_NOTICE TodoStyle=1 //通知
 	TS_MESSAGE TodoStyle=3//消息
 	TS_TASK TodoStyle=5 //任务
 )

 type TodoItem struct {
 	Owner string //用户 who
 	Title string //标题
 	Content string //内容
 	Style   TodoStyle //类型
 	Code string //业务类型
 	Request string //业务id
 	TheTime string //待办业务时间
 	CreateTime string //创建时间
 	UpdateTime string //更新时间
 }
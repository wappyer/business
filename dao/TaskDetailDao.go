package dao

import (
	. "business/common"
	"business/dao/model"
)

const (
	// 任务明细状态
	TaskDetailStatusInit    = "init"    // 待审核
	TaskDetailStatusPublic  = "public"  // 已发布
	TaskDetailStatusRunning = "running" // 进行中
	TaskDetailStatusDone    = "done"    // 已完成
)

/**
 * 获取任务明细列表
 */
type ListTaskDetailArgs struct {
	TaskId []int
}

func ListTaskDetail(args *ListTaskDetailArgs) (int, []model.TaskDetail) {
	var detailList []model.TaskDetail
	session := DbEngine.Where("1=1")
	if len(args.TaskId) > 0 {
		session.And("task_id in" + WhereInInt(args.TaskId))
	}
	count, err := session.FindAndCount(&detailList)
	if err != nil {
		panic(NewDbErr(err))
	}
	return int(count), detailList
}

func InsertTaskDetail(detail *model.TaskDetail) *model.TaskDetail {
	detail.SetStatus(TaskDetailStatusInit).SetPublishTime(GetNow())
	if row := detail.Insert(); row == 0 {
		panic(NewRespErr(ErrTaskDetailInsert, ""))
	}
	if !detail.Info() {
		panic(NewRespErr(ErrTaskDetailInsert, ""))
	}
	return detail
}

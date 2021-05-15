package dao

import (
	. "business/common"
)

const (
	// 订单状态
	OrderStatusInit    = "init"    // 待审核
	OrderStatusPublic  = "publish" // 已发布
	OrderStatusRunning = "running" // 进行中
	OrderStatusSend    = "send"    // 已发货
	OrderStatusDone    = "done"    // 已完成

	// 订单评论状态
	OrderCommentStatusInit    = "init"
	OrderCommentStatusComment = "comment"
	OrderCommentStatusAgain   = "again"
	OrderCommentStatusCancel  = "cancel"
)

var OrderStatusMap = MapStr{
	OrderStatusInit:    "待审核",
	OrderStatusPublic:  "已发布",
	OrderStatusRunning: "进行中",
	OrderStatusSend:    "已发货",
	OrderStatusDone:    "已完成",
}
var OrderStatusSlice = []string{OrderStatusInit, OrderStatusPublic, OrderStatusRunning, OrderStatusSend, OrderStatusDone}

var OrderCommentStatusMap = MapStr{
	OrderCommentStatusInit:    "init",
	OrderCommentStatusComment: "comment",
	OrderCommentStatusAgain:   "again",
	OrderCommentStatusCancel:  "cancel",
}

/**
 * 获取订单列表
 */
type ListOrderArgs struct {
	Id              int    `json:"id"`
	TaskId          int    `json:"task_id"`
	ShopId          int    `json:"shop_id"`
	UserId          int    `json:"user_id"`
	Status          string `json:"status"`
	CreateTimeStart string `json:"create_time_start"`
	CreateTimeEnd   string `json:"create_time_end"`
	Limit           int    `json:"limit"`
	Offset          int    `json:"offset"`
}

type ListOrderRet struct {
	Id                int    `json:"id"`
	UserId            int    `json:"user_id"`
	TaskId            int    `json:"task_id"`
	TaskDetailId      int    `json:"task_detail_id"`
	ShopId            int    `json:"shop_id"`
	OnlineOrderId     int    `json:"online_order_id"`
	Status            string `json:"status"`
	StatusDesc        string `json:"status_desc"`
	CommentStatus     string `json:"comment_status"`
	CommentStatusDesc string `json:"comment_status_desc"`
	CreateTime        string `json:"create_time"`
	UpdateTime        string `json:"update_time"`
}

func ListOrder(args *ListOrderArgs) (int, []ListOrderRet) {
	session := DbEngine.Table("b_order").Alias("bo").
		Select("*").
		Where("1=1")

	if args.Id > 0 {
		session.And("bo.id = ?", args.Id)
	}
	if args.UserId > 0 {
		session.And("bo.user_id = ?", args.UserId)
	}
	if args.TaskId > 0 {
		session.And("bo.task_id = ?", args.TaskId)
	}
	if args.ShopId > 0 {
		session.And("bo.shop_id = ?", args.ShopId)
	}
	if args.Status != "" {
		session.And("bo.status = ?", args.Status)
	}
	if args.CreateTimeStart != "" {
		session.And("bo.create_time >= ?", args.CreateTimeStart)
	}
	if args.CreateTimeEnd != "" {
		session.And("bo.create_time <= ?", args.CreateTimeEnd)
	}

	session.OrderBy("create_time desc").Limit(args.Limit, args.Offset)

	var orderList []ListOrderRet
	count, err := session.FindAndCount(&orderList)
	if err != nil {
		panic(NewDbErr(err))
	}
	return int(count), orderList
}

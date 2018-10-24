package scheduler

import (
	"fmt"
	"sync"
)

// Status 代表调度器状态的类型。
type Status uint8

//代表调度器状态的类型
type Statuss uint8

const (
	// SCHED_STATUS_UNINITIALIZED 代表未初始化的状态。
	SCHED_STATUS_UNINITIALIZED Status = 0
	// SCHED_STATUS_INITIALIZING 代表正在初始化的状态。
	SCHED_STATUS_INITIALIZING Status = 1
	// SCHED_STATUS_INITIALIZED 代表已初始化的状态。
	SCHED_STATUS_INITIALIZED Status = 2
	// SCHED_STATUS_STARTING 代表正在启动的状态。
	SCHED_STATUS_STARTING Status = 3
	// SCHED_STATUS_STARTED 代表已启动的状态。
	SCHED_STATUS_STARTED Status = 4
	// SCHED_STATUS_STOPPING 代表正在停止的状态。
	SCHED_STATUS_STOPPING Status = 5
	// SCHED_STATUS_STOPPED 代表已停止的状态。
	SCHED_STATUS_STOPPED Status = 6
)

const  (
	//代表未初始化状态
	SCHED_STATUS_UNINITIALIZEDD Statuss = iota
	//代表正在初始化状态
	SCHED_STATUS_INITIALIZINGG
	//已初始化状态
	SCHED_STATUS_INITIALIZEDD
	//正在启动
	SCHED_STATUS_STARTINGG
	//已经启动
	SCHED_STATUS_STARTEDD
	//正在停止
	SCHED_STATUS_STOPPINGG
	//已经停止
	SCHED_STATUS_STOPPEDD
)

//check用于检查状态
//currentstatus表示当前状态
//wantedstatus 表示想要的状态
//检查规则
//1 处于正在初始化，正在启动，正在停止状态时，不能从外部改变状态
//2想要的状态只能是正在初始化，正在启动，正在停止状态中的一个
//3处于未初始化状态时，不能改变为正在启动或正在停止
//4处于已启动时，不能改变为正在初始化或正在启动
//5只要未处于已启动状态就不能改变为正在停止状态

func checkStatuss(currentStatus Statuss,wantedStatus Statuss,lock sync.Locker) (err error) {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	//正在初始化，正在启动，正在停止时不能从外部改变
	switch currentStatus {
	case SCHED_STATUS_INITIALIZINGG:
		err  = genError("the scheduler is being initialized")
	case SCHED_STATUS_STARTINGG:
		err = genError("the scheduler is being started")
	case SCHED_STATUS_STOPPINGG:
		err = genError("the scheduler is being stopped")
	}
	if err != nil {
		return
	}
	if currentStatus == SCHED_STATUS_UNINITIALIZEDD &&
		(wantedStatus == SCHED_STATUS_STARTINGG ||
			wantedStatus == SCHED_STATUS_STOPPINGG) {
				err = genError("the scheduler has not yet been initialized")
				return
	}
	switch wantedStatus {
	case SCHED_STATUS_INITIALIZINGG:
		switch currentStatus {
		case SCHED_STATUS_STARTEDD:
			err = genError("the scheduler has been started")
		}
	case SCHED_STATUS_STARTINGG:
		switch currentStatus {
		case SCHED_STATUS_INITIALIZEDD:
			err = genError("")
		case SCHED_STATUS_STARTEDD:
			err = genError("")
		}
	case SCHED_STATUS_STOPPINGG:
		if currentStatus != SCHED_STATUS_STARTEDD {
			err = genError("")
		}
	default:
		errMsg :=
			fmt.Sprintf("unsupported wanted status for check! (wantedStatus: %d)",
				wantedStatus)
		err = genError(errMsg)
	}
	return
}

// checkStatus 用于状态的检查。
// 参数currentStatus代表当前的状态。
// 参数wantedStatus代表想要的状态。
// 检查规则：
//     1. 处于正在初始化、正在启动或正在停止状态时，不能从外部改变状态。
//     2. 想要的状态只能是正在初始化、正在启动或正在停止状态中的一个。
//     3. 处于未初始化状态时，不能变为正在启动或正在停止状态。
//     4. 处于已启动状态时，不能变为正在初始化或正在启动状态。
//     5. 只要未处于已启动状态就不能变为正在停止状态。
func checkStatus(
	currentStatus Status,
	wantedStatus Status,
	lock sync.Locker) (err error) {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	//正在初始化，正在启动，正在停止时都不能从外部改变
	switch currentStatus {
	case SCHED_STATUS_INITIALIZING:
		err = genError("the scheduler is being initialized!")
	case SCHED_STATUS_STARTING:
		err = genError("the scheduler is being started!")
	case SCHED_STATUS_STOPPING:
		err = genError("the scheduler is being stopped!")
	}
	if err != nil {
		return
	}
	if currentStatus == SCHED_STATUS_UNINITIALIZED &&
		(wantedStatus == SCHED_STATUS_STARTING ||
			wantedStatus == SCHED_STATUS_STOPPING) {
		err = genError("the scheduler has not yet been initialized!")
		return
	}
	switch wantedStatus {
	case SCHED_STATUS_INITIALIZING:
		switch currentStatus {
		case SCHED_STATUS_STARTED:
			err = genError("the scheduler has been started!")
		}
	case SCHED_STATUS_STARTING:
		switch currentStatus {
		case SCHED_STATUS_UNINITIALIZED:
			err = genError("the scheduler has not been initialized!")
		case SCHED_STATUS_STARTED:
			err = genError("the scheduler has been started!")
		}
	case SCHED_STATUS_STOPPING:
		if currentStatus != SCHED_STATUS_STARTED {
			err = genError("the scheduler has not been started!")
		}
	default:
		errMsg :=
			fmt.Sprintf("unsupported wanted status for check! (wantedStatus: %d)",
				wantedStatus)
		err = genError(errMsg)
	}
	return
}

// GetStatusDescription 用于获取状态的文字描述。
func GetStatusDescription(status Status) string {
	switch status {
	case SCHED_STATUS_UNINITIALIZED:
		return "uninitialized"
	case SCHED_STATUS_INITIALIZING:
		return "initializing"
	case SCHED_STATUS_INITIALIZED:
		return "initialized"
	case SCHED_STATUS_STARTING:
		return "starting"
	case SCHED_STATUS_STARTED:
		return "started"
	case SCHED_STATUS_STOPPING:
		return "stopping"
	case SCHED_STATUS_STOPPED:
		return "stopped"
	default:
		return "unkown"
	}
}

//用于获取状态的文字描述
func GetStatusDescriptionn(statuss Statuss) string {
	switch statuss {
	case SCHED_STATUS_UNINITIALIZEDD:
		return "uninitialized"
	case SCHED_STATUS_INITIALIZINGG:
		return "initializing"
	case SCHED_STATUS_INITIALIZEDD:
		return "initialized"
	case SCHED_STATUS_STARTINGG:
		return "starting"
	case SCHED_STATUS_STARTEDD:
		return "started"
	case SCHED_STATUS_STOPPINGG:
		return "stopping"
	case SCHED_STATUS_STOPPEDD:
		return "stopped"
	default:
		return "unkown"
	}
}
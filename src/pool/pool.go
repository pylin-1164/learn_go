package pool

import (
	"fmt"
)

type PoolThead struct {

}

var CONFIG_THREAD_LIMIT = 3
var workTask = make(chan Runnable,CONFIG_THREAD_LIMIT)
var poolNum,indexNum = 0,0
var callback = make(chan struct{})
var closePool = make(chan bool)

func NewPoolThread(){
	go func() {
		CloseWork()
	}()
	for  {
		task := <- taskQueue
		NewWorkRuner()
		workTask <- task
	}
}

func CloseWork(){
	for {
		<- callback
		if poolNum > 1 && len(taskQueue)<poolNum{
			poolNum --
			closePool <- true
		}else{
			closePool <- false
		}
	}
}

func NewWorkRuner(){
	//TODO 协程数量 <= 最大阈值（或者）正在下载的文件总大小 <= 最大阈值
	if poolNum < CONFIG_THREAD_LIMIT {
		poolNum++
		go func() {
			indexNum++
			currentNum := indexNum
			for  {
				t := <-workTask
				/*switch task := t.(type) {
					case *Task :
						task = t.(*Task)
						fmt.Printf("runner[%d] task id : %s-%s \n",currentNum,task.userName,task.id)
				}*/
				fmt.Printf("runner[%d] task id \n",currentNum)
				t.Run()
				callback <- struct{}{}
				close := <- closePool
				if close{
					break
				}
			}
		}()
	}

}

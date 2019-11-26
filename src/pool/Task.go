package pool

import (
	"fmt"
	"math/rand"
	"time"
)

type Task struct {
	UserName 	string
	Id 			string
}


func (*Task) Run() {
	fmt.Println("logger...")
	time.Sleep(time.Millisecond*time.Duration(rand.Intn(8000)))
}

type TaskQueue struct {
	size 	int
}

var CONFIG_QUEUE_LIMIT = 500
var taskQueue chan Runnable

func NewTaskQueue() *TaskQueue{
	config := TaskQueue{size: CONFIG_QUEUE_LIMIT}
	taskQueue = make(chan Runnable,config.size)
	return &config
}

func (t *TaskQueue) Push(r Runnable){
	taskQueue <- r
}



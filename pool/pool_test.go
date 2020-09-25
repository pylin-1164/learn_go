package pool

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestTheadPool(t *testing.T){

	//启动线程池
	go func() {
		NewPoolThread()
	}()

	queue := NewTaskQueue()
	go func() {
		i := 0
		for i<20{
			task := &Task{UserName:"pyl",Id:fmt.Sprintf("%d",i)}
			i++
			if i >3{
				time.Sleep(time.Millisecond*time.Duration(rand.Intn(2000)))
			}
			//waitJobs.push(job)
			queue.Push(task)
		}

	}()
	go func() {
		i := 0
		for i<20{
			task := &Task{UserName:"zhangsan",Id:fmt.Sprintf("%d",i)}
			i++
			if i>3{
				time.Sleep(time.Millisecond*time.Duration(rand.Intn(2000)))
			}
			//waitJobs.push(job)
			queue.Push(task)
		}

	}()
	time.Sleep(time.Minute*5)
}
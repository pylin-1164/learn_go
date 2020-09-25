# pool
* 1.  动态线程池，无任务时只保留一个线程。
* 2.  多条件限制，限制在一定条件下，即使线程数量未到阈值，仍然不再开放新线程。

``` go
//启动线程池
go func() {
  NewPoolThread()
}()

//继承Runnable
type Task struct {
}
func (*Task) Run() {
	panic("imp ...")
}

// 添加任务
queue.Push(task)

//修改该方法，自行扩展是否开启新线程的其他条件
func NewWorkRuner(){
	//TODO 协程数量 <= 最大阈值（或者）正在下载的文件总大小 <= 最大阈值
	if poolNum < CONFIG_THREAD_LIMIT {
      //...
  }
```

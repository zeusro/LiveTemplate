package goroutine

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

/*
阅读程序功能说明，补全缺失的代码。时间30分钟以内。
不要修改代码结构，仅需要补全注释标有todo的地方
功能说明
本程序的作用是，根据输入参数并发导出用户信息表中的数据。详情如下：
1、用户表结构为：uid，name。表名 user
2、uid是用户主键，正整数且是连续递增的。
3、该表有1亿条数据左右，注意索引的使用。

输入参数：

	-limit int
	      导出数据时数据库的分页大小 (default 3)
	-max int
	      导出数据的最大uid (default 10)
	-min int
	      导出数据的最小uid (default 1)
	-worker int
	      控制生产者并发数 (default 2

例如：
2个工作协程，sql分页大小为3，导出uid从1到10的数据
输入
go run ./main.go -worker=2 -limit=3 -min=1 -max=10
*/
func Main(worker, limit, uidMin, uidMax int) {
	// 处理输入参数（由调用方传入，而不是使用 flag）
	if worker <= 0 || limit <= 0 {
		log.Fatal("worker 和 limit 必须大于 0")
	}
	param := Input{
		limit:     limit,
		uidMin:    uidMin,
		uidMax:    uidMax,
		workerNum: worker,
	}

	var (
		queue = make(chan []*UserInfo, param.workerNum*2)
		// todo
		producerWg sync.WaitGroup
		consumerWg sync.WaitGroup
	)

	// todo 1个消费者（消费数据）
	consumer(queue, &consumerWg)

	// todo 多个生产者（导出数据）
	producer(param, queue, &producerWg)

	// todo 等待生产者退出
	producerWg.Wait()

	// todo 通知消费者退出（队列不会再写入新数据）
	close(queue)

	// todo 等待消费者退出
	consumerWg.Wait()

	log.Println("all jobs done")
}

// producer 获取数据（生产者）
// todo 如果有必要可以更改参数
func producer(param Input, queue chan []*UserInfo, wg *sync.WaitGroup) {
	// 拆分导出任务，按照输入的并发任务个数进行uid拆分
	jobs := splitJobs(param.uidMin, param.uidMax, param.workerNum)

	// 并发导出
	wg.Add(len(jobs))
	for i, j := range jobs {
		log.Printf("producer->worker:%d uid_min:%d uid_max:%d\n", i, j.UidMin, j.UidMax)
		go func(jj Job) {
			defer wg.Done()
			worker(jj, param.limit, queue)
		}(j)
	}
}

// splitJobs 拆分工作
// 因为uid是连续递增的，所以可以根据uid总数进行任务拆分
// 输入：uid起始范围，并发数
// 返回：每个任务的uid范围
func splitJobs(uidMin, uidMax int, workerNum int) []Job {
	if uidMax < uidMin {
		return []Job{}
	}
	// uid总数
	uidTotal := uidMax - uidMin + 1

	var jobSize int
	if uidTotal < workerNum { // 如果总的需要查询的数据数还没有worker多
		workerNum = uidTotal
		jobSize = 1
	} else {
		// 数据总数 除以 worker数 得到每个worker需要查询的数据大小
		jobSize = int(math.Ceil(float64(uidTotal) / float64(workerNum)))
	}

	// 拆分工作
	var jobs []Job
	for i := 0; i < workerNum; i++ {
		t := Job{
			UidMin: uidMin,
			UidMax: uidMin + jobSize - 1,
		}
		if t.UidMax > uidMax {
			t.UidMax = uidMax
		}
		jobs = append(jobs, t)
		uidMin = t.UidMax + 1
	}
	return jobs
}

// consumer 消费数据
// todo 如果有必要可以更改参数
func consumer(queue chan []*UserInfo, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for us := range queue {
			for _, u := range us {
				// 模拟消费
				log.Println(u.Name)
			}
		}
	}()
}

// worker 读取用户表数据协程
func worker(j Job, limit int, buf chan []*UserInfo) {
	var uidBegin = j.UidMin
	for {
		if uidBegin+limit > j.UidMax {
			limit = j.UidMax - uidBegin + 1
		}
		//todo 补齐SQL，注意需要使用索引来处理超大数据分页方式
		sql := fmt.Sprintf("SELECT `uid`,`name` FROM `user` WHERE `uid` >= %d AND `uid` <= %d ORDER BY `uid` LIMIT %d", uidBegin, j.UidMax, limit)
		// todo 查询数据
		users := queryDB(sql)
		if len(users) == 0 {
			break
		}
		buf <- users
		// todo 翻页
		uidBegin = users[len(users)-1].Uid + 1
		if uidBegin > j.UidMax {
			break
		}
	}
}

/***************************************** 分割线以下没有需要补齐的代码 *****************************************/

// UserInfo 用户表
type UserInfo struct {
	Uid  int
	Name string
}

// Job 导出任务配置
type Job struct {
	UidMin int
	UidMax int
}

// Input 用户输入参数
type Input struct {
	limit     int
	uidMin    int
	uidMax    int
	workerNum int
}

// queryDB 模拟查询数据库，此处可忽略
func queryDB(sql string) []*UserInfo {
	time.Sleep(time.Millisecond * 50)
	return []*UserInfo{{Uid: 1, Name: sql}}
}

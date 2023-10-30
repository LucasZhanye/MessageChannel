package pool

import (
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestOneWorker(t *testing.T) {
	convey.Convey("只有一个worker", t, func() {
		p := NewWorkerPool(1)

		go p.Run()

		target := []string{"apple", "watermelon", "peach", "banana"}
		result := []string{}

		for _, t := range target {
			str := t
			task := &Task{
				Fn: func(...any) {
					result = append(result, str)
				},
				Param: []any{},
			}

			p.Put(task)
		}

		p.Stop()
		fmt.Println(result)

		convey.So(result, convey.ShouldResemble, target)
	})
}

func TestSimple(t *testing.T) {
	convey.Convey("正常例子", t, func() {
		p := NewWorkerPool(2)

		go p.Run()

		target := []string{"apple", "watermelon", "peach", "banana"}
		var result = map[string]struct{}{}

		for _, t := range target {
			str := t
			task := &Task{
				Fn: func(...any) {
					result[str] = struct{}{}
				},
			}

			p.Put(task)
		}

		p.Stop()

		if len(target) != len(result) {
			t.Fatalf("Pool task error, element length not equal, target =  %d, result = %d", len(target), len(result))
		}

		for _, s := range target {
			if _, ok := result[s]; !ok {
				t.Fatalf("Pool task error, can not find element: %s", s)
			}
		}
	})
}

func TestWorker(t *testing.T) {
	convey.Convey("测试Worker重复使用", t, func() {
		p := NewWorkerPool(10)

		go p.Run()

		chanArr := [5]chan struct{}{}
		for i := 0; i < 5; i++ {
			index := i
			chanArr[index] = make(chan struct{})
			task := &Task{
				Fn: func(params ...any) {
					num := params[0].(int)
					<-chanArr[num]
				},
				Param: []any{
					index,
				},
			}

			p.Put(task)
		}

		workers := p.GetCurrentWorkerSize()
		convey.So(workers, convey.ShouldEqual, 5)

		for i := 0; i < 3; i++ {
			index := i
			chanArr[index] <- struct{}{}
		}

		time.Sleep(1 * time.Second)

		for i := 0; i < 3; i++ {
			index := i
			task := &Task{
				Fn: func(params ...any) {
					fmt.Println("index = ", params[0])
				},
				Param: []any{
					index,
				},
			}

			p.Put(task)
		}

		workers = p.GetCurrentWorkerSize()
		convey.So(workers, convey.ShouldEqual, 5)
	})
}

func TestWorkerBusy(t *testing.T) {
	convey.Convey("测试Worker忙", t, func() {
		p := NewWorkerPool(10)

		go p.Run()

		waitSize := p.GetWaitingTaskSize()

		convey.So(waitSize, convey.ShouldEqual, 0)

		chanArr := [10]chan struct{}{}
		for i := 0; i < 10; i++ {
			index := i
			chanArr[index] = make(chan struct{})
			task := &Task{
				Fn: func(params ...any) {
					num := params[0].(int)
					<-chanArr[num]
				},
				Param: []any{
					index,
				},
			}

			p.Put(task)
		}

		for i := 0; i < 5; i++ {
			index := i
			task := &Task{
				Fn: func(params ...any) {
					num := params[0].(int)
					<-chanArr[num]
				},
				Param: []any{
					index,
				},
			}

			p.Put(task)
		}
		time.Sleep(500 * time.Millisecond)

		workers := p.GetCurrentWorkerSize()
		convey.So(workers, convey.ShouldEqual, 10)

		waitSize = p.GetWaitingTaskSize()
		convey.So(waitSize, convey.ShouldEqual, 5)

		for i := 0; i < 4; i++ {
			index := i
			chanArr[index] <- struct{}{}
		}

		time.Sleep(taskIdleTimeout)

		waitSize = p.GetWaitingTaskSize()
		convey.So(waitSize, convey.ShouldEqual, 1)
	})
}

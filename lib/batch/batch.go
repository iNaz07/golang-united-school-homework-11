package batch

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {

	sem := semaphore.NewWeighted(pool)
	ctx := context.Background()
	mx := &sync.Mutex{}
	for i := 0; i < int(n); i++ {
		if err := sem.Acquire(ctx, 1); err != nil {
			fmt.Println(fmt.Errorf("wait for resources: %w", err))
		}
		go func(i int64) {
			u := getOne(i)
			mx.Lock()
			res = append(res, u)
			mx.Unlock()
			sem.Release(1)
		}(int64(i))
	}
	if err := sem.Acquire(ctx, pool); err != nil {
		fmt.Println(fmt.Errorf("wait for resources: %w", err))
	}
	//
	//var wg sync.WaitGroup
	//sema := make(chan struct{}, pool)
	//
	//for i := 0; i < int(n); i++ {
	//	wg.Add(1)
	//	sema <- struct{}{}
	//	go func(i int64) {
	//		res = append(res, getOne(i))
	//		<-sema
	//		wg.Done()
	//	}(int64(i))
	//}
	//wg.Wait()
	return
}

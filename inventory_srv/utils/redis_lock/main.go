package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

func main() {
	// Create a Redis client and pool
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "127.0.0.1:6379",
		DB:   1,
	})
	pool := goredis.NewPool(client)

	// Create an instance of redsync
	rs := redsync.New(pool)

	// Number of goroutines
	gNum := 11
	mutexName := "421"

	var wg sync.WaitGroup
	wg.Add(gNum)
	for i := 0; i < gNum; i++ {
		go func(id int) {
			defer wg.Done()
			mutex := rs.NewMutex(mutexName, redsync.WithExpiry(8*time.Second)) // 设置锁的超时时间为8秒

			for {
				fmt.Printf("Goroutine %d: 尝试获取锁\n", id)
				if err := mutex.Lock(); err != nil {
					fmt.Printf("Goroutine %d: 获取锁失败，重试中...\n", id)
					time.Sleep(500 * time.Millisecond) // 等待500毫秒后重试
					continue
				}

				fmt.Printf("Goroutine %d: 获取锁成功\n", id)
				time.Sleep(2 * time.Second) // 模拟操作
				fmt.Printf("Goroutine %d: 开始释放锁\n", id)
				if ok, err := mutex.Unlock(); !ok || err != nil {
					fmt.Printf("Goroutine %d: 释放锁失败\n", id)
					panic("unlock failed")
				}
				fmt.Printf("Goroutine %d: 释放锁成功\n", id)
				break
			}
		}(i)
	}
	wg.Wait()
}

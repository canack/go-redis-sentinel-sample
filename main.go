package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName: "mymaster",
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			addrAndPort := strings.Split(addr, ":")
			if len(addrAndPort) != 2 {
				return nil, fmt.Errorf("invalid address: %s", addr)
			}
			addr = ":" + addrAndPort[1]
			return net.Dial("tcp", addr)
		},
		SentinelAddrs:   []string{"redis-master:26380"},
		DialTimeout:     10 * time.Second,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		MaxRetries:      3,
		MinRetryBackoff: 1 * time.Second,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Redis connection error:", err)
		return
	}
	fmt.Println("Redis connection established:", pong)

	for i := 0; ; i++ {
		key := fmt.Sprintf("key%d", i)
		val := fmt.Sprintf("value%d", i)
		err = rdb.Set(ctx, key, val, 0).Err()
		if err != nil {
			fmt.Println("Set error:", err)
			return
		}
		fmt.Println("Value set successfully")

		val, err := rdb.Get(ctx, key).Result()
		if err != nil {
			fmt.Println("Get error:", err)
		} else {
			fmt.Println("Value retrieved successfully:", val)
		}

		fmt.Println("Sleeping for 1 second")
		time.Sleep(1 * time.Second)
	}
}

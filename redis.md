https://redis.io/docs/getting-started/

## install

https://redis.io/docs/getting-started/installation/install-redis-on-linux/

```
$redis-server -v
Redis server v=7.0.3 sha=00000000:0 malloc=jemalloc-5.2.1 bits=64 build=9e1e27bf47eab674
```

redis server 起動

````
$redis-server
略
3104:M 18 Jul 2022 07:09:59.156 * monotonic clock: POSIX clock_gettime
                _._
           _.-``__ ''-._
      _.-``    `.  `_.  ''-._           Redis 7.0.3 (00000000/0) 64 bit
  .-`` .-```.  ```\/    _.,_ ''-._
 (    '      ,       .-`  | `,    )     Running in standalone mode
 |`-._`-...-` __...-.``-._|'` _.-'|     Port: 6379
 |    `-._   `._    /     _.-'    |     PID: 3104
  `-._    `-._  `-./  _.-'    _.-'
 |`-._`-._    `-.__.-'    _.-'_.-'|
 |    `-._`-._        _.-'_.-'    |           https://redis.io
  `-._    `-._`-.__.-'_.-'    _.-'
 |`-._`-._    `-.__.-'    _.-'_.-'|
 |    `-._`-._        _.-'_.-'    |
  `-._    `-._`-.__.-'_.-'    _.-'
      `-._    `-.__.-'    _.-'
          `-._        _.-'
              `-.__.-'


````

```
$redis-cli -v
redis-cli 7.0.3
```

go redis
https://redis.uptrace.dev/
https://pkg.go.dev/github.com/go-redis/redis/v8#section-readme

go get

```
go get github.com/go-redis/redis/v8
```

https://redis.uptrace.dev/guide/go-redis.html#installation

```go
package main

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:	  "localhost:6379",
		Password: "", // no password set
		DB:		  0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
```

動作確認

```
$go run sample/redis.go
key value
key2 does not exist
```

https://redis.io/commands/getset/

```
$redis-cli
127.0.0.1:6379> get key
"value"
127.0.0.1:6379> get key2
(nil)
127.0.0.1:6379> get key
"value"
127.0.0.1:6379> set key2 value2
OK
127.0.0.1:6379> get key
"value"
127.0.0.1:6379> get key2
"value2"
```

set ttl

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 3*time.Second).Err()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 4; i++ {
		val, err := rdb.Get(ctx, "key").Result()
		if err == redis.Nil {
			fmt.Println("key does not exist")
		} else if err != nil {
			panic(err)
		}
		fmt.Printf("i: %d key: %s\n", i, val)
		time.Sleep(1 * time.Second)
	}
}
```

```
$go run sample/redis.go
i: 0 key: value
i: 1 key: value
i: 2 key: value
key does not exist
i: 3 key:
```

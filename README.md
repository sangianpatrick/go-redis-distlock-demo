# go-redis-distlock-demo
Simple Distributed Lock with Redis (Demo)

# How To Run
- `$ PORT=<your port> go run main.go`
- Make sure to hit the `/high` endpoint first and subsequently hit the `/low`. The `/high` will be the first request that acquire the lock and the `/low` will wait until the previous request is done or release the lock. If the second request still acquiring the lock and reach the maximum duration, it will returns ErrLocked
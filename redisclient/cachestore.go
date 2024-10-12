package redisclient

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"

    "github.com/go-redis/redis/v8"
)

// Config is the redis client config.
type Config struct {
    Addr     string
    Password string
    Database int
}

// Cache is a redis cache struct.
type Cache struct {
    cache *redis.Client
    mu    sync.RWMutex
}

// NewCache creates a new cache instance.
func NewCache(ctx context.Context, config Config, checkPing bool) (*Cache, error) {
    fmt.Println()
    result := &Cache{
        cache: redis.NewClient(&redis.Options{
            Addr:     config.Addr,
            Password: config.Password,
            DB:       config.Database,
        }),
    }
    if !checkPing {
        return result, nil
    }
    // If checkPing is true, do ping request on redis.
    _, err := result.cache.Ping(ctx).Result()
    if err != nil {
        return nil, err
    }
    fmt.Println("Connected to Redis")
    return result, nil
}

// GetUniqueIDsForCurrentMinute fetches the number unique IDs that processed in the current minute.
func (c *Cache) GetUniqueIDsForCurrentMinute(ctx context.Context) (int, error) {
    now := time.Now().UTC()
    return c.getUniqueIDsForMinute(ctx, now.Minute())
}

// GetUniqueIDsForPreviousMinute fetches the number unique IDs that processed in the previous minute.
func (c *Cache) GetUniqueIDsForPreviousMinute(ctx context.Context) (int, error) {
    prev := time.Now().UTC().Add(-1 * time.Minute)
    return c.getUniqueIDsForMinute(ctx, prev.Minute())
}

// getUniqueIDsForMinute fetches the number of unique IDs that were processed for a given minute.
func (c *Cache) getUniqueIDsForMinute(ctx context.Context, d int) (int, error) {
    var (
        n      int
        cursor uint64
        keys   []string
        err    error
    )
    c.mu.RLock()
    defer c.mu.RUnlock()
    key := fmt.Sprintf("%d*", d)
    for {
        keys, cursor, err = c.cache.Scan(ctx, cursor, key, 10001).Result()
        if err != nil {
            log.Fatal(err)
            return n, err
        }
        n += len(keys)
        if cursor == 0 {
            break
        }
    }
    return n, nil
}

// SetIDInCache stores the ID in the current minute. The key is minute_id. Returns (true, nil),
// in case the ID was not present for the current minute or (false, nil) otherwise.
func (c *Cache) SetIDInCache(ctx context.Context, id int64) (bool, error) {
    c.mu.Lock()
    defer c.mu.Unlock()
    now := time.Now().UTC()
    key := fmt.Sprintf("%d_%d", now.Minute(), id)

    added, err := c.cache.SetNX(ctx, key, true, 10*time.Minute).Result()
    if err != nil {
        log.Printf("Error interacting with Redis: %v", err)
        return false, err
    }
    return added, nil
}

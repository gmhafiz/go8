package cache

import (
	"context"
	"github.com/gmhafiz/go8/ent/gen"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack"

	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/domain/author/repository/database"
)

type Cache struct {
	service database.Repository
	cache   *redis.Client
}

type AuthorRedisService interface {
	List(ctx context.Context, f *author.Filter) ([]*gen.Author, int64, error)
	Update(ctx context.Context, toAuthor *author.Update) (*gen.Author, error)
	Delete(ctx context.Context, id int64) error
}

func NewRedisCache(service database.Repository, cache *redis.Client) *Cache {
	return &Cache{
		service: service,
		cache:   cache,
	}
}

func (c *Cache) List(ctx context.Context, f *author.Filter) ([]*gen.Author, int64, error) {
	// We want to store both list and the count together in one cache key.
	type result struct {
		List []*gen.Author `json:"list"`
		Num  int64         `json:"num"`
	}

	url := ctx.Value(author.CacheURL).(string)
	res := &result{}

	val, err := c.cache.Get(ctx, url).Result()
	if err == redis.Nil || err != nil {
		list, num, err := c.service.List(ctx, f)
		if err != nil {
			return nil, 0, err
		}
		res.List = list
		res.Num = num
		cacheEntry, err := msgpack.Marshal(res)
		//cacheEntry, err := json.Marshal(res)
		if err != nil {
			return c.service.List(ctx, f)
		}

		err = c.cache.Set(ctx, url, cacheEntry, 8*time.Second).Err()
		if err != nil {
			return c.service.List(ctx, f)
		}

		return list, num, nil
	}

	err = msgpack.Unmarshal([]byte(val), &res)
	//err = json.Unmarshal([]byte(val), &res)
	if err != nil {
		return c.service.List(ctx, f)
	}

	return res.List, res.Num, nil
}

func (c *Cache) Update(ctx context.Context, toAuthor *author.Update) (*gen.Author, error) {
	c.invalidate(ctx)

	return c.service.Update(ctx, toAuthor)
}

func (c *Cache) Delete(ctx context.Context, id int64) error {
	c.invalidate(ctx)

	return c.service.Delete(ctx, id)
}

func (c *Cache) invalidate(ctx context.Context) {
	url := ctx.Value(author.CacheURL).(string)
	split := strings.Split(url, "/")
	baseURL := strings.Join(split[:4], "/")

	keys, _ := c.cache.Keys(ctx, baseURL+"*").Result()
	for _, key := range keys {
		_ = c.cache.Del(ctx, key)
	}
}

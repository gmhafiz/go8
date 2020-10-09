package jobs

import (
	"database/sql"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog"
	"go8ddd/configs"
)

type Jobs struct {
	WorkerPool *work.WorkerPool
	Logger     *zerolog.Logger
	Enqueuer   *work.Enqueuer
}

func New(cfg *configs.Configs, log zerolog.Logger, db *sql.DB) *Jobs {
	redisClient := newPool(cfg)
	enqueuer := newQueuer(redisClient.Pool)
	pool := workerPool(redisClient.Pool, log, db)
	pool.Start()

	return &Jobs{
		WorkerPool: pool,
		Enqueuer:   enqueuer,
		Logger:     &log,
	}
}

func newQueuer(pool *redis.Pool) *work.Enqueuer {
	return work.NewEnqueuer("go8_jobs", pool)
}

func workerPool(pool *redis.Pool, log zerolog.Logger, db *sql.DB) *work.WorkerPool {
	workerPool := work.NewWorkerPool(Context{}, 10, "go8_jobs", pool)

	// Add middleware that will be executed for each job
	workerPool.Middleware(func(c *Context, job *work.Job, next work.NextMiddlewareFunc) error {
		c.Log = log
		c.DB = db
		log.Info().Msg("through middleware")
		return next()
	})

	// Map the name of jobs to handler functions
	workerPool.Job("send_email", (*Context).SendEmail)

	return workerPool
}

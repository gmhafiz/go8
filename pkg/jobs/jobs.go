package jobs

import (
	"log"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog"
)

type Jobs struct {
	WorkerPool *work.WorkerPool
	Logger     *zerolog.Logger
}

type Context struct {
	ID   int64
	Jobs *Jobs
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	return next()
}

func (c *Context) SendEmail(job *work.Job) error {
	time.Sleep(time.Second * 5)
	// Extract arguments:
	addr := job.ArgString("address")
	subject := job.ArgString("subject")
	if err := job.ArgError(); err != nil {
		return err
	}

	// Go ahead and send the email...
	// sendEmailTo(addr, subject)
	log.Println(addr)
	log.Println(subject)

	return nil
}

func (j *Jobs) NewQueuer(redisPool *redis.Pool) *work.Enqueuer {
	return work.NewEnqueuer("go8_jobs", redisPool)
}

func New(redisPool *redis.Pool, logger *zerolog.Logger) *Jobs {
	workerPool := work.NewWorkerPool(Context{}, 10, "go8_jobs", redisPool)

	// Add middleware that will be executed for each job
	workerPool.Middleware((*Context).Log)

	// Map the name of jobs to handler functions
	workerPool.Job("send_email", (*Context).SendEmail)

	return &Jobs{
		WorkerPool: workerPool,
		Logger:     logger,
	}
}

package jobs

import (
	"context"
	"database/sql"
	"time"

	"github.com/gocraft/work"
	"github.com/rs/zerolog"

	"go8ddd/internal/model"
)

type Context struct {
	ID  int64
	Log zerolog.Logger
	DB  *sql.DB
}

func (c *Context) SendEmail(job *work.Job) error {
	time.Sleep(time.Second * 5)
	// Extract arguments:
	addr := job.ArgString("address")
	subject := job.ArgString("subject")
	if err := job.ArgError(); err != nil {
		return err
	}

	books, err := model.Books().All(context.Background(), c.DB)
	if err != nil {

	}
	c.Log.Info().Msgf("%v", books[0])

	// Go ahead and send the email...
	// sendEmailTo(addr, subject)
	c.Log.Info().Msg(addr)
	c.Log.Info().Msg(subject)

	return nil
}

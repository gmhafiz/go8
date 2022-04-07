package repository

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/gmhafiz/go8/ent/gen/runtime"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"

	"github.com/gmhafiz/go8/ent/gen"
	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/utility/filter"
	parseTime "github.com/gmhafiz/go8/internal/utility/time"
)

var (
	dockerDB *db
)

type db struct {
	Conn *sql.DB
	Ent  *gen.Client
}

func (d db) migrate(upOrDown string) {
	driver, err := postgres.WithInstance(d.Conn, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../../../database/migrations/",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	switch upOrDown {
	case "down":
		if err = m.Down(); err != nil {
			log.Fatalln(err)
		}
	default:
		if err = m.Up(); err != nil {
			log.Fatalln(err)
		}
	}
}

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)
	log.Println(databaseUrl)

	_ = resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	dockerDB = &db{}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		dockerDB.Conn, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return dockerDB.Conn.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Performing a migration this way means all tests in this package shares
	// the same db schema across all unit test.
	// If isolation is needed, then do away with using `testing.M`. Do a
	// migration for each test handler instead.
	dockerDB.migrate("up")

	// We can access database with m.hostAndPort or m.databaseUrl
	// port changes everytime a new docker instance is run
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestAuthorRepository_Create(t *testing.T) {
	type args struct {
		author *author.CreateRequest
	}

	type want struct {
		author *gen.Author
		err    error
	}

	type test struct {
		name string
		args
		want
	}

	startTime := time.Now()

	var authorBooks []*gen.Book
	authorBooks = append(authorBooks, &gen.Book{
		ID:            1,
		Title:         "Title",
		PublishedDate: parseTime.Parse("2022-02-12T15:04:05Z"),
		Description:   "Description",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeletedAt:     &time.Time{},
		Edges:         gen.BookEdges{},
	})

	tests := []test{
		{
			name: "nil",
			args: args{nil},
			want: want{
				author: nil,
				err:    fmt.Errorf("request cannot be nil"),
			},
		},
		{
			name: "empty names",
			args: args{
				author: &author.CreateRequest{
					FirstName:  "",
					MiddleName: "",
					LastName:   "",
					Books:      nil,
				},
			},
			want: want{
				author: &gen.Author{
					ID:         1,
					FirstName:  "",
					MiddleName: "",
					LastName:   "",
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
					DeletedAt:  &time.Time{},
				},
				err: nil,
			},
		},
		{
			name: "normal",
			args: args{
				author: &author.CreateRequest{
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
					Books:      nil,
				},
			},
			want: want{
				author: &gen.Author{
					ID:         2,
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
					DeletedAt:  nil,
				},
				err: nil,
			},
		},
		{
			name: "add one book",
			args: args{
				author: &author.CreateRequest{
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
					Books: []author.Book{
						{
							Title:         "Title",
							PublishedDate: "2022-02-12T15:04:05Z",
							Description:   "Description",
						},
					},
				},
			},
			want: want{
				author: &gen.Author{
					ID:         3,
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
					DeletedAt:  nil,
					Edges: gen.AuthorEdges{
						Books: authorBooks,
					},
				},
				err: nil,
			},
		},
	}

	client := dbClient()
	repo := New(client)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			created, err := repo.Create(ctx, test.args.author)
			assert.Equal(t, err, test.want.err)

			if created == nil {
				return
			}
			assert.Equal(t, test.want.author.ID, created.ID)
			assert.Equal(t, test.want.author.FirstName, created.FirstName)
			assert.Equal(t, test.want.author.MiddleName, created.MiddleName)
			assert.Equal(t, test.want.author.LastName, created.LastName)
			assert.True(t, startTime.Before(created.CreatedAt) || startTime.Equal(created.CreatedAt))
			assert.True(t, startTime.Before(created.UpdatedAt) || startTime.Equal(created.UpdatedAt))
			assert.Nil(t, created.DeletedAt)

			if created.Edges.Books == nil {
				return
			}
			for i := range created.Edges.Books {
				assert.Equal(t, test.want.author.Edges.Books[i].Title, created.Edges.Books[i].Title)
				assert.Equal(t, test.want.author.Edges.Books[i].Description, created.Edges.Books[i].Description)
				assert.Equal(t, test.want.author.Edges.Books[i].PublishedDate, created.Edges.Books[i].PublishedDate)
				assert.True(t, startTime.Before(test.want.author.Edges.Books[i].CreatedAt) || startTime.Equal(test.want.author.Edges.Books[i].CreatedAt))
				assert.True(t, startTime.Before(test.want.author.Edges.Books[i].UpdatedAt) || startTime.Equal(test.want.author.Edges.Books[i].UpdatedAt))
				assert.Nil(t, created.DeletedAt)
			}
		})
	}
}

func TestRepository_Read(t *testing.T) {
	startTime := time.Now()

	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		args    args
		want    *gen.Author
		wantErr error
	}{
		{
			name: "Read one author without book(s)",
			args: args{
				ctx: context.Background(),
				id:  4,
			},
			want: &gen.Author{
				ID:         4,
				FirstName:  "First",
				MiddleName: "Middle",
				LastName:   "Last",
				CreatedAt:  time.Time{},
				UpdatedAt:  time.Time{},
				DeletedAt:  nil,
				Edges: gen.AuthorEdges{
					Books: nil,
				},
			},
			wantErr: nil,
		},
	}

	oneAuthorWithoutBooks := &author.CreateRequest{
		FirstName:  "First",
		MiddleName: "Middle",
		LastName:   "Last",
		Books:      nil,
	}

	client := dbClient()
	repo := New(client)
	ctx := context.Background()

	created, err := repo.Create(ctx, oneAuthorWithoutBooks)
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Read(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, err)

			assert.Equal(t, created.ID, got.ID)
			assert.Equal(t, created.FirstName, got.FirstName)
			assert.Equal(t, created.MiddleName, got.MiddleName)
			assert.Equal(t, created.LastName, got.LastName)
			assert.True(t, startTime.Before(created.CreatedAt) || startTime.Equal(created.CreatedAt))
			assert.True(t, startTime.Before(created.UpdatedAt) || startTime.Equal(created.UpdatedAt))
			assert.Nil(t, created.DeletedAt)
		})
	}
}

func TestRepository_List(t *testing.T) {
	startTime := time.Now().Add(time.Duration(-5) * time.Second)

	type fields struct {
		ent *gen.Client
	}
	type args struct {
		ctx context.Context
		f   *author.Filter
	}

	// From previous unit tests, 4 records have been inserted.
	want4 := []*gen.Author{
		{
			ID:         1,
			FirstName:  "",
			MiddleName: "",
			LastName:   "",
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
			DeletedAt:  nil,
			Edges:      gen.AuthorEdges{},
		},
		{
			ID:         2,
			FirstName:  "First",
			MiddleName: "Middle",
			LastName:   "Last",
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
			DeletedAt:  nil,
			Edges:      gen.AuthorEdges{},
		},
		{
			ID:         3,
			FirstName:  "First",
			MiddleName: "Middle",
			LastName:   "Last",
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
			DeletedAt:  nil,
			Edges: gen.AuthorEdges{
				Books: []*gen.Book{
					{
						ID:            1,
						Title:         "Title",
						PublishedDate: parseTime.Parse("2022-02-12T15:04:05Z"),
						Description:   "Description",
						CreatedAt:     time.Time{},
						UpdatedAt:     time.Time{},
						DeletedAt:     nil,
					},
				},
			},
		},
		{
			ID:         4,
			FirstName:  "First",
			MiddleName: "Middle",
			LastName:   "Last",
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
			DeletedAt:  nil,
			Edges:      gen.AuthorEdges{},
		},
	}

	client := dbClient()
	repo := New(client)
	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*gen.Author
		size    int
		wantErr error
	}{
		{
			name: "show all 4 currently been added",
			fields: fields{
				ent: client,
			},
			args: args{
				ctx: ctx,
				f:   author.Filters(nil),
			},
			want:    want4,
			size:    4,
			wantErr: nil,
		},
		{
			name: "filter with non-existent name",
			fields: fields{
				ent: client,
			},
			args: args{
				ctx: ctx,
				f: &author.Filter{
					Base: filter.Filter{
						Page:          1,
						Offset:        10,
						Limit:         10,
						DisablePaging: false,
						Sort:          nil,
						Search:        false,
					},
					FirstName:  "nonexistent",
					MiddleName: "nonexistent",
					LastName:   "nonexistent",
				},
			},
			want:    make([]*gen.Author, 0),
			size:    4,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, size, err := repo.List(tt.args.ctx, tt.args.f)
			assert.Equal(t, tt.wantErr, err)

			for i := range got {
				assert.Equal(t, tt.want[i].ID, got[i].ID)
				assert.Equal(t, tt.want[i].FirstName, got[i].FirstName)
				assert.Equal(t, tt.want[i].MiddleName, got[i].MiddleName)
				assert.Equal(t, tt.want[i].LastName, got[i].LastName)
				assert.True(t, startTime.Before(got[i].CreatedAt) || startTime.Equal(got[i].CreatedAt))
				assert.True(t, startTime.Before(got[i].UpdatedAt) || startTime.Equal(got[i].UpdatedAt))
				assert.Nil(t, got[i].DeletedAt)

				for j := range tt.want[i].Edges.Books {
					assert.Equal(t, tt.want[i].Edges.Books[j].ID, got[i].Edges.Books[j].ID)
					assert.Equal(t, tt.want[i].Edges.Books[j].Title, got[i].Edges.Books[j].Title)
					assert.Equal(t, tt.want[i].Edges.Books[j].Description, got[i].Edges.Books[j].Description)
					assert.True(t, startTime.Before(got[j].CreatedAt) || startTime.Equal(got[j].CreatedAt))
					assert.True(t, startTime.Before(got[j].UpdatedAt) || startTime.Equal(got[j].UpdatedAt))
					assert.Nil(t, got[j].DeletedAt)
				}
			}
			assert.Equal(t, tt.size, size)
		})
	}
}

func TestRepository_Update(t *testing.T) {
	startTime := time.Now().Add(time.Duration(-5) * time.Second)

	oneAuthor := &author.CreateRequest{
		FirstName:  "First",
		MiddleName: "Middle",
		LastName:   "Last",
		Books: []author.Book{
			{
				Title:         "Title",
				PublishedDate: "2022-02-12T15:04:05Z",
				Description:   "Description",
			},
		},
	}

	client := dbClient()
	repo := New(client)
	ctx := context.Background()

	_, err := repo.Create(ctx, oneAuthor)
	assert.Nil(t, err)

	type fields struct {
		ent *gen.Client
	}
	type args struct {
		ctx    context.Context
		author *author.Update
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *gen.Author
		wantErr error
	}{
		{
			name: "update names",
			fields: fields{
				ent: client,
			},
			args: args{
				ctx: ctx,
				author: &author.Update{
					ID:         1,
					FirstName:  "Updated First",
					MiddleName: "Updated Middle",
					LastName:   "Updated Last",
				},
			},
			want: &gen.Author{
				ID:         1,
				FirstName:  "Updated First",
				MiddleName: "Updated Middle",
				LastName:   "Updated Last",
				CreatedAt:  time.Time{},
				UpdatedAt:  time.Time{},
				DeletedAt:  nil,
				Edges: gen.AuthorEdges{
					Books: []*gen.Book{
						{
							Title:         "Title",
							PublishedDate: parseTime.Parse("2022-02-12T15:04:05Z"),
							Description:   "Description",
							CreatedAt:     time.Time{},
							UpdatedAt:     time.Time{},
							DeletedAt:     nil,
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				ent: tt.fields.ent,
			}
			got, err := r.Update(tt.args.ctx, tt.args.author)
			assert.Equal(t, err, tt.wantErr)

			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.FirstName, got.FirstName)
			assert.Equal(t, tt.want.MiddleName, got.MiddleName)
			assert.Equal(t, tt.want.LastName, got.LastName)
			assert.True(t, startTime.Before(got.CreatedAt) || startTime.Equal(got.CreatedAt))
			assert.True(t, startTime.Before(got.UpdatedAt) || startTime.Equal(got.UpdatedAt))
			assert.Nil(t, got.DeletedAt)
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	client := dbClient()
	repo := New(client)
	ctx := context.Background()

	oneAuthor := &author.CreateRequest{
		FirstName:  "First",
		MiddleName: "Middle",
		LastName:   "Last",
		Books: []author.Book{
			{
				Title:         "Title",
				PublishedDate: "2022-02-12T15:04:05Z",
				Description:   "Description",
			},
		},
	}

	created, err := repo.Create(ctx, oneAuthor)
	assert.Nil(t, err)

	type args struct {
		ctx      context.Context
		authorID uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		readErr error
	}{
		{
			name: "delete one",
			args: args{
				ctx:      ctx,
				authorID: 0,
			},
			wantErr: nil,
			readErr: &gen.NotFoundError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, created.ID)
			assert.Equal(t, tt.wantErr, err)

			_, err = repo.Read(ctx, 6)
			assert.NotNil(t, err)
		})
	}
}

func TestRepository_Search(t *testing.T) {}

func dbClient() *gen.Client {
	sqlxDB := sqlx.NewDb(dockerDB.Conn, "postgres")
	drv := entsql.OpenDB(dialect.Postgres, sqlxDB.DB)
	client := gen.NewClient(gen.Driver(drv))
	dockerDB.Ent = client

	return client
}

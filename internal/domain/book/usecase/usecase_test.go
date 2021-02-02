package usecase

import (
	"context"
	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/mock"
	"github.com/gmhafiz/go8/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/now"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"testing"
)

func TestBookUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockRepository(ctrl)
	uc := NewBookUseCase(repo)

	request := book.Request{
		Title:         "title",
		PublishedDate: "2006-01-02 15:04:05 +0000 UTC",
		ImageURL:      "https://example.com/image.png",
		Description:   "",
	}

	ctx := context.Background()

	expected := &models.Book{
		BookID:        0,
		Title:         request.Title,
		PublishedDate: now.MustParse(request.PublishedDate),
		ImageURL: null.String{
			String: request.ImageURL,
			Valid:  true,
		},
		Description: request.Description,
	}
	var err error
	var bookID int64
	repo.EXPECT().Create(ctx, gomock.Eq(expected)).Return(bookID, err).AnyTimes()
	repo.EXPECT().Find(ctx, gomock.Any()).Return(expected, err).AnyTimes()

	bookGot, err := uc.Create(ctx, request)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, bookGot.BookID, 0)
	assert.Equal(t, bookGot.Title, request.Title)
	assert.Equal(t, bookGot.PublishedDate.String(), request.PublishedDate)
	assert.Equal(t, bookGot.Description, request.Description)
	assert.Equal(t, bookGot.ImageURL.String, request.ImageURL)
}

//var (
//	repo book.Test
//)
//
//const uniqueDBName = "usecase"
//
//func TestMain(m *testing.M) {
//	// must go back to project's root path to get to the .env and ./database/migrations/ folder
//	changeDirTo := "../../../../"
//	err := os.Chdir(changeDirTo)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	err = godotenv.Load(".env")
//	if err != nil {
//		log.Println(err)
//	}
//	cfg := configs.DockerTestCfg()
//	cfg.Name = uniqueDBName
//
//	pool, err := dockertest.NewPool("")
//	if err != nil {
//		log.Fatalf("could not connect to docker: %s", err)
//	}
//
//	opts := dockertest.RunOptions{
//		Repository: "postgres",
//		Tag:        "13",
//		Env: []string{
//			"POSTGRES_USER=" + cfg.User,
//			"POSTGRES_PASSWORD=" + cfg.Pass,
//			"POSTGRES_DB=" + cfg.Name,
//			"TZ=UTC",
//			"PG_TZ=UTC",
//		},
//		ExposedPorts: []string{"5432"},
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"5432": {
//				{HostIP: "0.0.0.0", HostPort: cfg.Port},
//			},
//		},
//	}
//	resource, err := pool.RunWithOptions(&opts, func(config *docker.HostConfig) {
//		// set AutoRemove to true so that stopped container goes away by itself
//		config.AutoRemove = true
//		config.RestartPolicy = docker.RestartPolicy{
//			Name: "no",
//		}
//	})
//	if err != nil {
//		log.Fatalln("error running docker container")
//	}
//
//	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
//		cfg.Host, cfg.Port, cfg.User, cfg.Pass, uniqueDBName)
//	// Docker layer network is different on Mac
//	if runtime.GOOS == "darwin" {
//		cfg.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
//	}
//
//	if err = pool.Retry(func() error {
//		db, err := sqlx.Open(cfg.Dialect, dsn)
//		if err != nil {
//			return err
//		}
//		repo = postgres.NewBookRepository(db)
//		return db.Ping()
//	}); err != nil {
//		log.Fatalf("could not connect to docker: %s", err.Error())
//	}
//
//	defer func() {
//		repo.Close()
//	}()
//
//	dbCfg := &configs.Configs{
//		DockerTest: &configs.DockerTest{
//			Driver:  cfg.Dialect,
//			Host:    cfg.Host,
//			Port:    cfg.Port,
//			Name:    uniqueDBName,
//			User:    cfg.User,
//			Pass:    cfg.Pass,
//			SslMode: cfg.SslMode,
//		},
//	}
//	migrate.Up(dbCfg, ".")
//
//	code := m.Run()
//
//	if err := pool.Purge(resource); err != nil {
//		log.Printf("could not purge resource: %s", err)
//	}
//
//	os.Exit(code)
//}
//
//func TestBookUseCase_Create(t *testing.T) {
//	uc := NewBookUseCase(repo)
//
//	request := book.Request{
//		Title:         "title",
//		PublishedDate: "2006-01-02 15:04:05 +0000 UTC",
//		ImageURL:      "https://example.com/image.png",
//		Description:   "",
//	}
//
//	bookGot, err := uc.Create(context.Background(), request)
//	if err != nil {
//		t.Fatal(err)
//	}
//	assert.NotEqual(t, bookGot.BookID, 0)
//	assert.Equal(t, bookGot.Title, request.Title)
//	assert.Equal(t, bookGot.PublishedDate.String(), request.PublishedDate)
//	assert.Equal(t, bookGot.Description, request.Description)
//	assert.Equal(t, bookGot.ImageURL.String, request.ImageURL)
//}

package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	"github.com/gmhafiz/go8/cmd/extmigrate/migrate"
	"github.com/gmhafiz/go8/configs"
	bookRepo "github.com/gmhafiz/go8/internal/domain/book/repository/postgres"
	bookUseCase "github.com/gmhafiz/go8/internal/domain/book/usecase"
	"github.com/gmhafiz/go8/internal/domain/health/repository/postgres"
	"github.com/gmhafiz/go8/internal/domain/health/usecase"
	"github.com/gmhafiz/go8/internal/model"
	"github.com/gmhafiz/go8/internal/resource"
	"github.com/gmhafiz/go8/internal/server"
	"github.com/gmhafiz/go8/third_party/database"
)

type e2eTestSuite struct {
	suite.Suite
	dbConnectionStr string
	port            string
	db              *sqlx.DB
	//dbMigration     *migrate.Migrate
	app          *server.Server
	dockerClient *client.Client
	cfg          *configs.Configs
	Domain       *server.Domain
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {
	_ = os.Chdir("../../")
	cfg := configs.New()
	s.port = cfg.Api.Port
	//s.dockerClient = dbDocker(cfg)

	//cfg.Database.
	s.app = server.New("test")
	s.Init()

	migrate.Up(cfg, ".")
	//go s.app.Run(cfg, "testing")
}

// https://docs.docker.com/engine/api/sdk/examples/
func dbDocker(cfg *configs.Configs) (cl *client.Client) {
	imageName := "postgres:13"

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)

	return nil
}

func (s *e2eTestSuite) TearDownSuite() {
	p, _ := os.FindProcess(syscall.Getpid())
	_ = p.Signal(syscall.SIGINT)
}

func (s *e2eTestSuite) Test_EndToEnd_CreateBook() {
	reqBook := &resource.BookRequest{
		Title:         "title21",
		PublishedDate: "2020-07-31 15:04:05.1235 +0000 UTC",
		ImageURL:      "http://example.com/image.png",
		Description:   "descr",
	}
	reqStr, _ := json.Marshal(reqBook)

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://localhost:%s/api/v1/books", s.port),
		strings.NewReader(string(reqStr)))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusCreated, response.StatusCode)

	byteBody, err := ioutil.ReadAll(response.Body)
	s.NoError(err)

	respBook := model.Book{}
	err = json.Unmarshal(byteBody, &respBook)
	s.NoError(err)

	s.NotZero(respBook.BookID)
	s.Equal(reqBook.Title, respBook.Title)
	s.Equal(reqBook.Description, respBook.Description.String)
	s.Equal(reqBook.PublishedDate, respBook.PublishedDate.String())
	s.Equal(reqBook.ImageURL, respBook.ImageURL.String)

	_ = response.Body.Close()
}

func (s *e2eTestSuite) Init() {
	s.cfg = configs.New()
	db := database.NewSqlx(s.cfg)

	s.Domain.BookUC = bookUseCase.NewBookUseCase(bookRepo.NewBookRepository(db))
	s.Domain.HealthUC = usecase.NewHealthUseCase(postgres.NewHealthRepository(db))
}

package postgres_test

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"github.com/sp4rd4/ports/pkg/domain"
	"github.com/sp4rd4/ports/pkg/storage/postgres"
	"github.com/stretchr/testify/suite"

	// db driver
	_ "github.com/lib/pq"
)

var postgresHost string

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run(
		"postgres", "12.3-alpine",
		[]string{"POSTGRES_DB=ports", "POSTGRES_USER=postgres", "POSTGRES_PASSWORD=postgres"},
	)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	postgresHost = fmt.Sprintf(
		"postgresql://postgres:postgres@localhost:%s/ports?sslmode=disable", resource.GetPort("5432/tcp"),
	)

	if err := pool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", postgresHost)
		if err != nil {
			return err
		}
		defer db.Close()
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to postgres: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

type PostgresTestSuite struct {
	suite.Suite
	db      *sql.DB
	storage postgres.Storage
}

func (s *PostgresTestSuite) SetupSuite() {
	db, err := sql.Open("postgres", postgresHost)
	if err != nil {
		s.T().Fatalf("Unable to connect to postgres: %s", err)
	}
	dbMigrate, err := sql.Open("postgres", postgresHost)
	if err != nil {
		s.T().Fatalf("Unable to connect to postgres: %s", err)
	}
	storage := postgres.New(db)
	err = storage.Migrate(dbMigrate, "migrations")
	if !s.Nil(err, "Should migrate db with no error") {
		s.T().Fatal("Migration failed")
	}
	s.storage = storage
	s.db = db
}

func (s *PostgresTestSuite) TearDownSuite() {
	err := s.db.Close()
	if err != nil {
		s.T().Fatal("DB close failed")
	}
}

func (s *PostgresTestSuite) TearDownTest() {
	_, err := s.db.Exec("TRUNCATE ports;")
	if err != nil {
		s.T().Fatal("Test cleanup failed")
	}
}

func (s *PostgresTestSuite) TestSave() {
	port := &domain.Port{
		ID:      "PORTID",
		Name:    "Port",
		City:    "Boston",
		Country: "Belgium",
		Alias:   domain.StringArray{"PORTIDD"},
		Regions: domain.StringArray{"Provance", "Nova Scotia"},
		Coordinates: domain.Location{
			Latitude:  31.03351,
			Longitude: -17.8251657,
		},
		Province: "",
		Timezone: "Asia/Dubai",
	}
	err := s.storage.Save(port)
	s.Nil(err, "Should save port with no error")

	lPort, err := s.storage.Get(port.ID)
	s.Nil(err, "Should load port with no error")
	s.Equal(port, lPort, "Should load port equal to saved")
}

func (s *PostgresTestSuite) TestConflictingSaves() {
	port := &domain.Port{
		ID:      "PORTID",
		Name:    "Port",
		City:    "Boston",
		Country: "Belgium",
		Alias:   domain.StringArray{"PORTIDD"},
		Regions: domain.StringArray{"Provance", "Nova Scotia"},
		Coordinates: domain.Location{
			Latitude:  31.03351,
			Longitude: -17.8251657,
		},
		Province: "",
		Timezone: "Asia/Dubai",
	}
	err := s.storage.Save(port)
	s.Nil(err, "Should save port with no error")

	port.Name = "New Port"
	port.Country = "France"
	err = s.storage.Save(port)
	s.Nil(err, "Should save port with no error")

	lPort, err := s.storage.Get(port.ID)
	s.Nil(err, "Should load port with no error")
	s.Equal(port, lPort, "Should load port equal to modified")
}

func (s *PostgresTestSuite) TestMultipleSaves() {
	port1 := &domain.Port{
		ID:      "PORTID",
		Name:    "Port",
		City:    "Boston",
		Country: "Belgium",
		Alias:   domain.StringArray{"PORTIDD"},
		Regions: domain.StringArray{"Provance", "Nova Scotia"},
		Coordinates: domain.Location{
			Latitude:  31.03351,
			Longitude: -17.8251657,
		},
		Province: "",
		Timezone: "Asia/Dubai",
	}
	port2 := &domain.Port{
		ID:      "PORTID2",
		Name:    "Porting",
		City:    "Gyor",
		Country: "Slovakia",
		Coordinates: domain.Location{
			Latitude:  51.21303351,
			Longitude: 23.83451657,
		},
		Province: "",
		Timezone: "Asia/Beijing",
	}

	err := s.storage.Save(port1)
	s.Nil(err, "Should save port with no error")

	err = s.storage.Save(port2)
	s.Nil(err, "Should save port with no error")

	lPort, err := s.storage.Get(port1.ID)
	s.Nil(err, "Should load port with no error")
	s.Equal(port1, lPort, "Should load port equal to first")

	lPort, err = s.storage.Get(port2.ID)
	s.Nil(err, "Should load port with no error")
	s.Equal(port2, lPort, "Should load port equal to first")
}

func (s *PostgresTestSuite) TestGetMIssing() {
	lPort, err := s.storage.Get("id")
	s.Nil(lPort, "Should return nil port")
	s.True(errors.Is(err, domain.ErrNotFound), "Should return not found error")
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

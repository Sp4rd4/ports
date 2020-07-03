package postgres

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sp4rd4/ports/pkg/domain"
)

const errorTag = "postgres"

type storage struct {
	db *sqlx.DB
}

type Config struct {
	Host             string `env:"DATABASE_URL,required"`
	MigrationsFolder string `env:"MIGRATIONS_FOLDER" envDefault:"migrations"`
}

func New(conf Config) (domain.PortRepository, error) {
	db, err := sqlx.Connect("postgres", conf.Host)
	if err != nil {
		return nil, fmt.Errorf("%v] sql connect: %w", errorTag, err)
	}

	migrateDB, err := sqlx.Connect("postgres", conf.Host)
	if err != nil {
		return nil, fmt.Errorf("%v] sql migrate connect: %w", errorTag, err)
	}

	driver, err := postgres.WithInstance(migrateDB.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("%v] migrate driver: %w", errorTag, err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+conf.MigrationsFolder, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("%v] new migrate: %w", errorTag, err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("%v] migration: %w", errorTag, err)
	}

	return &storage{db: db}, nil
}

func (s storage) Save(port *domain.Port) error {
	_, err := s.db.NamedExec(`
	INSERT INTO ports (id, name, city, country, alias, regions, coordinates, province, timezone, unlocs, code)
		VALUES (:id, :name, :city, :country, :alias, :regions, :coordinates, :province, :timezone, :unlocs, :code)
	ON CONFLICT (id)
		DO UPDATE SET
			id=:id, name=:name, city=:city, country=:country, alias=:alias, regions=:regions,
			coordinates=:coordinates, province=:province, timezone=:timezone, unlocs=:unlocs, code=:code;`,
		port)
	if err != nil {
		return fmt.Errorf("%v] save: %w", errorTag, err)
	}
	return nil
}

func (s storage) Get(id string) (*domain.Port, error) {
	port := &domain.Port{}
	err := s.db.Get(port, `SELECT * FROM ports WHERE id=$1;`, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%v] get: %w", errorTag, domain.ErrNotFound)
	}
	return port, fmt.Errorf("%v] get: %w", errorTag, err)
}

// Client for ports storage in postgresql db.
package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/sp4rd4/ports/pkg/domain"

	// migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const errorTag = "postgres"

type Storage struct {
	db *sqlx.DB
}

var _ domain.PortRepository = Storage{}

func New(db *sql.DB) Storage {
	return Storage{db: sqlx.NewDb(db, "postgres")}
}

// Migrate creates migrations table if not exists and runs pending migrations,
// db connection will be closed after migrations done.
func (s Storage) Migrate(db *sql.DB, migrationsFolder string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("[%v] migrate driver: %w", errorTag, err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsFolder, "postgres", driver)
	if err != nil {
		return fmt.Errorf("[%v] new migrate: %w", errorTag, err)
	}

	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("[%v] migration: %w", errorTag, err)
	}
	return nil
}

func (s Storage) Save(port *domain.Port) error {
	_, err := s.db.NamedExec(`
	INSERT INTO ports (id, name, city, country, alias, regions, coordinates, province, timezone, unlocs, code)
		VALUES (:id, :name, :city, :country, :alias, :regions, :coordinates, :province, :timezone, :unlocs, :code)
	ON CONFLICT (id)
		DO UPDATE SET
			id=:id, name=:name, city=:city, country=:country, alias=:alias, regions=:regions,
			coordinates=:coordinates, province=:province, timezone=:timezone, unlocs=:unlocs, code=:code;
		`, port)
	if err != nil {
		return fmt.Errorf("[%v] save: %w", errorTag, err)
	}
	return nil
}

func (s Storage) Get(id string) (*domain.Port, error) {
	port := &domain.Port{}
	err := s.db.Get(port, `SELECT * FROM ports WHERE id=$1;`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("[%v] get: %w", errorTag, domain.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("[%v] get: %w", errorTag, err)
	}
	return port, nil
}

package address_repo

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mariocoski/address-service/internal/config"
	domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	"github.com/mariocoski/address-service/internal/shared/database/postgres_driver"
)

type AddressesRepository interface {
	GetAll() ([]*domain.Address, error)
	// GetAllPaginated(page int, pageSize int) ([]*domain.Address, error)
	// GetById(id int) (*domain.Address, error)
	// Save(address *domain.Address) error
}

type AddresssRepoDependencies struct {
	RepoType string
	Config   *config.Config
}

type postgresAddressesRepo struct {
	conn *sql.DB
}

func NewAddressesRepository(dependencies *AddresssRepoDependencies) (AddressesRepository, error) {
	if dependencies.RepoType == "postgres" {
		conn, err := postgres_driver.ConnectSQL(dependencies.Config.PostgresConnectionUrl)
		log.Println("conn", dependencies.Config.PostgresConnectionUrl)
		if err != nil {
			return &postgresAddressesRepo{}, err
		}

		return &postgresAddressesRepo{
			conn: conn.SQL,
		}, nil
	}

	return nil, fmt.Errorf("unknown RepoType, found: %v", dependencies.RepoType)
}

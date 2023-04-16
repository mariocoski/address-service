package address_repo

import (
	"database/sql"
	"fmt"

	"github.com/mariocoski/address-service/internal/config"
	domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	"github.com/mariocoski/address-service/internal/shared/core/pagination"
	"github.com/mariocoski/address-service/internal/shared/database/postgres_driver"
)

type AddressesRepository interface {
	GetAllPaginated(currentPage int, perPage int) (pagination.PaginationResult[domain.Address], error)
	GetById(addressId string) (domain.Address, error)
	// Save(address *domain.Address) error
}

type AddresssRepoDependencies struct {
	Config *config.Config
}

type postgresAddressesRepo struct {
	conn *sql.DB
}

func NewAddressesRepository(config config.Config) (AddressesRepository, error) {
	if config.RepositoryType == "postgres" {
		conn, err := postgres_driver.ConnectSQL(config.PostgresConnectionUrl)
		if err != nil {
			return &postgresAddressesRepo{}, err
		}

		return &postgresAddressesRepo{
			conn: conn.SQL,
		}, nil
	}

	return nil, fmt.Errorf("unknown RepositoryType, found: %v", config.RepositoryType)
}

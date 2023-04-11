package address_repo

import (
	"context"
	"fmt"
	"log"

	addresses "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	addresses_domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
)

func (r *postgresAddressesRepo) GetAll() ([]*domain.Address, error) {
	rows, err := r.conn.QueryContext(
		context.Background(),
		`SELECT 
			id,
			address_line_1,
			address_line_2,
			address_line_3,
			city,
			county,
			state,
			postcode,
			country
		FROM 
			addresses
	`)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	addresses := make([]*addresses.Address, 0)

	for rows.Next() {
		address := &addresses_domain.Address{}

		err := rows.Scan(
			&address.Id,
			&address.AddressLine1,
			&address.AddressLine2,
			&address.AddressLine3,
			&address.City,
			&address.County,
			&address.State,
			&address.Postcode,
			&address.Country,
		)
		log.Println("address", address)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}
	return addresses, nil
}

// func (r *AddressesRepository) GetAllPaginated(page int, pageSize int) ([]*models.Address, error) {
// 	// TODO: Implement pagination logic using OFFSET and LIMIT clauses in SQL query
// 	return nil, fmt.Errorf("not implemented")
// }

// func (r *AddressesRepository) GetById(id int) (*models.Address, error) {
// 	address := &models.Address{}
// 	err := r.conn.QueryRow(context.Background(), "SELECT * FROM addresses WHERE id = $1", id).Scan(
// 		&address.ID,
// 		&address.AddressLine1,
// 		&address.AddressLine2,
// 		&address.City,
// 		&address.County,
// 		&address.State,
// 		&address.Postcode,
// 		&address.Country,
// 	)
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			return nil, fmt.Errorf("address not found")
// 		}
// 		return nil, fmt.Errorf("error querying database: %v", err)
// 	}

// 	return address, nil
// }

// func (r *AddressesRepository) Save(address *models.Address) error {
// 	_, err := r.conn.Exec(
// 		context.Background(),
// 		"INSERT INTO addresses (address_line_1, address_line_2, city, county, state, postcode, country) VALUES ($1, $2, $3, $4, $5, $6, $7)",
// 		address.AddressLine1,
// 		address.AddressLine2,
// 		address.City,
// 		address.County,
// 		address.State,
// 		address.Postcode,
// 		address.Country,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("error inserting address: %v", err)
// 	}

// 	return nil
// }

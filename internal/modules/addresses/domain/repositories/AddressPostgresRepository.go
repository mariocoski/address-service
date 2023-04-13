package address_repo

import (
	"context"
	"fmt"
	"log"

	addresses "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	addresses_domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	domain "github.com/mariocoski/address-service/internal/modules/addresses/domain"
	"github.com/mariocoski/address-service/internal/shared/core/pagination"
)

func (r *postgresAddressesRepo) GetAllPaginated(currentPage int, perPage int) (pagination.PaginationResult[domain.Address], error) {
	offset := (currentPage - 1) * perPage
	limit := perPage

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
		ORDER BY id DESC
		OFFSET $1 LIMIT $2
	`, offset, limit)

	addresses := make([]addresses.Address, 0)
	default_pagination_result := pagination.PaginationResult[domain.Address]{
		Data:       addresses,
		Pagination: pagination.DEFAULT_PAGINATION,
	}

	if err != nil {
		return default_pagination_result, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		address := addresses_domain.Address{}

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
			return default_pagination_result, fmt.Errorf("error scanning row: %v", err)
		}
		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return default_pagination_result, fmt.Errorf("error iterating rows: %v", err)
	}

	var total int

	row := r.conn.QueryRowContext(
		context.Background(),
		`SELECT COUNT(*) FROM (
			SELECT 1 FROM addresses
		) t`)

	err = row.Scan(&total)
	if err != nil {
		return default_pagination_result, fmt.Errorf("error counting addresses rows: %v", err)
	}

	lastPage := total / perPage
	if total%perPage != 0 {
		lastPage++
	}

	var to int

	if currentPage == lastPage {
		to = total
	} else {
		to = offset + len(addresses)
	}

	paginationData := pagination.Pagination{
		LastPage:    lastPage,
		Total:       total,
		CurrentPage: currentPage,
		PerPage:     perPage,
		From:        offset,
		To:          to,
	}

	paginatedResult := pagination.PaginationResult[addresses_domain.Address]{
		Data:       addresses,
		Pagination: paginationData,
	}

	return paginatedResult, nil
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

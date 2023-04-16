package address_repo

import (
	"context"
	"database/sql"
	"fmt"

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
		ORDER BY created_at DESC
		OFFSET $1 LIMIT $2
	`, offset, limit)

	addresses := make([]domain.Address, 0)
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

		var addressLine2, addressLine3, county, state sql.NullString

		err := rows.Scan(
			&address.Id,
			&address.AddressLine1,
			&addressLine2,
			&addressLine3,
			&address.City,
			&county,
			&state,
			&address.Postcode,
			&address.Country,
		)

		if err != nil {
			return default_pagination_result, fmt.Errorf("error scanning row: %v", err)
		}

		if addressLine2.Valid {
			address.AddressLine2 = addressLine2.String
		}

		if addressLine3.Valid {
			address.AddressLine3 = addressLine3.String
		}

		if county.Valid {
			address.County = county.String
		}

		if state.Valid {
			address.State = state.String
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
		to = offset + perPage
	}

	// TODO: investigate the validity of the computations and add tests
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

func (r *postgresAddressesRepo) GetById(addressId string) (domain.Address, error) {

	row := r.conn.QueryRowContext(
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
		WHERE id=$1::uuid
		LIMIT 1
	`, addressId)

	address := addresses_domain.Address{}

	var addressLine2, addressLine3, county, state sql.NullString

	err := row.Scan(
		&address.Id,
		&address.AddressLine1,
		&addressLine2,
		&addressLine3,
		&address.City,
		&county,
		&state,
		&address.Postcode,
		&address.Country,
	)

	if err != nil {
		return domain.Address{}, fmt.Errorf("error get address by id: %v", err)
	}

	if addressLine2.Valid {
		address.AddressLine2 = addressLine2.String
	}

	if addressLine3.Valid {
		address.AddressLine3 = addressLine3.String
	}

	if county.Valid {
		address.County = county.String
	}

	if state.Valid {
		address.State = state.String
	}

	return address, nil
}

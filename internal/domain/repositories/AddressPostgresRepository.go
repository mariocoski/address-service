package address_repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	domain "github.com/mariocoski/address-service/internal/domain"
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
		WHERE deleted_at IS NULL
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
		address := domain.Address{}

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
			SELECT 1 FROM addresses WHERE deleted_at IS NULL
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

	paginatedResult := pagination.PaginationResult[domain.Address]{
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
		WHERE id=$1::uuid AND deleted_at IS NULL
		LIMIT 1
	`, addressId)

	address := domain.Address{}

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
		return domain.Address{}, domain.ErrAddressNotFound
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

func (r *postgresAddressesRepo) Save(addressDBO domain.AddressInitializer) (domain.Address, error) {

	row := r.conn.QueryRowContext(
		context.Background(),
		`INSERT INTO addresses (  
			address_line_1,
			address_line_2,
			address_line_3,
			city,
			county,
			state,
			postcode,
			country
		)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		addressDBO.AddressLine1,
		addressDBO.AddressLine2,
		addressDBO.AddressLine3,
		addressDBO.City,
		addressDBO.County,
		addressDBO.State,
		addressDBO.Postcode,
		addressDBO.Country,
	)

	var addressId string
	err := row.Scan(
		&addressId,
	)

	if err != nil {
		return domain.Address{}, fmt.Errorf("error scanning address id returned by the save method: %v", err)
	}

	savedAddress, err := r.GetById(addressId)
	if err != nil {
		return domain.Address{}, fmt.Errorf("error retrieving saved address: %v", err)
	}

	return savedAddress, nil
}

func (r *postgresAddressesRepo) Delete(addressId string) (string, error) {
	result, err := r.conn.ExecContext(
		context.Background(),
		`UPDATE addresses SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL
	`, addressId)

	if err != nil {
		return "", fmt.Errorf("error in delete method when soft deleting: %v", err)
	}
	numberOfAffectedRows, err := result.RowsAffected()

	if err != nil {
		return "", fmt.Errorf("error in delete method when checking RowsAffected(): %v", err)
	}

	if numberOfAffectedRows != 1 {
		return "", domain.ErrAddressNotFound
	}

	return addressId, nil
}

type UpdateQuery = string

func getUpdateQueryAndParams(addressId string, addressPatch domain.AddressPatch) (UpdateQuery, pgx.NamedArgs) {

	newPatch := pgx.NamedArgs{}
	var query UpdateQuery = "UPDATE addresses SET updated_at = NOW()"
	if addressPatch.AddressLine1 != nil {
		newPatch["address_line_1"] = *addressPatch.AddressLine1
		query += fmt.Sprintf(", %v = %v", "address_line_1", "@address_line_1")
	}
	if addressPatch.AddressLine2 != nil {
		newPatch["address_line_2"] = *addressPatch.AddressLine2
		query += fmt.Sprintf(", %v = %v", "address_line_2", "@address_line_2")
	}
	if addressPatch.AddressLine3 != nil {
		newPatch["address_line_3"] = *addressPatch.AddressLine3
		query += fmt.Sprintf(", %v = %v", "address_line_3", "@address_line_3")
	}
	if addressPatch.City != nil {
		newPatch["city"] = *addressPatch.City
		query += fmt.Sprintf(", %v = %v", "city", "@city")
	}
	if addressPatch.County != nil {
		newPatch["county"] = *addressPatch.County
		query += fmt.Sprintf(", %v = %v", "county", "@county")

	}
	if addressPatch.Country != nil {
		newPatch["country"] = *addressPatch.Country
		query += fmt.Sprintf(", %v = %v", "country", "@country")
	}
	if addressPatch.State != nil {
		newPatch["state"] = *addressPatch.State
		query += fmt.Sprintf(", %v = %v", "state", "@state")
	}
	if addressPatch.Postcode != nil {
		newPatch["postcode"] = *addressPatch.Postcode
		query += fmt.Sprintf(", %v = %v", "postcode", "@postcode")
	}
	newPatch["id"] = addressId
	query += " WHERE id = @id::uuid AND deleted_at IS NULL"

	return query, newPatch
}

func (r *postgresAddressesRepo) Update(addressId string, addressPatch domain.AddressPatch) (domain.Address, error) {
	query, params := getUpdateQueryAndParams(addressId, addressPatch)

	result, err := r.conn.ExecContext(context.Background(), query, params)
	if err != nil {
		return domain.Address{}, fmt.Errorf("error in update method: %v", err)
	}

	numberOfAffectedRows, err := result.RowsAffected()
	if err != nil {
		return domain.Address{}, fmt.Errorf("error in update method when checking RowsAffected(): %v", err)
	}

	if numberOfAffectedRows != 1 {
		return domain.Address{}, domain.ErrAddressNotFound
	}

	updatedAddress, err := r.GetById(addressId)
	if err != nil {
		return domain.Address{}, fmt.Errorf("error in update method when fetching updated address: %v", err)
	}

	return updatedAddress, nil
}

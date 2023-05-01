package address_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	domain "github.com/mariocoski/address-service/internal/domain"
	"github.com/mariocoski/address-service/internal/shared/core/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *mongoAddressesRepo) GetAllPaginated(currentPage int, perPage int) (pagination.PaginationResult[domain.Address], error) {
	// Calculate the offset based on the current page and the number of items per page.
	offset := (currentPage - 1) * perPage

	// Set the limit to the number of items per page.
	limit := perPage

	// Build the MongoDB filter to exclude soft deleted addresses.
	filter := bson.D{{Key: "deleted_at", Value: nil}}

	// Set the options for the MongoDB query.
	findOptions := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	// Execute the MongoDB query.
	cursor, err := r.collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return pagination.PaginationResult[domain.Address]{}, fmt.Errorf("error querying database: %v", err)
	}
	defer cursor.Close(context.Background())

	// Read the addresses from the MongoDB cursor.
	addresses := make([]domain.Address, 0)
	for cursor.Next(context.Background()) {
		address := domain.Address{}
		err := cursor.Decode(&address)
		if err != nil {
			return pagination.PaginationResult[domain.Address]{}, fmt.Errorf("error decoding address: %v", err)
		}
		addresses = append(addresses, address)
	}
	if err := cursor.Err(); err != nil {
		return pagination.PaginationResult[domain.Address]{}, fmt.Errorf("error iterating addresses: %v", err)
	}

	// Count the total number of addresses.
	totalInt64, err := r.collection.CountDocuments(context.Background(), filter, nil)
	if err != nil {
		return pagination.PaginationResult[domain.Address]{}, fmt.Errorf("error counting addresses: %v", err)
	}
	total := int(totalInt64)
	// Calculate the pagination information.
	lastPage := total / perPage
	if total%perPage != 0 {
		lastPage++
	}
	from := offset + 1
	to := offset + len(addresses)
	if to > total {
		to = total
	}
	paginationData := pagination.Pagination{
		LastPage:    lastPage,
		Total:       int(total),
		CurrentPage: currentPage,
		PerPage:     perPage,
		From:        from,
		To:          to,
	}

	// Build and return the pagination result.
	paginatedResult := pagination.PaginationResult[domain.Address]{
		Data:       addresses,
		Pagination: paginationData,
	}
	return paginatedResult, nil
}

func (r *mongoAddressesRepo) GetById(addressId string) (domain.Address, error) {
	uuid, err := uuid.Parse(addressId)

	if err != nil {
		return domain.Address{}, err
	}

	filter := bson.M{
		"_id":       uuid,
		"deletedAt": nil,
	}

	var address domain.Address
	if err := r.collection.FindOne(context.Background(), filter).Decode(&address); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Address{}, domain.ErrAddressNotFound
		}

		return domain.Address{}, fmt.Errorf("error retrieving address by id: %v", err)
	}

	return address, nil
}

func (r *mongoAddressesRepo) Save(addressDBO domain.AddressInitializer) (domain.Address, error) {

	var uuid uuid.UUID = uuid.New()

	document := bson.M{
		"_id":            uuid,
		"address_line_1": addressDBO.AddressLine1,
		"address_line_2": addressDBO.AddressLine2,
		"address_line_3": addressDBO.AddressLine3,
		"city":           addressDBO.City,
		"county":         addressDBO.County,
		"state":          addressDBO.State,
		"postcode":       addressDBO.Postcode,
		"country":        addressDBO.Country,
		"created_at":     time.Now().Format(time.RFC3339),
		"updated_at":     time.Now().Format(time.RFC3339),
		"deleted_at":     nil,
	}

	_, err := r.collection.InsertOne(context.Background(), document)
	if err != nil {
		return domain.Address{}, fmt.Errorf("error saving address: %v", err)
	}

	createdAddress, err := r.GetById(uuid.String())
	if err != nil {
		return domain.Address{}, fmt.Errorf("error in save when fetching updated address: %v", err)
	}

	return createdAddress, nil
}

func (r *mongoAddressesRepo) Delete(addressId string) (string, error) {
	uuid, err := uuid.Parse(addressId)

	if err != nil {
		return "", err
	}

	filter := bson.M{
		"_id":        uuid,
		"deleted_at": nil,
	}

	result, err := r.collection.UpdateOne(context.Background(), filter, bson.M{"$set": bson.M{"deleted_at": time.Now().UTC()}})
	if err != nil {
		return "", fmt.Errorf("error in delete method when soft deleting: %v", err)
	}

	if result.ModifiedCount == 0 {
		return "", domain.ErrAddressNotFound
	}

	return addressId, nil
}

func (r *mongoAddressesRepo) Update(addressId string, addressPatch domain.AddressPatch) (domain.Address, error) {
	uuid, err := uuid.Parse(addressId)

	if err != nil {
		return domain.Address{}, err
	}

	filter := bson.M{"_id": uuid}
	update := bson.M{
		"$set": bson.M{
			"address_line_1": addressPatch.AddressLine1,
			"address_line_2": addressPatch.AddressLine2,
			"address_line_3": addressPatch.AddressLine3,
			"city":           addressPatch.City,
			"county":         addressPatch.County,
			"state":          addressPatch.State,
			"postcode":       addressPatch.Postcode,
			"country":        addressPatch.Country,
			"updated_at":     time.Now().Format(time.RFC3339),
		},
	}

	result, err := r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return domain.Address{}, fmt.Errorf("error in update method: %v", err)
	}

	if result.MatchedCount != 1 {
		return domain.Address{}, domain.ErrAddressNotFound
	}

	updatedAddress, err := r.GetById(addressId)
	if err != nil {
		return domain.Address{}, fmt.Errorf("error in update method when fetching updated address: %v", err)
	}

	return updatedAddress, nil
}

package address_repo

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/mariocoski/address-service/internal/config"
	domain "github.com/mariocoski/address-service/internal/domain"
	"github.com/mariocoski/address-service/internal/shared/core/pagination"
	"github.com/mariocoski/address-service/internal/shared/database/postgres_driver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AddressesRepository interface {
	GetAllPaginated(currentPage int, perPage int) (pagination.PaginationResult[domain.Address], error)
	GetById(addressId string) (domain.Address, error)
	Save(address domain.AddressInitializer) (domain.Address, error)
	Delete(addressId string) (string, error)
	Update(addressId string, patch domain.AddressPatch) (domain.Address, error)
}

type AddresssRepoDependencies struct {
	Config *config.Config
}

type postgresAddressesRepo struct {
	conn *sql.DB
}

type mongoAddressesRepo struct {
	collection *mongo.Collection
}

var (
	tUUID       = reflect.TypeOf(uuid.UUID{})
	uuidSubtype = byte(0x04)

	mongoRegistry = bson.NewRegistryBuilder().
			RegisterTypeEncoder(tUUID, bsoncodec.ValueEncoderFunc(uuidEncodeValue)).
			RegisterTypeDecoder(tUUID, bsoncodec.ValueDecoderFunc(uuidDecodeValue)).
			Build()
)

func uuidEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tUUID {
		return bsoncodec.ValueEncoderError{Name: "uuidEncodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}
	b := val.Interface().(uuid.UUID)
	return vw.WriteBinaryWithSubtype(b[:], uuidSubtype)
}

func uuidDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tUUID {
		return bsoncodec.ValueDecoderError{Name: "uuidDecodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}

	var data []byte
	var subtype byte
	var err error
	switch vrType := vr.Type(); vrType {
	case bsontype.Binary:
		data, subtype, err = vr.ReadBinary()
		if subtype != uuidSubtype {
			return fmt.Errorf("unsupported binary subtype %v for UUID", subtype)
		}
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a UUID", vrType)
	}

	if err != nil {
		return err
	}
	uuid2, err := uuid.FromBytes(data)
	if err != nil {
		return err
	}
	val.Set(reflect.ValueOf(uuid2))
	return nil
}

func NewAddressesRepository(config config.Config) (AddressesRepository, error) {

	if config.RepositoryType == "mongo" {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.MongoDBConnectionUrl).SetRegistry(mongoRegistry))
		if err != nil {
			return &mongoAddressesRepo{}, err
		}

		return &mongoAddressesRepo{
			collection: client.Database(config.MongoDBName).Collection("addresses"),
		}, nil
	}

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

--liquibase formatted sql
--changeset user:20230414145234-create_addresses_table splitStatements:false
CREATE TABLE addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    address_line_1 VARCHAR(255) NOT NULL,
    address_line_2 VARCHAR(255),
    address_line_3 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    county VARCHAR(100),
    state VARCHAR(100),
    postcode VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TRIGGER update_address_updated_at_trigger
BEFORE UPDATE ON addresses
FOR EACH ROW
EXECUTE FUNCTION update_address_updated_at_column();

--rollback DROP table addresses;


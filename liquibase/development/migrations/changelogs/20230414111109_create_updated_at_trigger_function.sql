--liquibase formatted sql
--changeset user:20230414192409-create_updated_at_trigger_function splitStatements:false

CREATE EXTENSION IF NOT EXISTS pgcrypto; 

CREATE OR REPLACE FUNCTION update_address_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

begin;
 
--rollback DROP FUNCTION update_address_updated_at_column();



--liquibase formatted sql
--changeset user:20230414150000-insert_addresses
INSERT INTO addresses (
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
('123 Main St', '', '', 'Los Angeles', 'Los Angeles County', 'California', '90001', 'USA'),
('456 Elm St', '', '', 'San Francisco', 'San Francisco County', 'California', '94101', 'USA'),
('789 Oak St', '', '', 'San Diego', 'San Diego County', 'California', '92101', 'USA'),
('10 Pine St', '', '', 'Boston', 'Suffolk County', 'Massachusetts', '02108', 'USA'),
('20 Maple St', '', '', 'New York', 'New York County', 'New York', '10001', 'USA'),
('30 Cedar St', '', '', 'Chicago', 'Cook County', 'Illinois', '60601', 'USA'),
('40 Birch St', '', '', 'Seattle', 'King County', 'Washington', '98101', 'USA'),
('50 Willow St', '', '', 'Portland', 'Multnomah County', 'Oregon', '97201', 'USA'),
('60 Spruce St', '', '', 'Austin', 'Travis County', 'Texas', '78701', 'USA'),
('70 Walnut St', '', '', 'Denver', 'Denver County', 'Colorado', '80202', 'USA'),
('80 Chestnut St', '', '', 'Philadelphia', 'Philadelphia County', 'Pennsylvania', '19107', 'USA'),
('90 Ash St', '', '', 'Phoenix', 'Maricopa County', 'Arizona', '85001', 'USA'),
('100 Elm St', '', '', 'Las Vegas', 'Clark County', 'Nevada', '89101', 'USA'),
('110 Oak St', '', '', 'Houston', 'Harris County', 'Texas', '77002', 'USA'),
('120 Pine St', '', '', 'San Antonio', 'Bexar County', 'Texas', '78205', 'USA'),
('130 Maple St', '', '', 'Orlando', 'Orange County', 'Florida', '32801', 'USA'),
('140 Cedar St', '', '', 'Atlanta', 'Fulton County', 'Georgia', '30303', 'USA'),
('150 Birch St', '', '', 'New Orleans', 'Orleans Parish', 'Louisiana', '70112', 'USA'),
('160 Willow St', '', '', 'Seattle', 'King County', 'Washington', '98101', 'USA'),
('170 Spruce St', '', '', 'Dallas', 'Dallas County', 'Texas', '75201', 'USA'),
('180 Walnut St', '', '', 'San Diego', 'San Diego County', 'California', '92101', 'USA'),
('190 Chestnut St', '', '', 'Miami', 'Miami-Dade County', 'Florida', '33101', 'USA'),
('200 Ash St', '', '', 'Nashville', 'Davidson County', 'Tennessee', '37201', 'USA'),
('210 Elm St', '', '', 'San Francisco', 'San Francisco County', 'California', '94101', 'USA'),
('220 Oak St', '', '', 'Boston', 'Suffolk County', 'Massachusetts', '02108', 'USA'),
('230 Pine St', '', '', 'Chicago', 'Cook County', 'Illinois', '60601', 'USA'),
('240 Maple St', '', '', 'New York', 'New York County', 'New York', '10001', 'USA');

--rollback DELETE FROM addresses;


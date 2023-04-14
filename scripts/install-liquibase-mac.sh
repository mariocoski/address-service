#!/bin/bash

# Set installation path
INSTALL_PATH=${INSTALL_PATH:-"/usr/local/bin"}

# Check if brew is installed
if ! command -v brew &> /dev/null
then
    echo "Brew not found, installing..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

# Install liquibase with PostgreSQL support
brew install liquibase
wget https://github.com/liquibase/liquibase-postgresql/releases/download/liquibase-postgresql-4.12.0/liquibase-postgresql-4.12.0.jar
sudo mkdir -p /usr/local/Cellar/liquibase/lib/
sudo mv liquibase-postgresql-4.12.0.jar /usr/local/Cellar/liquibase/lib/

echo "Liquibase with PostgreSQL support installed successfully!"
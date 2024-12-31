#! /bin/bash

mysql -h localhost -u root < ./create-db.sql &&
mysql -u root < ./create-user.sql &&
mysql -h localhost -u root snippetbox < ../internal/models/db-scripts/setup.sql 
mysql -h localhost -u root snippetbox < ../internal/models/db-scripts/test-data/dummy-data.sql

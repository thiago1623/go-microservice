@echo off

REM CSV file path on your local machine
set "csv_file=C:\Users\thiago\Documents\Golang\go-microservice\test_bia.csv"  REM Replace this path with your own CSV file path

REM Copy csv file to db container (assuming the container name is "db_go_testing")
docker cp "%csv_file%" db_go_testing:/test_bia.csv

REM Run command inside db container to create tables and import data
docker exec -i db_go_testing psql -d microservice_db_testing -U db_user_testing -p 5432 -h localhost <<EOF
CREATE TABLE addresses (
    id UUID PRIMARY KEY,
    meter_id INT NOT NULL,
    address VARCHAR(255),
    start_date DATE,
    end_date DATE
);

CREATE TABLE energy_consumptions (
    id UUID PRIMARY KEY,
    meter_id INT NOT NULL,
    active_energy DOUBLE PRECISION,
    reactive_energy DOUBLE PRECISION,
    capacitive_reactive DOUBLE PRECISION,
    solar DOUBLE PRECISION,
    date DATE,
    address_id UUID REFERENCES addresses(id)
);

COPY energy_consumptions(id, meter_id, active_energy, reactive_energy, capacitive_reactive, solar, date)
FROM '/test_bia.csv' WITH (FORMAT csv, DELIMITER ',', HEADER true, NULL 'null');
\q
EOF

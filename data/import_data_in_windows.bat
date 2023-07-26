@echo off

REM Replace this path with your own CSV file path
set "csv_file=C:\Users\thiago\Documents\Golang\go-microservice\test_bia.csv"

REM copy csv file to db container
docker cp "%csv_file%" db_go:/test_bia.csv

REM run command inside db container and save data
docker exec -i db_go psql -d microservice_db -U db_user_admin -p 5432 -h localhost << EOF
COPY energy_consumptions(id, meter_id, active_energy, reactive_energy, capacitive_reactive, solar, date)
FROM '/test_bia.csv' WITH (FORMAT csv, DELIMITER ',', HEADER true, NULL 'null');
\q
EOF
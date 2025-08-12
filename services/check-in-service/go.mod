module check-in-service

go 1.24

replace smart-fit/proto => ../../services/db-gateway-service/proto

require (
	github.com/gorilla/mux v1.8.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
)

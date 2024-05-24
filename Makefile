sqlup:
	migrate -path migrate -database "postgresql://root:root@localhost:5432/root?sslmode=disable" -verbose up
sqldown:
	migrate -path migrate -database "postgresql://root:root@localhost:5432/root?sslmode=disable" -verbose down
.PHONY: sqlup sqldown
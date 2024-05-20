.PHONY: create-keypair migrate-create migrate-up migrate-down migrate-force 

PWD = $(shell pwd)
WIN_PWD = $(shell powershell -Command "(Get-Location).Path")
ACCTPATH = $(PWD)/account
MPATH_WIN = $(WIN_PWD)\account\migrations 
MPATH = $(PWD)/account/migrations

PORT = 5433

N = 1

to_unix_path = $(subst \,/,$1)

create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(ACCTPATH)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(ACCTPATH)/rsa_private_$(ENV).pem -pubout -out $(ACCTPATH)/rsa_public_$(ENV).pem

migrate-create:
	@echo "---Creating migration files---"
	migrate create -ext sql -dir $(MPATH_WIN) -seq -digits 5 $(NAME)

migrate-up:
	migrate -source "file://$(call to_unix_path,$(MPATH_WIN))" -database "postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable" up $(N)

migrate-down:
	migrate -source "file://$(call to_unix_path,$(MPATH_WIN))" -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable down $(N)

migrate-force:
	migrate -source file:/$(MPATH) -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable force $(VERSION)
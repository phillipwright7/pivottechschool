run-server:
	go run /Users/phillipwright/go/src/github.com/phillipwright7/pivottechschool/cmd/server/main.go

get-products:
	curl -i -X GET http://localhost:8080/products

get-product:
	curl -i -X GET http://localhost:8080/products/56

add-product:
	curl -i -X POST http://localhost:8080/products -d '{"name": "Apple", "description": "This is a test.", "price": 3}'

update-product:
	curl -i -X PUT http://localhost:8080/products/99 -d '{"name": "New iPhone", "description": "This is a new product.", "price": 749}'

delete-product:
	curl -i -X DELETE http://localhost:8080/products/99

run-marvel:
	cd cmd/marvel && go run main.go -p=../../.env

run-seeder:
	cd cmd/seeder && go run main.go -p=../../products.db

calculator-build:
	cd cmd/calculator
	go build -o calculator

calculator-test:
	cd calculator
	go test -v ./...

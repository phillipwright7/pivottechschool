run-server:
	cd cmd/server && go run main.go -db=../seeder/products.db

get-products:
	curl -i -X GET http://localhost:8080/products?limit=0 && printf "\n"

	curl -i -X GET http://localhost:8080/products?limit=-5 && printf "\n"

	curl -i -X GET http://localhost:8080/products?limit=a && printf "\n"

	curl -i -X GET http://localhost:8080/products?limit=99 && printf "\n"

get-product:
	curl -i -X GET http://localhost:8080/products/0 && printf "\n"

	curl -i -X GET http://localhost:8080/products/-99 && printf "\n"

	curl -i -X GET http://localhost:8080/products/b && printf "\n"

	curl -i -X GET http://localhost:8080/products/98 && printf "\n"

add-product:
	curl -i -X POST http://localhost:8080/products -d '{"price": 2}' && printf "\n"

	curl -i -X POST http://localhost:8080/products -d '{"name": "Apple"}' && printf "\n"

	curl -i -X POST http://localhost:8080/products -d '{"name": "Apple", "price": hfhgdg}' && printf "\n"

	curl -i -X POST http://localhost:8080/products -d '{"name": "Apple", "price": "hfhgdg"}' && printf "\n"

	curl -i -X POST http://localhost:8080/products -d '{"name": "Apple", "price": 2}' && printf "\n"

update-product:
	curl -i -X PUT http://localhost:8080/products/10000 -d '{"name": "iPhone 14", "price": 799}' && printf "\n"

	curl -i -X PUT http://localhost:8080/products/b -d '{"name": "iPhone 14", "price": 799}' && printf "\n"

	curl -i -X PUT http://localhost:8080/products/100 -d '{"name": "iPhone 14"}' && printf "\n"

	curl -i -X PUT http://localhost:8080/products/100 -d '{"price": 799}' && printf "\n"

	curl -i -X PUT http://localhost:8080/products/100 -d '{"name": "iPhone 14", "price": 799}' && printf "\n"

delete-product:
	curl -i -X DELETE http://localhost:8080/products/10000 && printf "\n"

	curl -i -X DELETE http://localhost:8080/products/c && printf "\n"

	curl -i -X DELETE http://localhost:8080/products/99 && printf "\n"

run-marvel:
	cd cmd/marvel && go run main.go -p=../../.env

run-seeder:
	cd cmd/seeder && go run main.go -db=../../products.db -json=../server/products.json

calculator-build:
	cd cmd/calculator && go build -v

calculator-test:
	cd calculator && go test -v ./...

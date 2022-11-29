run-server:
	cd cmd/server && go run main.go -db=../seeder/products.db

get-products:
	curl -i -X GET http://localhost:8080/products?limit=10

get-product:
	curl -i -X GET http://localhost:8080/products/99

add-product:
	curl -i -X POST http://localhost:8080/products -d '{"name": "Apple", "price": 2}'

update-product:
	curl -i -X PUT http://localhost:8080/products/5 -d '{"name": "iPhone 14", "price": 799}'

delete-product:
	curl -i -X DELETE http://localhost:8080/products/100

run-marvel:
	cd cmd/marvel && go run main.go -p=../../.env

run-seeder:
	cd cmd/seeder && go run main.go -db=products.db -json=../server/products.json

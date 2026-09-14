.PHONY: test clean generate

test:
	./test.sh

clean:
	docker compose down -v --remove-orphans

generate:
	sqlc generate
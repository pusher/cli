.PHONY: build-mockgen
build-mockgen:
	docker build -f mockgen.Dockerfile --tag pusher_cli_mockgen .

.PHONY: mocks
mocks: build-mockgen
	docker run --rm --volume "$$(pwd):/src" pusher_cli_mockgen

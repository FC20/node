all: lint install

install: go.sum
		GO111MODULE=on go install ./cmd/ouroborosd
		GO111MODULE=on go install ./cmd/ouroboroscli

		mkdir -p ~/.ouroborosd/config

		cp -r ./installation/genesis.json ~/.ouroborosd/config/
		cp -r ./installation/config.toml ~/.ouroborosd/config/

		ouroboroscli config chain-id ouroboros
		ouroboroscli config output json
		ouroboroscli config indent true
		ouroboroscli config node tcp://localhost:26657
		ouroboroscli config trust-node true

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify
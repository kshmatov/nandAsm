MAIN=ntasm
ifeq ($(OS),Windows_NT)
    BIN=./bin/$(MAIN).exe
else
    BIN=./bin/$(MAIN)
endif

SAMPLE1=samples/Add.asm
SAMPLE2=samples/Rect.asm
SAMPLE_PONG=samples/Pong.asm

.PHONY: ALL
all: deps test build

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: vet
	go test ./...

.PHONY: build
build: vet
	go build -o $(BIN) ./cmd/main.go

.PHONY: deps
deps:
	go mod tidy

init:
	go mod init

run: build
	$(BIN) -i $(SAMPLE1)
	$(BIN) -i $(SAMPLE2)

run-pong: build
	$(BIN) -i $(SAMPLE_PONG) -o pong.hack
include goarch.mk

.PHONY: all install clean

LD_FLAGS := "-s -w"
BINARY_NAME := whippet-cli
INSTALL_PATH := /usr/bin

all: $(BINARY_NAME)

$(BINARY_NAME): 
	mkdir -p .go_work/build/bin
	go build -ldflags $(LD_FLAGS) -o $(BINARY_NAME) .

install: $(BINARY_NAME)
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)/whippet

test:
	go test . ./...

cover:
	rm /tmp/test-cover || true
	go test . ./... -coverprofile=/tmp/test-cover
	go tool cover -html /tmp/test-cover

clean:
	rm -f $(BINARY_NAME)

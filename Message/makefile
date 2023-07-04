.PHONY: all clean build run

# Компилятор Go
GO = go

# Флаги компиляции
GOFLAGS = -v

# Имя исполняемого файла
BINARY = messenger

all: build

build: $(BINARY)

$(BINARY):
	$(GO) build $(GOFLAGS) -o $(BINARY) ./main.go

run:
	$(GO) run $(GOFLAGS) main.go

clean:
	rm -f $(BINARY)

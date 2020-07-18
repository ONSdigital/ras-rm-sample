projectName := file-uploader

.PHONY: install
install: checkgo clean link
	go build -i -v -o main
	chmod 755 main

.PHONY: checkgo
checkgo:
ifndef GOPATH
	$(error GOPATH is not set)
endif

.PHONY: clean
clean:
	rm -f main

.PHONY: link
link:
	mkdir -p $(GOPATH)/src/github.com/ONSdigital
	ln -sf $(shell pwd) $(GOPATH)/src/github.com/ONSdigital/

.PHONY: test
test:
	go test -parallel=1 -race -coverprofile=coverage.txt ./...

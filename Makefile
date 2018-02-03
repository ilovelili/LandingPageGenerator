VERSION = $(shell grep 'version =' version.go | sed -E 's/.*"(.+)"$$/\1/')

default: all

all: build

deps:
	go get -d -v -u github.com/jlaffaye/ftp
	go get -d -v -u github.com/skip2/go-qrcode	
	go get -d -v -u github.com/gocarina/gocsv	

build: deps
	go build -o lpg
	cp -r config.json template output ./ansible/playbooks/files/compiled
	mv lpg ./ansible/playbooks/files/compiled
	rm ./ansible/playbooks/files/compiled/index.html

version:
	@echo $(VERSION)

.PTHONY: all deps build version
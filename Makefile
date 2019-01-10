OUT := ns
PKG := gitlab.havana/BDIO/notification-service-v2
DEB_TEMPLATE := notification-service*.deb

PL_VERSION := ${CI_PIPELINE_ID}.${SHA_VER}
ifeq ("${PL_VERSION}", ".")
	PL_VERSION:=manual
endif
VERSION := $(shell cat VERSION).${PL_VERSION}_($(shell date +%d.%m.%Y))
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: build

build: clean
	go build -i -v -o ${OUT} -ldflags="-X main.version=${VERSION}"

deb: build
	sh packaging/packaging.sh

# test:
# 	@go test -short ${PKG_LIST}

# vet:
# 	@go vet ${PKG_LIST}

# lint:
# 	@for file in ${GO_FILES} ;  do \
# 		golint $$file ; \
# 	done

# vendor:
# 	@godep save

# static: vet lint
# 	go build -i -v -o ${OUT}-static -tags netgo -ldflags="-extldflags \"-static\" -X main.version=${VERSION} -w -s -X main.version=${VERSION}"
# 	go build -i -v -o ${SHITLIB_PKG}/${SHITLIB_OUT}-static -tags netgo -ldflags="-extldflags \"-static\" -X main.version=${VERSION} -w -s -X main.version=${VERSION}" ${SHITLIB_PKG}

# run: server
# 	./${OUT}

deploy: build deb deploy54

# deployvirt: build deb deploy101 deploy102 deploy201

deploy54:
	sh -c "scp ${DEB_TEMPLATE} astra@10.46.2.54:/tmp/"
	ssh astra@10.46.2.54 "sudo dpkg -i /tmp/${DEB_TEMPLATE} && rm -f /tmp/${DEB_TEMPLATE}; sudo apt-get install -f -y"

# deploy101:
# 	sh -c "scp ${DEB_TEMPLATE} astra@192.168.100.101:/tmp/"
# 	ssh astra@192.168.100.101 "sudo dpkg -i /tmp/${DEB_TEMPLATE} && rm -f /tmp/${DEB_TEMPLATE}; sudo apt-get install -f -y"

# deploy102:
# 	sh -c "scp ${DEB_TEMPLATE} astra@192.168.100.102:/tmp/"
# 	ssh astra@192.168.100.102 "sudo dpkg -i /tmp/${DEB_TEMPLATE} && rm -f /tmp/${DEB_TEMPLATE}; sudo apt-get install -f -y "

# deploy201:
# 	sh -c "scp ${DEB_TEMPLATE} astra@192.168.100.201:/tmp/"
# 	ssh astra@192.168.100.201 "sudo dpkg -i /tmp/${DEB_TEMPLATE} && rm -f /tmp/${DEB_TEMPLATE}; sudo apt-get install -f -y"

clean:
	-@rm ${OUT} ${DEB_TEMPLATE}

# .PHONY: run server static vet lint
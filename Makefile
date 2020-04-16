REALIZE := $(GOPATH)/bin/realize
DEP := $(GOPATH)/bin/dep
GQL_GEN := $(GOPATH)/bin/gqlgen
COBRA := $(GOPATH)/bin/cobra

start:
	realize start

build-binary:
	 CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lingva-server ./src/lingva

build:
	docker build -t lingva-server .

remove-container:
	docker rm --force lingva-server || true

run:
	docker run -d \
	-p 80:8080 \
	-v student-images:/files \
	--name lingva-server \
	--restart unless-stopped \
	lingva-server

gqlgen:
	cd src/lingva/gql && gqlgen

clean: 
	rm -rf src/github.com src/golang.org src/gopkg.in pkg

all:
	make gqlgen
	make build-binary
	make build
	make remove-container
	make run

$(DEP):
	@echo Downloading dep.
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

$(REALIZE):
	@echo Downloading realize.
	go get github.com/oxequa/realize

$(GQL_GEN):
	@echo Downloading gqlgen.
	go get github.com/vektah/gqlgen

$(COBRA): 
	@echo Downloading cobra.	
	go get -u github.com/spf13/cobra/cobra
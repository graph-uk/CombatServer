FROM golang:1.13.7
WORKDIR /go
RUN PATH=$PATH:~/go/bin &&\
	apt-get update &&\
	curl -sL https://deb.nodesource.com/setup_10.x | bash -&&\
	apt-get install -y nodejs
COPY . /go
RUN npm install
RUN npm run iter
#	npm install
#	mkdir -p /go/src/github.com/graph-uk &&\
#	cd /go/src/github.com/graph-uk &&\
#	git clone https://github.com/graph-uk/malibu-server.git &&\
#	go get -u github.com/gobuffalo/packr/packr &&\
#	apt-get update &&\
#	curl -sL https://deb.nodesource.com/setup_10.x | bash -&&\
#	apt-get install -y nodejs &&\
#	cd /go/src/github.com/graph-uk/malibu-server &&\
#	GOPATH=$GOPATH:/go/src/github.com/graph-uk/malibu-server &&\
#	npm install &&\

#	npm run build-assets &&\
#	cd /go/src/github.com/graph-uk/malibu-server/src/malibu && go build &&\
#	go install &&\
#	cp /go/bin/malibu . &&\

#	cd /go/src/github.com/graph-uk/malibu-server/src/malibu-client && go build &&\
#	cd /go/src/github.com/graph-uk/malibu-server/src/malibu-worker && go build &&\

#	cd /go/src/github.com/graph-uk/malibu-server/src/malibu && GOOS=windows GOARCH=amd64 go build &&\
#	cd /go/src/github.com/graph-uk/malibu-server/src/malibu-client && GOOS=windows GOARCH=amd64 go build &&\
#	cd /go/src/github.com/graph-uk/malibu-server/src/malibu-worker && GOOS=windows GOARCH=amd64 go build &&\

#	cd /go/src/github.com/graph-uk/malibu-server/src/malibu-server &&\
#	go get golang.org/x/text/secure/bidirule &&\
#	npm run copy-client-worker-to-assets &&\
#	packr build

WORKDIR /go/src/github.com/graph-uk/malibu-server/src/malibu-server
ENTRYPOINT /go/src/github.com/graph-uk/malibu-server/src/malibu-server/malibu-server
EXPOSE 3133
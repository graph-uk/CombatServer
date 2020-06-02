FROM golang:1.13.7 as builder
WORKDIR /go
RUN PATH=$PATH:~/go/bin &&\
	apt-get update &&\
	curl -sL https://deb.nodesource.com/setup_10.x | bash -&&\
	apt-get install -y nodejs
COPY . /go
RUN npm install &&\
	chmod +x ./node_modules/packr-win-lin/packr &&\
	chmod +x ./node_modules/malibu-win-lin/malibu
RUN npm run iter

FROM golang:1.13.7 as prod
COPY --from=builder /go/src/malibu-server/assets/_/dist/malibu/malibu /bin/malibu
COPY --from=builder /go/int-tests/src/Tests/twoSessions/server/malibu-server /bin/malibu-server

WORKDIR /malibu-server
ENTRYPOINT malibu-server
EXPOSE 3133
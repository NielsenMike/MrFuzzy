#------------   Base   ------------#
FROM alpine:latest


#------------ GO BUILD ------------#
RUN apk update && apk add go gcc bash musl-dev 
RUN wget https://golang.org/dl/go1.17.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.17.linux-amd64.tar.gz
RUN cd /usr/local/go/src && ./make.bash
ENV PATH=$PATH:/usr/local/go/bin
RUN rm go1.17.linux-amd64.tar.gz
RUN apk del go
RUN go version
#------------ GO CLIENT ------------#
WORKDIR /app



COPY go.mod ./
RUN go mod download

#git clone#
COPY *.go ./ 				

RUN go build -o /app/mf_client

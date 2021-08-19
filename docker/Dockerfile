#------------   Base   ------------#
FROM alpine:latest


#------------ GO BUILD ------------#
RUN apk update && apk add go gcc bash musl-dev git
RUN wget https://golang.org/dl/go1.17.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.17.linux-amd64.tar.gz
RUN cd /usr/local/go/src && ./make.bash
ENV PATH=$PATH:/usr/local/go/bin
RUN rm go1.17.linux-amd64.tar.gz
RUN apk del go
RUN go version
#------------ GO CLIENT ------------#
RUN git clone https://gitlab.enterpriselab.ch/mnielsen/mf_client.git
WORKDIR /mf_client
RUN go mod download
RUN go build -o /mf_client

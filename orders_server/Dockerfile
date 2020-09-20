FROM golang:1.12.17

ENV GO111MODULE=on

# create project directory
RUN mkdir -p /go/src/github.com/avvarikrish/chefcurrygobbles/config

# copy project into container
COPY ./orders_server /go/src/github.com/avvarikrish/chefcurrygobbles/orders_server

# copy config into container
COPY ./config/orders_server.yml /go/src/github.com/avvarikrish/chefcurrygobbles/config

# copy pkg and proto
COPY ./pkg /go/src/github.com/avvarikrish/chefcurrygobbles/pkg
COPY ./proto /go/src/github.com/avvarikrish/chefcurrygobbles/proto

# Set the working directory
WORKDIR /go/src/github.com/avvarikrish/chefcurrygobbles/

# create the go modules file
RUN go mod init

# add any dependencies
RUN go mod tidy

# give executable permission to run command
RUN chmod +x orders_server/run.sh

EXPOSE 50051

CMD ["orders_server/run.sh"]

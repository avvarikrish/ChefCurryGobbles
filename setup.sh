# Generate the protobuf files
for PROTODIR in proto/*
do
    ./generate.sh $(basename $PROTODIR)
done

# Generate the client binaries
go build -o bin/ccgobbles ccgobbles_client/ccgobbles/ccgobbles_client.go 
go build -o bin/metrics ccgobbles_client/metrics/metrics_client.go
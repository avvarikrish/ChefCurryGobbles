for PROTODIR in proto/*
do
    ./generate.sh $(basename $PROTODIR)
done
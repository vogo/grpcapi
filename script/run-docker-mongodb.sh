
mkdir -p /tmp/mongodata

docker run -it --rm -d \
	--network grpc \
	-p 27017:27017 \
	-e MONGO_INITDB_ROOT_USERNAME=admin \
	-e MONGO_INITDB_ROOT_PASSWORD=grpcpass \
	-v /tmp/mongodata:/data/db \
	-h grpc-mongodb \
	--name grpc-mongodb \
	mongo


	

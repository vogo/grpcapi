
docker run -it --rm \
	--network=grpc \
	-p 9002:9002 \
	-h grpc-hello \
	--name grpc-hello \
	vogo/grpc-hello


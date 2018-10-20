
docker run -it --rm \
	--network=grpc \
	-p 9001:9001 \
	-h grpc-echo \
	--name grpc-echo \
	vogo/grpc-echo



docker run -it --rm \
	--network grpc \
	-p 8080:8080 \
	-h grpc-apigateway \
	--name grpc-apigateway \
	vogo/grpc-apigateway


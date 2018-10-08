package config

import "flag"

var (
	SecretKey = "ADkajfdJIALDSFJJkadf"

	APIGatewayAddress   = ":8080"
	TokenServiceAddress = ":9001"
	EchoServiceAddress  = ":9002"
	HelloServiceAddress = ":9003"

	EchoServiceEndpoint  = flag.String("echo_endpoint", EchoServiceAddress, "endpoint of echo service")
	HelloServiceEndpoint = flag.String("hello_endpoint", HelloServiceAddress, "endpoint of hello service")
)

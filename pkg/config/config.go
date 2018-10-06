
package config

import "flag"

var (
	APIGatewayAddress   = ":8080"
	EchoServiceAddress  = ":9001"
	HelloServiceAddress = ":9002"

	EchoServiceEndpoint  = flag.String("echo_endpoint", EchoServiceAddress, "endpoint of echo service")
	HelloServiceEndpoint = flag.String("hello_endpoint", HelloServiceAddress, "endpoint of hello service")
)

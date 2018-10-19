# Enhance grpc-ecosystem/grpc-gateway

[grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) is a grpc extension, it can generate a restful api proxy for grpc service.

see [code generate template](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/protoc-gen-grpc-gateway/gengateway/template.go)


grpc-gateway only provide a proxy ,but there is not a proxy intercepter.

I define a `google.protobuf.MethodOptions` named `allow_roles` in all methods, which is used for authorization.
Then it can get a map whoes key is the full path of method, and value is the defined allow roles.

I think the authorization can be done in grpc-gateway.

First, add method path into handler and provide it when registering.
```
type handler struct {
	pat     Pattern
	method  string    // <--------the full path of method /pkg.ServiceName/MethodName
	h       HandlerFunc
}
```

Second, add intercepter in ServeMux and provide method to initial it.

```
type ServeMux struct {
	// handlers maps HTTP method to a list of handlers.
	handlers               map[string][]handler
	forwardResponseOptions []func(context.Context, http.ResponseWriter, proto.Message) error
	marshalers             marshalerRegistry
	incomingHeaderMatcher  HeaderMatcherFunc
	outgoingHeaderMatcher  HeaderMatcherFunc
	metadataAnnotators     []func(context.Context, *http.Request) metadata.MD
	protoErrorHandler      ProtoErrorHandlerFunc
	intercepter            func(context.Context,http.ResponseWriter, *http.Request, method string) error  // <----- add intercepter definition
}
```

Third, call intercepter before calling handler:
```
for _, h := range s.handlers[r.Method] {
	pathParams, err := h.pat.Match(components, verb)
	if err != nil {
		continue
	}
	//-------------------------
	if s.intercepter != nil {
	  err = intercepter(ctx, w, r, h.method) //  <----------- call intercepter before calling handler
	  if err != nil {
		//TODO HANDLE ERROR
		return
	  }
	}
	//-------------------------

	h.h(w, r, pathParams)
	return
}
```


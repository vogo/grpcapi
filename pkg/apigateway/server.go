// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apigateway

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pborman/uuid"
	"github.com/vogo/grpcapi/pkg/apigateway/spec"
	"github.com/vogo/grpcapi/pkg/auth"
	"github.com/vogo/grpcapi/pkg/config"
	"github.com/vogo/grpcapi/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	//Authorization header
	Authorization = "Authorization"
	//RequestIDKey header
	RequestIDKey = "X-Request-Id"
)

type register struct {
	name     string
	f        func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	endpoint string
}

var (
	registers = []register{
		{"Token Manager", pb.RegisterTokenManagerHandlerFromEndpoint, config.TokenServiceAddress},
		{"Echo Service", pb.RegisterEchoServiceHandlerFromEndpoint, config.EchoServiceAddress},
		{"Hello Service", pb.RegisterHelloServiceHandlerFromEndpoint, config.HelloServiceAddress},
	}
	//ClientOptions grpc client options
	ClientOptions = []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}
)

func log() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New()
		c.Request.Header.Set(RequestIDKey, requestID)
		c.Writer.Header().Set(RequestIDKey, requestID)

		t := time.Now()

		// process request
		c.Next()

		latency := time.Since(t)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path

		logStr := fmt.Sprintf("%s | %3d | %v | %s | %s %s %s",
			requestID,
			statusCode,
			latency,
			clientIP, method,
			path,
			c.Errors.String(),
		)

		switch {
		case statusCode >= 400 && statusCode <= 499:
			glog.Info(logStr)
		case statusCode >= 500:
			glog.Error(logStr)
		default:
			glog.Info(logStr)
		}
	}
}

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//httprequest, _ := httputil.DumpRequest(c.Request, false)
				//logger.Critical(nil, "Panic recovered: %+v\n%s", err, string(httprequest))
				glog.Errorf("Panic recovered:%+v\n", err)
				c.JSON(500, gin.H{
					"title": "Error",
					"err":   err,
				})
			}
		}()
		c.Next() // execute all the handlers
	}
}
func serveGatewayMux(mux *runtime.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if req.URL.Path == "/api/v1/oauth2/token" {
			// skip auth requester
			mux.ServeHTTP(w, req)
			return
		}

		var stat *status.Status
		ctx := req.Context()
		_, outboundMarshaler := runtime.MarshalerForRequest(mux, req)

		token := strings.SplitN(req.Header.Get(Authorization), " ", 2)
		if token[0] != "Bearer" {
			stat = status.New(codes.Unauthenticated, "valid token required")
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, stat.Err())
			return
		}

		requester, err := auth.Validate(config.SecretKey, token[1])
		if err != nil {
			if err == auth.ErrExpired {
				stat = status.New(codes.Unauthenticated, "token expired")
			} else {
				stat = status.New(codes.Unauthenticated, "token invalid")
			}
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, stat.Err())
			return
		}
		req.Header.Set(auth.RequesterKey, requester.ToJSON())
		req.Header.Del(Authorization)

		mux.ServeHTTP(w, req)
	})
}

func swaggerHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(spec.Files["api.swagger.json"]))
	})
}

func mainHandler() http.Handler {
	var gwmux = runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			return metadata.Pairs(
				auth.RequesterKey, req.Header.Get(auth.RequesterKey),
				RequestIDKey, req.Header.Get(RequestIDKey),
			)
		}),
	)

	for _, r := range registers {
		glog.Infof("%v %v", r.name, r.endpoint)
		err := r.f(context.Background(), gwmux, r.endpoint, ClientOptions)
		if err != nil {
			glog.Fatal(err)
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/", serveGatewayMux(gwmux))
	return mux
}

func ginRun() error {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(log())
	r.Use(recovery())
	r.Any("/swagger/api/v1", gin.WrapH(swaggerHandler()))
	r.Any("/api/v1/*filepath", gin.WrapH(mainHandler()))

	return r.Run(config.APIGatewayAddress)
}

// Serve to start api gateway
func Serve() {
	flag.Parse()
	defer glog.Flush()

	if err := ginRun(); err != nil {
		glog.Fatal(err)
	}
}

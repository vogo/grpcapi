// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apigateway

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go2s/o2s/o2"
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
		if strings.HasPrefix(req.URL.Path, "/oauth2/") {
			// skip auth requester
			mux.ServeHTTP(w, req)
			return
		}
		ctx := req.Context()
		_, outboundMarshaler := runtime.MarshalerForRequest(mux, req)

		oauth2Svr := o2.GetOauth2Svr()
		tokenInfo, err := oauth2Svr.ValidationBearerToken(req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, status.New(codes.Unauthenticated, err.Error()).Err())
			return
		}
		req.Header.Set(auth.KeyUserID, tokenInfo.GetUserID())
		req.Header.Set(auth.KeyScope, tokenInfo.GetScope())
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
				auth.KeyUserID, req.Header.Get(auth.KeyUserID),
				auth.KeyScope, req.Header.Get(auth.KeyScope),
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

type server struct {
	cfg *config.Config
}

func (s *server) run() error {
	if !s.cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(log())
	r.Use(recovery())
	r.Any("/api/v1/*filepath", gin.WrapH(mainHandler()))
	r.GET("/swagger/api/v1", gin.WrapH(swaggerHandler()))

	s.initOauth2(r)

	return r.Run(config.APIGatewayAddress)
}

// Serve to start api gateway
func Serve(cfg *config.Config) {
	s := &server{cfg: cfg}
	if err := s.run(); err != nil {
		glog.Fatal(err)
	}
}

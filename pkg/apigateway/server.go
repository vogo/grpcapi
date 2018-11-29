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

	"github.com/vogo/grpcapi/pkg/apigateway/spec"
	"github.com/vogo/grpcapi/pkg/config"
	"github.com/vogo/grpcapi/pkg/constants"
	"github.com/vogo/grpcapi/pkg/identity"
	"github.com/vogo/grpcapi/pkg/pb"

	"github.com/gin-gonic/gin"
	"github.com/go2s/o2s/o2"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pborman/uuid"
	"github.com/vogo/clog"
	"github.com/vogo/grpcapi/pkg/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	//ClientOptions grpc client options
	ClientOptions = []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithUnaryInterceptor(authorizer),
	}
)

func log() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New()
		c.Request.Header.Set(constants.KeyRequestID, requestID)
		c.Writer.Header().Set(constants.KeyRequestID, requestID)

		t := time.Now()

		// process request
		c.Next()

		latency := time.Since(t)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path

		logStr := fmt.Sprintf("%3d | %v | %s | %s %s %s",
			statusCode,
			latency,
			clientIP, method,
			path,
			c.Errors.String(),
		)

		switch {
		case statusCode >= 400 && statusCode <= 499:
			clog.Log(clog.WarnLevel, requestID, logStr)
		case statusCode >= 500:
			clog.Log(clog.ErrorLevel, requestID, logStr)
		default:
			clog.Log(clog.InfoLevel, requestID, logStr)
		}
	}
}

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				clog.Fatal(c, "Panic recovered:%+v", err)
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
		token, ok := oauth2Svr.BearerAuth(req)
		if !ok {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, status.Error(codes.Unauthenticated, "token required"))
			return
		}
		claims, err := oauth2Svr.ParseJWTAccessToken(token)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, status.Error(codes.Unauthenticated, err.Error()))
			return
		}

		span := tracing.Span(req)
		defer span.Finish()

		userID := claims.Subject
		scopes := strings.Split(claims.Scope, ",")
		clog.Debug(ctx, "request user id %v", userID)

		//TODO --> currently use a struct as an identity object, CHANGE IT as business required
		identity := identity.New(userID, []pb.Role{pb.Role_USER}, scopes)
		req.Header.Set(constants.KeyIdentity, identity.String())
		req.Header.Del(constants.Authorization)

		mux.ServeHTTP(w, req)

	})
}

func swaggerHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(spec.Files["api.swagger.json"]))
	})
}

func mainHandler(cfg *config.Config) http.Handler {
	var gwmux = runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			md := metadata.Pairs(
				constants.KeyIdentity, req.Header.Get(constants.KeyIdentity),
				constants.KeyRequestID, req.Header.Get(constants.KeyRequestID),
			)
			for name, value := range req.Header {
				if tracing.SpanHTTPHeader(name) {
					md[strings.ToLower(name)] = value
				}
			}
			clog.Info(ctx, "grpc metedata:%v", md)
			return md
		}),
	)

	for _, r := range registers {
		clog.Info(nil, "proxy %v", r.endpoint)
		ctx := context.Background()
		err := r.f(ctx, gwmux, r.endpoint, ClientOptions)
		if err != nil {
			clog.Fatal(ctx, "%v", err)
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
	if clog.GlobalLevel() > clog.DebugLevel {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(log())
	r.Use(recovery())
	r.Any("/api/v1/*filepath", gin.WrapH(mainHandler(s.cfg)))
	r.GET("/swagger/api/v1", gin.WrapH(swaggerHandler()))

	s.initOauth2(r)

	tracing.Start()

	return r.Run(fmt.Sprintf(":%d", config.PortAPIGateway))
}

// Serve to start api gateway
func Serve(cfg *config.Config) {
	s := &server{cfg: cfg}
	if err := s.run(); err != nil {
		clog.Fatal(nil, "%v", err)
	}
}

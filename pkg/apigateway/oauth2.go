// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apigateway

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go2s/o2s/jwtex"
	"github.com/go2s/o2s/o2"
	"github.com/go2s/o2s/o2m"
	"github.com/vogo/clog"
	"gopkg.in/oauth2.v3"
)

func (s *server) initOauth2(r *gin.Engine) {
	clog.Info(nil, "init oauth2 server")

	database := s.cfg.Mongo.Database
	mgoClient := o2m.NewMongoClient(&s.cfg.Mongo)

	ts := o2m.NewTokenStore(mgoClient, database, "token")
	cs := o2m.NewClientStore(mgoClient, database, "client")
	us := o2m.NewUserStore(mgoClient, database, "user", o2m.DefaultMgoUserCfg())
	as := o2m.NewAuthStore(mgoClient, database, "auth")

	o2Cfg := o2.DefaultServerConfig()

	o2Cfg.ServerName = "Oauth2 Server"
	o2Cfg.JWTSupport = true
	o2Cfg.JWT = jwtex.JWTConfig{
		SignedKey:     []byte(s.cfg.SignKey),
		SigningMethod: jwt.SigningMethodHS512,
	}
	svr := o2.InitOauth2Server(cs, ts, us, as, o2Cfg, func(method, pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
		r.Handle(method, pattern, func(c *gin.Context) {
			handler(c.Writer, c.Request)
		})
	})

	//set allowed grant types
	svr.Config.AllowedGrantTypes = []oauth2.GrantType{
		oauth2.PasswordCredentials,
		oauth2.Refreshing,
	}

}

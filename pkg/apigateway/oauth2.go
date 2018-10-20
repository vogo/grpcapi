// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apigateway

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go2s/o2m"
	"github.com/go2s/o2s/o2"
	"github.com/vogo/clog"
	oauth2 "gopkg.in/oauth2.v3"
)

func (s *server) initOauth2(r *gin.Engine) {
	clog.Info(nil, "init oauth2 server")

	database := s.cfg.Mongo.Database
	mgoSession := o2m.NewMongoSession(&s.cfg.Mongo)

	ts := o2m.NewTokenStore(mgoSession, database, "token")
	cs := o2m.NewClientStore(mgoSession, database, "client")
	us := o2m.NewUserStore(mgoSession, database, "user", o2m.DefaultMgoUserCfg())
	as := o2m.NewAuthStore(mgoSession, database, "auth")

	o2Cfg := o2.DefaultServerConfig()

	o2Cfg.ServerName = "Oauth2 Server"
	o2Cfg.JWT = o2.JWTConfig{
		Support:    true,
		SignKey:    []byte(s.cfg.SignKey),
		SignMethod: jwt.SigningMethodHS512,
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

// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apigateway

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go2s/o2m"
	"github.com/go2s/o2s/o2"
	oauth2 "gopkg.in/oauth2.v3"
)

func (s *server) initOauth2(r *gin.Engine) {
	database := s.cfg.Mongo.Database
	mgoSession := o2m.NewMongoSession(&s.cfg.Mongo)

	ts := o2m.NewTokenStore(mgoSession, database, "token")
	cs := o2m.NewClientStore(mgoSession, database, "client")
	us := o2m.NewUserStore(mgoSession, database, "user", o2m.DefaultMgoUserCfg())
	as := o2m.NewAuthStore(mgoSession, database, "auth")

	cfg := o2.DefaultServerConfig()
	cfg.ServerName = "Oauth2 Server"

	svr := o2.InitOauth2Server(cs, ts, us, as, cfg, func(method, pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
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

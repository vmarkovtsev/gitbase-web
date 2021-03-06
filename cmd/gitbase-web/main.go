package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/src-d/gitbase-web/server"
	"github.com/src-d/gitbase-web/server/handler"
	"github.com/src-d/gitbase-web/server/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
)

// version will be replaced automatically by the CI build.
// See https://github.com/src-d/ci/blob/v1/Makefile.main#L56
var version = "dev"

// Note: maxAllowedPacket must be explicitly set for go-sql-driver/mysql v1.3.
// Otherwise gitbase will be asked for the max_allowed_packet column and the
// query will fail.
// The next release should make this parameter optional for us:
// https://github.com/go-sql-driver/mysql/pull/680
type appConfig struct {
	Env             string `envconfig:"ENV" default:"production"`
	Host            string `envconfig:"HOST" default:"0.0.0.0"`
	Port            int    `envconfig:"PORT" default:"8080"`
	ServerURL       string `envconfig:"SERVER_URL"`
	DBConn          string `envconfig:"DB_CONNECTION" default:"root@tcp(localhost:3306)/none?maxAllowedPacket=4194304"`
	SelectLimit     int    `envconfig:"SELECT_LIMIT" default:"100"`
	BblfshServerURL string `envconfig:"BBLFSH_SERVER_URL" default:"127.0.0.1:9432"`
	FooterHTML      string `envconfig:"FOOTER_HTML"`
}

func main() {
	// main configuration
	var conf appConfig
	envconfig.MustProcess("GITBASEPG", &conf)

	// logger
	logger := service.NewLogger(conf.Env)

	// database
	db, err := sql.Open("mysql", conf.DBConn)
	if err != nil {
		logger.Fatalf("error opening the database: %s", err)
	}
	defer db.Close()

	static := handler.NewStatic("build/public", conf.ServerURL, conf.SelectLimit, conf.FooterHTML)

	// start the router
	router := server.Router(logger, static, version, db, conf.BblfshServerURL)
	logger.Infof("listening on %s:%d", conf.Host, conf.Port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Host, conf.Port), router)
	logger.Fatal(err)
}

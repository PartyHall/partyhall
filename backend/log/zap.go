package log

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var LOG *zap.SugaredLogger
var DB *sqlx.DB

func Load(isInDev bool) {
	if isInDev {
		log, _ := zap.NewDevelopment()
		LOG = log.Sugar()
	} else {
		log, _ := zap.NewProduction()
		LOG = log.Sugar()
	}
}

package repo_wrappers

import (
	"database/sql"
	"minitwit/src/datalayer"
	"minitwit/src/models"
	"minitwit/src/utils"
)

var userRepo *datalayer.Repository[models.User]
var messageRepo *datalayer.Repository[models.Message]
var followerRepo *datalayer.Repository[models.Follower]
var latestProcessedRepo *datalayer.Repository[models.LatestProcessed]

func InitRepos(db *sql.DB) {
	db.SetMaxOpenConns(utils.GetEnvInt("DB_MAX_OPEN_CONS", 25))
	db.SetMaxIdleConns(utils.GetEnvInt("DB_MAX_IDLE_CONS", 10))
	db.SetConnMaxLifetime(utils.GetEnvDuration("DB_MAX_CONN_LIFETIME", "60m"))
	
	userRepo = datalayer.NewRepository[models.User](db, "users")
	messageRepo = datalayer.NewRepository[models.Message](db, "message")
	followerRepo = datalayer.NewRepository[models.Follower](db, "follower")
	latestProcessedRepo = datalayer.NewRepository[models.LatestProcessed](db, "latest_processed")
}

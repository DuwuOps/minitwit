package repo_wrappers

import (
	"database/sql"
	"minitwit/src/datalayer"
	"minitwit/src/models"
)

var userRepo *datalayer.Repository[models.User]
var messageRepo *datalayer.Repository[models.Message]
var followerRepo *datalayer.Repository[models.Follower]

func InitRepos(db *sql.DB) {
	userRepo = datalayer.NewRepository[models.User](db, "users")
	messageRepo = datalayer.NewRepository[models.Message](db, "message")
	followerRepo = datalayer.NewRepository[models.Follower](db, "follower")
}

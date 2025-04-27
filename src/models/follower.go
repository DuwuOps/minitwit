package models

type Follower struct {
	FollowerID  int `db:"follower_id"`  
	FollowingID int `db:"following_id"` 
}

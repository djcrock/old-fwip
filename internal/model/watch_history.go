package model

type WatchHistory struct {
	UserId      int64 `json:"user_id"`
	TitleId     int64 `json:"title_id"`
	Watched     bool  `json:"watched"`
	WantToWatch int64 `json:"want_to_watch"`
}

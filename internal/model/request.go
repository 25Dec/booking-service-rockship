package model

type PaginationParams struct {
	Limit  int `form:"limit,default=10"`
	Offset int `form:"offset,default=0"`
}

// form to is timestamp without timezone
type FromToParams struct {
	From string `form:"from" time_format:"2006-01-02 15:04:05" default:"2019-01-01 00:00:00+07"`
	To   string `form:"to" time_format:"2006-01-02 15:04:05" default:"2222-01-01 00:00:00+07"`
}

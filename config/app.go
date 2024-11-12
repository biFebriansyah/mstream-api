package config

type ResultWarp struct {
	Data any
	Meta any
}

type Meta struct {
	Next  int32 `json:"next"`
	Prev  int32 `json:"prev"`
	Total int32 `json:"total"`
}

type FilterParams struct {
	Param  string
	Column string
}

type Pagination struct {
	Name  string
	Page  int32
	Limit int32
}

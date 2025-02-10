package infrastructure

import "gorm.io/gorm"

type Filter interface {
	GenerateMods(db *gorm.DB) *gorm.DB
}

type OrderBy struct {
	Column string
	Order  string
}

func SetOrderBy(s string) OrderBy {
	var o OrderBy
	if s == "" {
		return o
	}
	if s[:1] == "-" {
		o.Order = "DESC"
		o.Column = s[1:]
		return o
	}
	o.Order = "ASC"
	o.Column = s
	return o
}

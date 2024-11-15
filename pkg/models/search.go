package models

type Search struct {
	Text string `form:"search"`
}

type Filter struct {
	Text string `form:"filter"`
}

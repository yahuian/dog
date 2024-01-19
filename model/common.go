package model

type ListOption struct {
	Page  int
	Size  int
	Query any
	Args  []any
	Order string
}

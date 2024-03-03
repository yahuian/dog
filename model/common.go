package model

type ListOption struct {
	Page  int
	Size  int
	Query any
	Args  []any
	Order string
}

type AllOption struct {
	Query any
	Args  []any
	Order string
}

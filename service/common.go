package service

type Paging struct {
	Page int `form:"page,default=1" validate:"min=1"`          // 页码
	Size int `form:"size,default=10" validate:"min=1,max=100"` // 每页数量
}

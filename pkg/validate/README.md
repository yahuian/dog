## 开发背景

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type user struct {
	Name string `binding:"required"`
	Age  int    `binding:"gt=18"`
}

func main() {
	r := gin.Default()

	r.POST("/ping", func(c *gin.Context) {
		var u user
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": u})
	})

	r.Run()
}
```

```
curl -X POST http://localhost:8080/ping -H 'content-type: application/json' -d '{ "name": "tom" }'

{"msg":"Key: 'user.Age' Error:Field validation for 'Age' failed on the 'gt' tag"}
```

可以看到 gin 默认的 validate 报错信息非常不利于调用者排查错误

## 改造效果

```
{
    "code": 400,
    "msg": "Name is a required field, Age must be greater than 18",
    "data": null
}
```

需要注意标签名 `binding` 修改为了 `validate`

```go
type UpdateIn struct {
	ID uint `validate:"required"`
}
```

此外增加了 `validate.Struct()` 和 `validate.Var()` 方法，可以对结构体和变量进行校验

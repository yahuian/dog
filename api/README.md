## 开发规范

### url 路径中增加统一的 api 前缀和版本号

```
https://www.example.com/api/v1
```

### 简单的 CRUD 采用 RESTful 风格

```
查询 GET /user
新增 POST /user
编辑 PUT /user
删除 DELETE /user
```

### GET 请求参数放入 url 中，其余方法均采用 json 格式放入 body 中

```
GET /user?Name=小明&Age=18

POST /user
{
  "Name":"小明",
  "Age":18
}

DELETE /user
{
  "ID": [1,2,3]
}
```
特殊情况：如果查询参数太多超出 url 最大长度，此时可以放入 body 中

### path 不可动态变化，需要参数则增加 url 查询参数

```
❌ GET /user/1
✅ GET /user?ID=1
```
因为动态变化的 path 对监控系统不友好，类似微信小程序的 page 页面不允许变化


### 复杂一点的场景在 path 中增加 action 辅助说明，且全部使用 POST 请求

```
用户登录 POST /user/login
停止任务 POST /task/stop
问题置顶 POST /task/top
```

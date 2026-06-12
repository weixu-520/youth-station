## gateway.go

### 命令行参数解析

```go
 var configFile = flag.String("f", "etc/gateway-api.yaml", "the config file")

func main() {
    flag.Parse()
}
```

flag.String 定义了一个名为 f 的命令行标志，默认值为 "etc/gateway-api.yaml"，帮助信息为 "the config file"，返回字符串指针。

flag.Parse() 解析 os.Args 中的命令行参数，将用户输入的 -f xxx.yaml 赋值给 *configFile。

作用：让配置文件路径可被外部灵活指定，避免硬编码

### 加载配置文件

```go
var c config.Config
conf.MustLoad(*configFile, &c)
```

conf.MustLoad（通常来自 github.com/zeromicro/go-zero/core/conf）：

读取 YAML/JSON 配置文件（路径为 *configFile），反序列化到结构体 c 中。

Must 前缀表示：如果加载失败（如文件不存在、格式错误），会直接 panic 终止程序。

作用：将配置（如服务端口、数据库连接、Redis 地址等）载入内存，供后续使用。

### 创建REST服务器

```go
server := rest.MustNewServer(c.RestConf)
defer server.Stop()
```

rest.MustNewServer（来自 github.com/zeromicro/go-zero/rest）：

根据 c.RestConf 中的配置（如监听地址、超时、TLS 等）创建一个 HTTP 服务器对象。

同样 Must 表示创建失败时 panic。

defer server.Stop()：确保在 main 函数退出前优雅关闭服务器（释放资源、完成正在处理的请求）。

### 初始化服务上下文

```go
ctx := svc.NewServiceContext(c)
```

svc.NewServiceContext 是业务层自定义的构造函数（一般在 internal/svc/servicecontext.go 中定义）。

作用：根据配置 c 初始化服务所需的依赖，例如：

数据库连接（gorm、sqlx）

Redis 客户端

其他 RPC 客户端

日志、缓存、消息队列等
并将它们封装在 ctx 结构体中，供后续业务逻辑使用。

### 注册路由处理器

```go
handler.RegisterHandlers(server, ctx)
```

handler.RegisterHandlers 是路由注册函数（一般在 internal/handler 包中定义）。

作用：将 API 路径（如 /api/user/info）与对应的处理函数（Handler）绑定到 server 上，同时把依赖 ctx 传递给这些 Handler。

常见写法：handler.RegisterHandlers(server, ctx) 内部会调用 server.AddRoute(...) 或 server.POST(...)。

## Handler.go(以adminGetApplicationListHandler为例)

### http.HandlerFunc

`AdminGetApplicationListHandler` 是一个**工厂函数**，它接收 `*svc.ServiceContext`，返回一个 `http.HandlerFunc`。

- `http.HandlerFunc` 是 Go 标准库中 `net/http` 包的一个类型：`type HandlerFunc func(ResponseWriter, *Request)`，它实现了 `http.Handler` 接口。

- 这意味着返回的函数可以直接用于 `http.HandleFunc` 或 go-zero 的路由注册。

**设计目的**：通过闭包捕获 `svcCtx`，使得真正处理请求的匿名函数在运行时能够访问到应用级别的依赖（数据库、缓存、配置等）。

### 解析请求参数 (`httpx.Parse`)

```go
if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
```

### `httpx.Parse(r, &req)` 的详细行为

该函数实现了智能参数绑定，逻辑如下：

- 首先检查 `r.Header.Get("Content-Type")`：
  
  - 如果是 `application/json` → 使用 `json.NewDecoder(r.Body).Decode(&req)`
  
  - 如果是 `application/x-www-form-urlencoded` → 调用 `r.ParseForm()` 然后绑定 form 值到结构体
  
  - 如果是 `multipart/form-data` → 同样绑定 form（但不处理文件）
  
  - 如果 `Content-Type` 不匹配，会尝试从 `r.URL.Query()` 解析（即 query string）

- 无论哪种情况，最终都会将**query string 中的参数**合并到结构体（优先级低于 body/form）。

- 支持字段标签中的 `optional`（可选）和 `default`（默认值），如果未提供则使用默认值。

- 若解析失败（如 required 字段缺失、类型转换错误），返回非 nil 错误。

### 错误响应 `httpx.ErrorCtx`

- 它自动判断错误类型：
  
  - 如果错误是 `*errors.CodeMsg`（go-zero 自定义错误），则提取业务码和消息。
  
  - 否则，默认使用 HTTP 状态码 400（Bad Request）或 500（Internal Server Error）。

- 构造 JSON 响应体，格式固定为：`{"code": <整数码>, "msg": "<错误信息>"}`。

- 通过 `r.Context()` 传递上下文，以便在中间件或日志中关联请求追踪 ID（trace-id）。

- 最终调用 `w.WriteHeader()` 和 `json.NewEncoder(w).Encode(errBody)`。

这一行之后的 `return` 确保匿名函数立即退出，不再往下执行。

## 创建 Logic 对象

```go
l := admin.NewAdminGetApplicationListLogic(r.Context(), svcCtx)
```

### `NewAdminGetApplicationListLogic` 构造函数

典型实现如下：

```go
type AdminGetApplicationListLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewAdminGetApplicationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetApplicationListLogic {
    return &AdminGetApplicationListLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}
```

它只是简单地将 `context` 和 `svcCtx` 保存到结构体中，并不执行任何业务。

### 为什么要把 ctx 和 svcCtx 分开传？

ctx 代表本次请求的上下文（包含超时、取消信号、trace 信息），每个请求不同。

svcCtx 是全局应用上下文（包含数据库连接池、redis 客户端、配置），所有请求共享。



AdminGetApplicationListHandler(svcCtx)` 返回的 `http.HandlerFunc` 最终会被注册到路由中，如：

```go
server.AddRoute(rest.Route{
    Method: http.MethodGet,
    Path:   "/admin/application/list",
    Handler: AdminGetApplicationListHandler(svcCtx),
})
```

此时 `svcCtx` 被闭包捕获，路由处理时不会再改变。

goctl api go -api api/gateway.api -dir apps/gateway -style gozero

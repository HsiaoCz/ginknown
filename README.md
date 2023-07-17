# ginknown

### 1、mysql 预处理

普通 sql 语句执行过程

- 客户端对 sql 语句进行占位符替换得到完整的 sql 语句
- 客户端发送完整 sql 语句到 mysql 服务端
- mysql 服务端执行完整的 sql 语句并将结果返回给客户端

预处理的执行过程:

- 把 sql 语句分成两部分，命令部分与数据部分
- 先把命令部分发送给 mysql 服务端,mysql 服务端进行 sql 预处理
- 然后把数据部分发送给 mysql 服务端,mysql 服务端对 sql 语句进行占位符替换
- mysql 服务端执行完整的 sql 语句并将结果返回给客户端

什么时候需要使用占位符

当批量的执行一条 mysql 语句，除了数据不同，别的都相同的时候

sql 注入的问题

当我们直接使用用户输入的内容来执行 sql 语句的时候，就容易发送 sql 注入
一个原则就是，不要自己拼接 sql 语句
还有一个原则，不要相信用户输入的内容

### 2、绑定变量(bindvars)

查询占位符，在内部称为 bindvas，应该始终使用它向数据库发送值，因为它们可以防止 sql 注入攻击，不过，不要用来占位表名

### 3、pipeline

pipeline 主要是一种网络优化，它本质上一位着客户端缓冲一堆命令并一次性将他们发送到服务器。RTT(往返的时延)
这些命令不能保证在事务中执行。
这样做的好处是节省了每个命令的网络返回时间

```go
pipe:=rdb.Pipeline()

incr:=pipe.Incr("pipeline_counter")
pipe.Expire("pipeline_counter",time.Hour)

_,err:=pipe.Exec()
fmt.Println(incr.Val(),err)
```

pipeline 可以将三个命令一起发送，RTT 只有一个

**redis 事务**

Redis 时单线程的，因此单个命令始终是原子的，但是来自不同客户端的两个给定命令可以依次执行
multi/exec 能够确保在 multi/exec 两个语句之间的命令之间没有其他客户端正在执行命令

在这种场景下我们需要使用 TxPipeline,TxPipeline 总体上类似于 pipeline,但是它内部会使用 multi/exec 包裹排队的命令

```go
pipe:=rdb.TxPipeline()

incr:=pipe.Incr("tx_pipeline_counter")
pipe.Expire("tx_pipeline_counter",time.Hour)

_,err:=pipe.Exec()
fmt.Println(incr.Val(),err)

```

**Wathc**

某些场景下，我们除了要使用 Multi/Exec 命令外，还需要配合 Watch 命令
在用户使用 Watch 命令监视某个键之后，直到该用户执行 exec 命令的这段时间里，如果有其他用户抢先对被监视的键进行了替换\更新\删除等操作，那么用户尝试执行 exec 的时候，事务将失败并返回一个错误，用户可以根据这个错误选择重试事务或者放弃事务

```go
// watch watch_count的值，并在值不变的前提下将其值+1
key:="watch_count"
err=client.Watch(func(tx *redis.Tx)error{
    n,err:=tx.Get(key).Int()
    if err!=nil&&err!=redis.Nil{
        return err
    }
_,err=tx.Pipeline(func(pipe redis.Pipeline)error{
    pipe.Set(key,n+1,0)
    return nil
})
return err
},key)

```

### 4、zap 日志库

一个好的日志记录器能够:

- 能够将事件记录到文件，而不是应用程序控制台
- 日志切割，能够根据文件大小、时间或间隔等来切割日志文件
- 支持不同的日志级别，例如 INFO，DEBUG、ERROR 等
- 能够打印基本信息，如调用文件/函数名和行号,日志时间等

### 5、定义错误码

```go
/*

{
    code:1001
    msg:请求成功
    data:{}
}
*/

// 对于响应，我们可以定义一个结构体
type ResponseData struct{
    Code ResCode  `json:"code"`
    Msg  interface{} `json:"msg"`
    Data interface{}   `json:"data"`
}

func ResponseErr(c *gin.Context,code ResCode){
   responseData:=&ResponseData{
    Code:code,
    Msg:code.Msg(),
    Data:nil,
   }
   c.JSON(http.StatusOk,responseData)
}

func ResponseErrorWithMsg(c *gin.Context,code ResCode,msg interface{}){
    c.JSON(http.StatusOK,&ResponseData{
        Code:code,
        Msg:msg,
        Data:nil,
    })
}

func ResponseSuccess(c *gin.Context,data interface{}){
    responseData:=&ResponseData{
        Code:CodeSuccess,
        Msg:CodeSuccess.Msg(),
        Data:data,
    }
    c.JSON(http.StatusOK,responseData)
}

// 定义错误码
type ResCode int

const (
   CodeSuccess=1000+iota
   CodeInvalidParam
   CodeUserExist
   CodeUserNotExist
   CodeInvalidPassword
   CodeServerBusy
)

var codeMsgMap=map[ResCode]string{
    CodeSuccess:"success",
    CodeInvalidParam:"请求参数错误",
    CodeUserExist:"用户名已存在",
    CodeUserNotExist:"用户名不存在",
    CodeInvalidPassword:"用户名或密码错误",
    CodeServerBusy:"服务繁忙",
}

func (c ResCode)Msg()string{
    msg,ok:=codeMsgMap[c]
    if !ok{
        msg=codeMsgMap[CodeServerBusy]
    }
    return msg
}

```

### 6、用户认证模式

HTTP 是一个无状态的协议，一次请求结束后，下次发送服务器就不知道这个请求是谁发来的了。

Cookie-Session 模式

- 客户端使用用户名、密码进行认证
- 服务端验证用户名、密码正确后生成并存储 Session，将 SessionID 通过 Cookie 返回给客户端
- 客户端访问需要认证的的接口时在 cookie 中携带 sessionID
- 服务端通过 SessionID 查找 Session 并进行鉴权，返回给客户端需要的数据

Session 和 Cookie 中存在多种问题

可以使用 Token,无状态的鉴权

### 7、限制同一时间同一用户只能登录一台设备

解决这个问题，在生成 token 的时候，拿到用户的 id,将这个 id 与 token 的对应关系存储到 redis 里面，
后续用户登录的时候，除了验证 token 是否有效，还可以通过用户的 ID 与 redis 里面的 token 时否对应，如果不一致，重新登录

### 8、AIR 实现实时热重载

### 9、分页展示

### 10、解决传递给前端数字 ID 数据失真的问题

RESTful API 数据通过 JSON 格式的数据
前端 js number 能表示的数字的范围是-(2^53-1)到(2^53-1)之间
但是后端 go 的 int64 能表示的数字的范围是-(2^63-1)到(2^63-1)

解决办法：前端传递的数据转成字符串传递给后端，后端传递给前端的数据在序列化的时候转成字符串

```go
type Person struct{
    ID int64 `json:"id,string"`
    Username string `json:"username"`
}
```

在 json 的 tag 里面加给,string 就可以解决这个问题

### 11、使用 Swagger 生成接口文档

模型后面写注释可以在 swagger 文档里面显示这个注释
在 tag 里面写 explame 可以在生成文档的时候显示出来
有一个点，每个文档返回的数据不同，可以为接口定义一个专门的模型

例如:

```go
type _ResponsePostList struct{

}
```

### 12、为接口编写单元测试

比如有这样一段代码

```go

package main

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
```

测试:

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
```

在 gin 对接口进行测试

```go
func TestCreatePostHandler(t *testing.T){
    gin.SetMode(gin.TestMode)
    r:= gin.Default()
    url:="/api/v1/post"
    r.Post(url,CreatePostHandler)

    // 这个接口有一些依赖数据
    // 我们自己造一些
    body:=`{
       "commuit_id":1,
       "title":"test",
       "content":"just a test"
    }`

    req,_:=http.NewRequest(http.MethodPost,url,bytes.NewReader([]byte(body)))
    w:=httptest.NewRecorder()
    r.ServeHTTP(w,req)

    assert.Equal(t,200,w.Code)
    // 判断响应的内容是不是按照预期返回了需要登录的错误
    // 方法1：判断响应内容中是不是包含指定的字符串
    assert.Contains(t,w.Body.String(),"需要登录")

    // 方法二
    res:=new(ResponseData)
    if err:=json.Unmarshal(w.Body.Bytes(),res);err!=nil{
        t.Fatal("json.Unmarshal w.Body failed,err:%v\n",err)
    }
    assert.Equal(t,res.Code,CodeNeedLogin)
}
```

单元测试有一点需要注意:

我们测试一个需要操作数据库的接口的时候，直接在测试函数里面执行这个函数是跑不通的

因为单元测试只会执行这个函数，而数据库操作依赖一个 DB,他会报空指针引用
解决办法：

```go
// 在test文件中创建一个init函数

func init(){
    // 这里填上我们需要的mysql配置文件信息
    mysqlcfg:=&conf.MysqlConfig{

    }
    // 这里执行数据库初始化操作
    // 初始化db
    err:=Init(mysqlcfg)
    if err!=nil{
        panic(err)
    }

}
```

### 13、常用的 HTTP 压力测试

压力测试相关术语

- 响应时间(RT):指系统对请求做出响应的时间
- 吞吐量:指系统在单位时间内处理的请求的数量
- qps：每秒查询率，是一台服务器每秒能够响应的查询次数，是对一个特点的服务器在规定时间内所处理流量多少的衡量标准
- TPS:每秒钟系统能够处理的交易或事务的数量
- 并发连接数：某个时刻服务器能接收的请求总数

压力测试工具:

ab
wrk

### 14、限流策略

漏桶和令牌桶

漏桶按照固定的速率去处理，有点像削峰填谷
但是它并不能很好的处理有大量突发请求的场景
毕竟在某些情况下我们可能需要提高系统的处理效率，而不是一味的按照固定速率处理请求

令牌桶

不断的往桶里放令牌，生成令牌的速度是恒定的，请求在桶里拿到令牌后可以被服务端处理，这种限流策略的好处是，有时候突然有大量请求时也能处理，当然，桶是有容量限制的，大量请求的时候，没拿到令牌的请求就需要等待

### 15、使用 pprof 对 go 进行调优

### 16、使用 docker 进行部署

```dockerfile
# 基于的基础镜像
FROM golang:alpine

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录: /build

WORKDIR /build

# 将代码复制到容器中

COPY . .

# 将我们的代码编译成二进制可执行文件app

RUN go build -o app .

# 移动到用于存放生成的二进制文件的/dist 目录
WORKDIR /dist

# 将二进制文件从/build 目录复制到这里

RUN cp /build/app .

# 声明服务端口
EXPOSE 8888

# 启动容器时运行的命令
CMD ["/dist/app"]

```

docker build . -t goweb_app
-t 指定名字

docker run -p 8888:8888 goweb_app

-p 指的是，容器里面的端口映射到系统的端口 8888

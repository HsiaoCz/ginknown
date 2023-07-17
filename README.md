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

HTTP是一个无状态的协议，一次请求结束后，下次发送服务器就不知道这个请求是谁发来的了。

Cookie-Session模式

- 客户端使用用户名、密码进行认证
- 服务端验证用户名、密码正确后生成并存储Session，将SessionID通过Cookie返回给客户端
- 客户端访问需要认证的的接口时在cookie中携带sessionID
- 服务端通过SessionID查找Session并进行鉴权，返回给客户端需要的数据

Session和Cookie中存在多种问题

可以使用Token,无状态的鉴权

### 7、限制同一时间同一用户只能登录一台设备

解决这个问题，在生成token的时候，拿到用户的id,将这个id与token的对应关系存储到redis里面，
后续用户登录的时候，除了验证token是否有效，还可以通过用户的ID与redis里面的token时否对应，如果不一致，重新登录


### 8、AIR实现实时热重载

### 9、分页展示

### 10、解决传递给前端数字ID数据失真的问题

RESTful API 数据通过JSON格式的数据
前端js number 能表示的数字的范围是-(2^53-1)到(2^53-1)之间
但是后端go的int64能表示的数字的范围是-(2^63-1)到(2^63-1)
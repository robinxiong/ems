#Responder
Responder根据request中accept mime type的不同，提供不同的响应内容，可以查看auth/handlers.go的DefaultLoginHandler

##Usage
### Register mime type
```go
import "github.com/qor/responder"

responder.Register("text/html", "html")
responder.Register("application/json", "json")
responder.Register("application/xml", "xml")
```
response默认注册了以上的三种类型，你可以通过Register函数注册更多的函数, 它接收两个参数

1. The mime type, like `text/html`
2. The format of the mime type, like `html`

### 响应注册的mime types

```go
func handler(writer http.ResponseWriter, request *http.Request) {
  responder.With("html", func() {
    writer.Write([]byte("this is a html request"))
  }).With([]string{"json", "xml"}, func() {
    writer.Write([]byte("this is a json or xml request"))
  }).Respond(request)
})
```

如果没有找到相应的mime type, 第一个html将作为默认的响应


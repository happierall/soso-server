# Soso server (Go)
### Comfortable, fast, bidirectional protocol over websocket instead REST

## Client lib
[soso-client](https://github.com/happierall/soso-client)

##Install
```
  go get -u github.com/happierall/soso-server
```

##Usage
```go
  import (
  	"fmt"

  	soso "github.com/happierall/soso-server"
  )

  func main() {

    // Simple Use:
    Router := soso.Default()

    Router.CREATE("message", func (m *soso.Msg) {
      fmt.Println(m.RequestMap)

      m.Success(map[string]interface{}{
        "id": 1,
      })
    })

    Router.Run(4000)

  }
```

```go
  // Add routes as list:
  var Routes = soso.Routes{}
  Routes.Add("message", "create", ChatSendMessage)

  // Handler:
  func ChatSendMessage(m *soso.Msg) {

    // Send direct message:
    soso.SendMsg("notify", "created", m.Session,
      map[string]interface{}{
        "text": "Congratulation for first message",
      },
    )

    m.Success(map[string]interface{}{
      "message": "message hi",
      "id": m.RequestMap["id"],
    })

  }

  Router := soso.Default()
  Router.HandleRoutes(Routes)
  Router.Run(4000)
```


```go
// Custom listener:
Router := soso.Default()
Router.Handle("message", "create", func (m *soso.Msg) {
  fmt.Println(m.RequestMap)
})
http.HandleFunc("/soso", Router.receiver)
http.ListenAndServe("localhost:4000", nil)
```

## Client request (if use without [soso-client](https://github.com/happierall/soso-client))
```javascript
  // javascript pure:
  var sock = new WebSocket("ws://localhost:4000/soso")

  var data = {
        data_type: "message",
        action_str: "create",
        log_map: {},
        request_map: {msg: "hello world"},
        trans_map: {},
  }

  sock.onopen = () => {

    sock.send( JSON.stringify( data ) )

  }
```

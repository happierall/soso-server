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
  Routes.Add("user", "create", UserCreate)

  // Handler:
  func UserCreate(m *soso.Msg) {

    type Data struct {
      ID int64 `json:"id"`
    }
    data := &Data{}
    m.ReadData(data)

    // Send direct message:
    soso.SendMsg("notify", "created", m.Session,
      map[string]interface{}{
        "text": "Congratulation for first message",
      },
    )

    m.Success(map[string]interface{}{
      "id": data.ID,
    })

  }

  Router := soso.Default()
  Router.HandleRoutes(Routes)
  Router.Run(4000)
```


```go
// Custom listener:
Router := soso.Default()
Router.Handle("user", "create", func (m *soso.Msg) {})
http.HandleFunc("/soso", Router.Receiver)
http.ListenAndServe("localhost:4000", nil)
```

## Client request (if use without [soso-client](https://github.com/happierall/soso-client))
```javascript
  // javascript pure:
  var sock = new WebSocket("ws://localhost:4000/soso")

  var data = {
      model: "message",
      action: "create",
      data: {msg: "hello world"},
      log: {},
      other: {},
  }

  sock.onopen = () => {

    sock.send( JSON.stringify( data ) )

  }
```

# Golang WebSocket Framework (Warning! Not maintained)
### Soso-server

### Additional libs
[soso-client](https://github.com/happierall/soso-client)
[soso-auth](https://github.com/happierall/soso-auth)

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
  Routes.CREATE("user", UserCreate)

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

  // You can handle net/http handlers
  http.HandleFunc("/oauth/callback/github", callbackGithub)

  // And run
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

## Middleware and Auth
```go

import (
  soso "github.com/happierall/soso-server"
  jose "github.com/dvsekhvalnov/jose2go"
)

func main() {
  Router := soso.Default()

  Router.Middleware.Before(func (m *soso.Msg, start time.Time) {
  	token, uid, err := readToken(m)

    // User is blank, you can use it
  	m.User.ID = strconv.FormatInt(uid, 10)
  	m.User.Token = token
  	m.User.IsAuth = true
  	m.User.IsAnonymous = true
  })

  Router.SEARCH("user", func(m *soso.Msg) {
    if m.User.IsAuth {
      fmt.Println(m.User.Token, m.User.ID)
    }
  })

  Router.Run(4000)
}
```


## Events and Sessions

```go
  func main() {
    Router := soso.Default()

    soso.Sessions.OnOpen(func(session soso.Session) {
      fmt.Println("Client connected")
    })

    soso.Sessions.OnClose(func(session soso.Session) {
      fmt.Println("Client disconnected")
    })

    Router.Middleware.Before(func (m *soso.Msg, start time.Time) {
      
      uid := AuthUser()

      if m.User.IsAuth {
        soso.Sessions.Push(m.Session, uid)) 
      }

    })


    Router.Run(4000)
  }
```

### License
[MIT](http://opensource.org/licenses/MIT)

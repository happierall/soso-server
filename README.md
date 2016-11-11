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
  })

  Router.SEARCH("user", func(m *soso.Msg) {
    if m.User.IsAuth {
      fmt.Println(m.User.Token, m.User.ID)
    }
  })

  Router.Run(4000)
}


// Example of readToken. It's may be absolutely different ;)
func readToken(m *soso.Msg) (string, int64, error) {
	type Other struct {
		Token string `json:"token"` // Recommend to use "token" name
  }
	other := Other{}
	err := m.ReadOther(&other)

	payload, _, err := jose.Decode(other.Token, []byte("Secret_key"))

	type tokenData struct {
		UID int64
	}
	var td tokenData
	json.Unmarshal([]byte(payload), &td)

	return other.Token, td.UID, nil
}

```

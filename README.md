# Soso server (Go)
## Comfortable, fast, bidirectional protocol over websocket instead REST

### Client lib
[soso-client](https://github.com/happierall/soso-client)

###Install
```
  go get -u github.com/happierall/soso-server
```

###Usage
```go

//Simple Use:
  Router := soso.Default()
  Router.CREATE("message", ChatSendMessage)
  Router.Run(4000)


//Add routes as list:
  var Routes = soso.Routes{}
  Routes.Add("create", "message", ChatSendMessage)

  Router := soso.Default()
  Router.HandleRoutes = Routes
  Router.Run(4000)


//Custom listener:
  Router := soso.Default()
  Router.Handle("message", "create", ChatSendMessage)
  http.HandleFunc("/soso", Router.receiver)
  http.ListenAndServe("localhost:4000", nil)


//Handler:
  func ChatSendMessage(m *soso.Msg) {

    m.Success(map[string]interface{}{
      "message": "message hi",
      "id": m.RequestMap["id"],
    })

  }


//Send direct message:
  soso.SendMsg("message", "created", session,
    map[string]interface{}{
      "id": "1",
    },
  )

```

### Client request (if use without soso-client)
#### Reccomend lib:
[soso-client](https://github.com/happierall/soso-client)

```javascript
  // javascript pure:
  var sock = new WebSocket("ws://localhost:4000/soso")

  var data = {
        data_type: "message",
        action_str: "create",
        log_map: {},
        request_map: {msg: "hello world"},
        trans_map: {}
  }

  sock.onopen = () => {

    sock.send( JSON.stringify( data ) )

  }
```

package views

func init() {
  Routes.Add("users", "retrieve", UserRetrieve)
  Routes.Add("users", "create", UserCreate)
}

func UserRetrieve(m *soso.Msg) {

  m.Success(map[string]interface{}{
    "id": m.RequestMap["id"],
    "Name": "Mary"
  })

}

func UserCreate(m *soso.Msg) {

  m.Success(map[string]interface{}{
    "message": "message hi",
    "id": m.RequestMap["id"],
  })

}

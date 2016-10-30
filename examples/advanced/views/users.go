package views

import soso "github.com/happierall/soso-server"

func init() {
	Routes.Add("users", "retrieve", UserRetrieve)
	Routes.Add("users", "create", UserCreate)
}

func UserRetrieve(m *soso.Msg) {

	type request struct {
		ID int64 `json:"id"`
	}
	req := &request{}
	m.ReadRequest(req)

	m.Success(map[string]interface{}{
		"id":   req.ID,
		"Name": "Mary",
	})

}

func UserCreate(m *soso.Msg) {

	m.Success(map[string]interface{}{
		"message": "message hi",
	})

}

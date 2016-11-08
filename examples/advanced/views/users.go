package views

import soso "github.com/happierall/soso-server"

func init() {
	Routes.Add("user", "get", UserGet)
	Routes.Add("user", "create", UserCreate)
}

func UserGet(m *soso.Msg) {

	type Data struct {
		ID int64 `json:"id"`
	}
	data := &Data{}
	m.ReadData(data)

	m.Success(map[string]interface{}{
		"id":   data.ID,
		"Name": "Mary",
	})

}

func UserCreate(m *soso.Msg) {

	m.Success(map[string]interface{}{
		"message": "message hi",
	})

}

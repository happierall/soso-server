package views

import soso "../../.."

func init() {
	Routes.GET("user", userGet)
	Routes.CREATE("user", userCreate)
}

func userGet(m *soso.Msg) {

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

func userCreate(m *soso.Msg) {

	m.Success(map[string]interface{}{
		"message": "message hi",
	})

}

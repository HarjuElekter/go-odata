package odata

import "testing"

func Test_Client(t *testing.T) {
	c := NewClient()
	c.SetBaseURL("http://10.110.51.102:13548/BC210_WS/ODataV4/Company('HARJU ELEKTER')/")
	c.SetBaseCredentials(&BaseAuthorization{
		Name:     "app",
		Password: "YxHEEnJkaip/eby/pMnHC9zhD5iuQ7Q+WiawjgUClEQ=",
	})

	body, _, err := c.Get("BI_ItemCategories")
	if err != nil {
		t.Error(err)
	}

	t.Log(string(body))
}

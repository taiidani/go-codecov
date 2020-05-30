package codecov

import "net/http"

func (c *Client) addAuthorization(request *http.Request) {
	request.Header.Add("Authorization", "token "+c.token)
}

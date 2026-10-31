package ngx

import (
	"context"
	"encoding/json"
	"net/http"
)

type BodyDecoder func(resp *http.Response) error

func JSONDecoder(v interface{}) BodyDecoder {
	return func(resp *http.Response) error {
		return json.NewDecoder(resp.Body).Decode(v)
	}
}

func (r *Request) Decode(ctx context.Context, decoder BodyDecoder) (*http.Response, error) {
	var resp, err = r.Do(ctx)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()
	//defer io.Copy(io.Discard, resp.Body)

	return resp, decoder(resp)
}

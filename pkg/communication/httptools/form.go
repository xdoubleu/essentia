package httptools

import (
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

// WriteForm writes the provided data to [url.Values].
func WriteForm(src any) (url.Values, error) {
	values := url.Values{}

	encoder := schema.NewEncoder()
	err := encoder.Encode(src, values)
	if err != nil {
		return nil, err
	}

	return values, nil
}

// ReadForm reads url-encoded form data and assigns this to dst.
func ReadForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	decoder := schema.NewDecoder()
	err = decoder.Decode(dst, r.PostForm)
	if err != nil {
		return err
	}

	return nil
}

package external

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func buildUrl(path string, queryParams map[string]string) (resUrl string, err error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to parse path: %s", path)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to parse query params: %s", u.RawQuery)
	}

	for key, param := range queryParams {
		q.Add(key, param)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func jsonUnmarshal(t interface{}, r *http.Response) (err error) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&t); err != nil {
		return errors.Wrapf(err, "Failed to Unmarshal data to type: %T", &t)
	}

	return
}

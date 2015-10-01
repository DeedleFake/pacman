package aur

import (
	"encoding/json"
	"github.com/DeedleFake/pacman"
	"net/http"
	"net/url"
)

var (
	AURURL = "https://aur.archlinux.org"
)

type Client struct {
	Client *http.Client
	URL    string
}

var DefaultClient = &Client{
	Client: http.DefaultClient,
	URL:    AURURL + "/rpc.php",
}

func (c *Client) buildURL(t string, a ...string) string {
	arg := "arg"
	if len(a) > 1 {
		arg = "arg[]"
	}

	q := make(url.Values, 2)
	q.Set("type", t)
	for _, a := range a {
		q.Add(arg, a)
	}

	return c.URL + "?" + q.Encode()
}

func (c *Client) get(data interface{}, t string, a ...string) error {
	rsp, err := c.Client.Get(c.buildURL(t, a...))
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	d := json.NewDecoder(rsp.Body)

	var raw json.RawMessage
	err = d.Decode(&raw)
	if err != nil {
		return err
	}

	var rs struct {
		Type    string          `json:"type"`
		Results json.RawMessage `json:"results"`
	}
	err = json.Unmarshal(raw, &rs)
	if err != nil {
		return err
	}
	if rs.Type == "error" {
		var e Error
		err = json.Unmarshal(rs.Results, &e)
		if err != nil {
			return err
		}

		return e
	}

	return json.Unmarshal(raw, data)
}

func Search(name string) ([]*PkgInfo, error) {
	return DefaultClient.Search(name)
}

func (c *Client) Search(name string) ([]*PkgInfo, error) {
	var data struct {
		Results []*PkgInfo `json:"results"`
	}
	err := c.get(&data, "search", name)
	if err != nil {
		return nil, err
	}

	return data.Results, nil
}

func Info(name string) (*PkgInfo, error) {
	return DefaultClient.Info(name)
}

func (c *Client) Info(name string) (*PkgInfo, error) {
	var data struct {
		ResultCount int             `json:"resultcount"`
		Results     json.RawMessage `json:"results"`
	}
	err := c.get(&data, "info", name)
	if err != nil {
		return nil, err
	}

	if data.ResultCount == 0 {
		return nil, nil
	}

	var result PkgInfo
	err = json.Unmarshal(data.Results, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func MultiInfo(names ...string) ([]*PkgInfo, error) {
	return DefaultClient.MultiInfo(names...)
}

func (c *Client) MultiInfo(names ...string) ([]*PkgInfo, error) {
	var data struct {
		Results []*PkgInfo `json:"results"`
	}
	err := c.get(&data, "multiinfo", names...)
	if err != nil {
		return nil, err
	}

	return data.Results, nil
}

func MSearch(name string) ([]*PkgInfo, error) {
	return DefaultClient.MSearch(name)
}

func (c *Client) MSearch(name string) ([]*PkgInfo, error) {
	var data struct {
		Results []*PkgInfo `json:"results"`
	}
	err := c.get(&data, "msearch", name)
	if err != nil {
		return nil, err
	}

	return data.Results, nil
}

type PkgInfo struct {
	ID             int64
	Name           string
	PkgBaseID      int64  `json:"PackageBaseID"`
	PkgBase        string `json:"PackageBase"`
	Version        pacman.Version
	Desc           string `json:"Description"`
	URL            string
	NumVotes       int
	outOfDate      int `json:"OutOfDate"`
	Maintainer     string
	FirstSubmitted int64
	LastModified   int64
	License        string
	URLPath        string
	CategoryID     int
	Popularity     float64
}

type Error string

func (err Error) Error() string {
	return "RPC Error: " + string(err)
}

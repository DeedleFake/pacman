package aur

import (
	"encoding/json"
	"github.com/DeedleFake/pacman"
	"net/http"
	"net/url"
)

const (
	// The default URL of the AUR.
	AURURL = "https://aur.archlinux.org"
)

// Client provides an interface for accessing the AUR.
type Client struct {
	Client *http.Client
	RPC    string
}

// DefaultClient is a Client with default values.
var DefaultClient = &Client{
	Client: http.DefaultClient,
	RPC:    AURURL + "/rpc.php",
}

func (c *Client) buildRPC(t string, a ...string) string {
	arg := "arg"
	if len(a) > 1 {
		arg = "arg[]"
	}

	q := make(url.Values, 2)
	q.Set("type", t)
	for _, a := range a {
		q.Add(arg, a)
	}

	return c.RPC + "?" + q.Encode()
}

func (c *Client) get(data interface{}, t string, a ...string) error {
	rsp, err := c.Client.Get(c.buildRPC(t, a...))
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

// Search calls DefaultClient.Search().
func Search(arg string) ([]*PkgInfo, error) {
	return DefaultClient.Search(arg)
}

// Search searches the AUR for packages whose names or descriptions
// contain arg. If no packages are found, it returns an empty slice
// and nil.
func (c *Client) Search(arg string) ([]*PkgInfo, error) {
	var data struct {
		Results []*PkgInfo `json:"results"`
	}
	err := c.get(&data, "search", arg)
	if err != nil {
		return nil, err
	}

	return data.Results, nil
}

// Info calls DefaultClient.Info().
func Info(arg string) (*PkgInfo, error) {
	return DefaultClient.Info(arg)
}

// Info retrieves information about one package from the AUR. If no
// such package exists, it returns nil, nil.
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

// MultiInfo calls DefaultClient.MultiInfo().
func MultiInfo(args ...string) ([]*PkgInfo, error) {
	return DefaultClient.MultiInfo(args...)
}

// MultiInfo gets information about multiple packages from the AUR. If
// a package doesn't exist, it is simply omitted from the returned
// slice. If none of the packages exist, the returned slice will be
// empty.
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

// MSearch calls DefaultClient.MSearch().
func MSearch(arg string) ([]*PkgInfo, error) {
	return DefaultClient.MSearch(arg)
}

// MSearch searches the AUR for all packages that have the given
// maintainer. If no such packages exist, it returns an empty slice
// and nil.
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

// PkgInfo contains information about a package entry in the AUR.
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

// OutOfDate wraps the OutOfDate parameter returned by the AUR's RPC
// mechanism. The AUR returns an int, so this method is a convience to
// convert that to a bool.
func (pkg PkgInfo) OutOfDate() bool {
	return pkg.outOfDate != 0
}

// Error represents an error returned from an RPC call to the AUR.
type Error string

func (err Error) Error() string {
	return "RPC Error: " + string(err)
}

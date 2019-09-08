package basecampclient

import (
	"fmt"
	"net/http"

	"github.com/rs/xid"
	"golang.org/x/oauth2"
)

var (
	// AccessTypeWebserver adds a query parameter that sets type to web_server
	accessTypeWebserver = oauth2.SetAuthURLParam("type", "web_server")
)

// Client wraps all necessary data for accessing basecamp api
type Client struct {
	email       string
	appName     string
	oauthConfig *oauth2.Config
	state       string
	code        string
	token       *oauth2.Token
	httpclient  *http.Client
}

// New returns a new basecamp client
func New(clientID, clientSecret, email, appName, callbackURL string) *Client {
	return &Client{
		appName: appName,
		state:   xid.New().String(),
		code:    "",
		token:   &oauth2.Token{},
		oauthConfig: &oauth2.Config{
			RedirectURL:  callbackURL,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{},
			Endpoint: oauth2.Endpoint{
				AuthStyle: oauth2.AuthStyleAutoDetect,
				AuthURL:   "https://launchpad.37signals.com/authorization/new",
				TokenURL:  "https://launchpad.37signals.com/authorization/token",
			},
		},
		httpclient: http.DefaultClient,
	}
}

// AuthCodeURL returns a url, that asks bascamp for access
func (c *Client) AuthCodeURL() string {
	return c.oauthConfig.AuthCodeURL(c.state, accessTypeWebserver)
}

// HandleCallback extracts the code from basecamp and exchages it for a token
func (c *Client) HandleCallback(request *http.Request) {
	code := request.FormValue("code")
	state := request.FormValue("state")
	if state != c.state {
		fmt.Println("[basecampclient/client.go/handleCallback] state doesn't match")
		return
	}
	t, err := c.oauthConfig.Exchange(oauth2.NoContext, code, accessTypeWebserver)
	if err != nil {
		fmt.Println("[basecampclient/client.go/handleCallback] couldn't exchange token: " + err.Error())
		return
	}
	c.token = t
}

func (c *Client) addHeader(request *http.Request) {
	request.Header.Add("Authorization", "Bearer "+c.token.AccessToken)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", fmt.Sprintf("%s (%s)", c.appName, c.email))
}

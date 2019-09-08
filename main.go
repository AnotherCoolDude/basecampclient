package main

import (
	"fmt"
	"net/http"

	"github.com/AnotherCoolDude/basecampclient/basecampclient"
)

var (
	client *basecampclient.Client
)

func main() {
	client = basecampclient.New(
		"34a8b658744a7f190a03765149bbcf2282161add",
		"--",
		"christian.hovenbitzer@selinka-schmitz.de",
		"ProadInterface",
		"http://localhost:8080/callback")

	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleBasecampLogin)
	http.HandleFunc("/callback", handleBasecampCallback)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html>
<body>
	<a href="/login">Basecamp Log In</a>
</body>
</html>`
	fmt.Fprintf(w, htmlIndex)
}

func handleBasecampLogin(w http.ResponseWriter, r *http.Request) {
	url := client.AuthCodeURL()
	//url := oauthConfig.AuthCodeURL(oauthStateString, accessTypeWebserver)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleBasecampCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Code value " + r.FormValue("code"))
	client.HandleCallback(r)
	//content, err := getAuthInfo(r.FormValue("state"), r.FormValue("code"))
	client.Get("https://launchpad.37signals.com/authorization.json")

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }
	// fmt.Fprintf(w, "Content: %s\n", content)
}

// func getAuthInfo(state string, code string) ([]byte, error) {
// 	if state != oauthStateString {
// 		return nil, fmt.Errorf("invalid oauth state")
// 	}
// 	token, err := oauthConfig.Exchange(oauth2.NoContext, code, accessTypeWebserver)
// 	if err != nil {
// 		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
// 	}
// 	authURL := "https://launchpad.37signals.com/authorization.json"
// 	client := http.DefaultClient
// 	response, err := client.Do(basecampGET(authURL, token.AccessToken))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
// 	}
// 	defer response.Body.Close()
// 	contents, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
// 	}
// 	return contents, nil
// }

func basecampGET(url, token string) *http.Request {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		fmt.Println("error making request")
		fmt.Println(err)
		return &http.Request{}
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "ProadInterface (-.-@selinka-schmitz.de)")
	return req
}

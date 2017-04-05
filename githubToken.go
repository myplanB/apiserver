package apiserver

import (
    "fmt"
    "net/http"

    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
    githuboauth "golang.org/x/oauth2/github"
    "context"
)

var (
    // You must register the app at https://github.com/settings/applications/new
    // Set callback to http://127.0.0.1:7000/github_oauth_cb
    // Set ClientId and ClientSecret to
    oauthConf = &oauth2.Config{
        ClientID:     "17f03800ee9e0b7caf95",
        ClientSecret: "542dd6ffef2da3ac66ed493e82fba8fd0de089ba",
        //ClientID:     "",
        //ClientSecret: "",
        // select level of access you want https://developer.github.com/v3/oauth/#scopes
        Scopes:       []string{"user:email", "repo"},
        Endpoint:     githuboauth.Endpoint,
    }
    // random string for oauth2 API calls to protect against CSRF
    oauthStateString = "thisshouldberandom"
)

const htmlIndex = `<html><body>
Logged in with <a href="/login">GitHub</a>
</body></html>
`

// /
func handleMain(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(htmlIndex))
}

// /login
func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
    url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// /github_oauth_cb. Called by github after authorization is granted
func handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
    state := r.FormValue("state")
    if state != oauthStateString {
        fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    code := r.FormValue("code")
    token, err := oauthConf.Exchange(oauth2.NoContext, code)
    if err != nil {
        fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    oauthClient := oauthConf.Client(oauth2.NoContext, token)

    client := github.NewClient(oauthClient)

    user, _, err := client.Users.Get(context.Background(),"")
    if err != nil {
        fmt.Printf("client.Users.Get() faled with '%s'\n", err)
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }
    fmt.Printf("Logged in as GitHub user: %s\n", *user.Login)
	// list all repositories for the authenticated user
    repos, _, err := client.Repositories.List(context.Background(), "", nil)

    for k,v := range repos{
        fmt.Println(k,v)
    }

    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}


func main() {
    http.HandleFunc("/", handleMain)
    http.HandleFunc("/login", handleGitHubLogin)
    http.HandleFunc("/github_oauth_cb", handleGitHubCallback)
    fmt.Print("Started running on http://127.0.0.1:7000\n")
    fmt.Println(http.ListenAndServe(":7000", nil))
}

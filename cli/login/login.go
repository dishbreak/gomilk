package login

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/urfave/cli"

	"github.com/pkg/browser"

	"github.com/dishbreak/gomilk/api"

	"github.com/dishbreak/gomilk/api/auth"

	"github.com/mitchellh/go-homedir"
)

/*
GetUserTokenError commmunicates where GetUserToken tried looking for a token.
*/
type GetUserTokenError struct {
	TokenFilePath string
}

func (e *GetUserTokenError) Error() string {
	return fmt.Sprintf("Failed to find token at path: %s", e.TokenFilePath)
}

func tokenFilePath() string {
	userdir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	tokenPath := filepath.Join(userdir, ".gomilk", "token")

	return tokenPath
}

/*
GetUserToken retrieves the token from the filesystem.

Don't presume this token is correct unless you've called IsAuthenticated() first
*/
func GetUserToken() (string, error) {
	tokenPath := tokenFilePath()
	buffer, err := ioutil.ReadFile(tokenPath)

	if err != nil {
		return "", &GetUserTokenError{tokenPath}
	}

	return string(buffer), nil
}

/*
Token is the authentication token for RTM. Don't presume it's valid or it exists!
Use IsAuthenticated() to check for validity.
*/
var Token string

/*
Setup will initialize Token with the value from the filesystem.
This will let downstream code use Token to make authenticated requests.
*/
func Setup() error {
	Token, err := GetUserToken()
	_ = Token // This is a module-level variable, so I'm "using" it here.
	return err
}

func setUserToken(token string) error {
	tokenPath := tokenFilePath()
	gomilkDir := path.Dir(tokenPath)
	if stat, err := os.Stat(gomilkDir); os.IsNotExist(err) {
		err = os.Mkdir(gomilkDir, 0755)
		if err != nil {
			return err
		}
	} else if mode := stat.Mode(); !mode.IsDir() {
		return fmt.Errorf("needed to create dir '%s' but it exists already as a file", gomilkDir)
	}

	err := ioutil.WriteFile(tokenPath, []byte(token), 0600)
	if err != nil {
		return err
	}

	return nil
}

/*
IsAuthenticated will check if we have a valid token present.
*/
func IsAuthenticated() bool {
	token, err := GetUserToken()
	if err != nil {
		if _, ok := err.(*GetUserTokenError); ok {
			return false
		}
		panic(err)
	}

	_, err = auth.CheckToken(token)
	if err != nil {
		if err, ok := err.(*api.RTMAPIError); ok {
			if err.Rsp.Err.Code != auth.ERROR_INVALID_AUTH_TOKEN {
				panic(err)
			}
			return false
		}
	}
	return true
}

const (
	rootURL = "https://www.rememberthemilk.com/services/auth/"
)

/*
Attempt to log the user into Remember the Milk
*/
func Login(c *cli.Context) {
	if IsAuthenticated() {
		fmt.Println("Logged in! We're good to go. Moo!")
		return
	}

	u, err := url.Parse(rootURL)
	if err != nil {
		panic(err)
	}

	frob, err := auth.GetFrob()
	if err != nil {
		panic(err)
	}

	params := map[string]string{
		"api_key": api.APIKey,
		"perms":   auth.Delete.String(),
		"frob":    frob.Rsp.Frob,
	}

	q := u.Query()

	for param, val := range params {
		q.Set(param, val)
	}

	q.Set("api_sig", api.SignRequest(params))

	u.RawQuery = q.Encode()

	fmt.Println("We're going to have you log in to Remember the Milk in a browser window.")
	fmt.Println("When you're done, come back to this window and hit Enter to complete the login.")
	fmt.Println("Press [Enter] to continue...")
	fmt.Scanln()

	fmt.Println("Sending you to Remember the Milk now...")
	browser.OpenURL(u.String())
	fmt.Println("Press [Enter] to finish logging in.")
	fmt.Scanln()

	token, err := auth.GetToken(frob.Rsp.Frob)
	if err != nil {
		panic(err)
	}

	err = setUserToken(token.Rsp.Auth.Token)
	if err != nil {
		panic(err)
	}

	fmt.Printf("You've now logged in!")
	fmt.Printf("If you'd like us to forget your login, delete the following file:")
	fmt.Printf("\t%s\n", tokenFilePath())
	fmt.Println("Don't forget to buy some milk. :)")
}

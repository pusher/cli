package channels

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/pusher/cli/api"
	"github.com/pusher/cli/commands"
	"github.com/pusher/cli/commands/auth"
	pusher "github.com/pusher/pusher-http-go"
	"github.com/spf13/cobra"
)

var localAuthServerPort int

//LocalAuthServer starts a server locally that authenticates all requests.
var LocalAuthServer = &cobra.Command{
	Use:   "auth-server",
	Short: "Run a local auth server that authenticates all requests",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if !auth.APIKeyValid() {
			fmt.Println("Your API key isn't valid. Add one with the `login` command.")
			os.Exit(1)
			return
		}

		if commands.AppID == "" {
			fmt.Fprintf(os.Stderr, "Please supply --app-id\n")
			os.Exit(1)
			return
		}

		app, err := api.GetApp(commands.AppID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get the app: %s\n", err.Error())
			os.Exit(1)
			return
		}

		token, err := api.GetToken(commands.AppID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get app token: %s\n", err.Error())
			os.Exit(1)
			return
		}

		pClient := pusher.Client{
			AppID:   fmt.Sprintf("%d", app.ID),
			Key:     token.Key,
			Secret:  token.Secret,
			Cluster: app.Cluster,
		}

		http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
			params, _ := ioutil.ReadAll(req.Body)
			response, err := pClient.AuthenticatePrivateChannel(params)

			if err != nil {
				fmt.Println("Invalid request", err)
				return
			}

			resp.Header().Set("Access-Control-Allow-Origin", "*")

			fmt.Fprintf(resp, string(response))
		})

		portColor := color.New(color.FgBlue)
		fmt.Print("Started local server. Listening on port ")
		portColor.Printf("%d", localAuthServerPort)
		fmt.Print(".\n")

		http.ListenAndServe(fmt.Sprintf(":%d", localAuthServerPort), nil)
	},
}

func init() {
	LocalAuthServer.PersistentFlags().IntVar(&localAuthServerPort, "port", 8080, "")
	LocalAuthServer.PersistentFlags().StringVar(&commands.AppID, "app-id", "", "Pusher App ID")
}

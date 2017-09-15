package commands

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/pusher/pusher-cli/api"
	"github.com/pusher/pusher-cli/config"
	"github.com/pusher/pusher-http-go"
	"github.com/spf13/cobra"
)

var localAuthServerPort int
var LocalAuthServer = &cobra.Command{
	Use:   "local-auth-server",
	Short: "Run a local auth server that authenticates all requests",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if config.Get().Email == "" || config.Get().Password == "" {
			fmt.Printf("Not logged in as '%s'.\n", config.Get().Email)
			os.Exit(1)
			return
		}

		if appId == "" {
			fmt.Fprintf(os.Stderr, "Please supply --app-id\n")
			os.Exit(1)
			return
		}

		app, err := api.GetApp(appId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get the app: %s\n", err.Error())
			os.Exit(1)
			return
		}

		token, err := api.GetToken(appId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get token: %s\n", err.Error())
			os.Exit(1)
			return
		}

		pClient := pusher.Client{
			AppId:   fmt.Sprintf("%d", app.Id),
			Key:     token.Key,
			Secret:  token.Secret,
			Cluster: app.Cluster,
		}

		http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
			params, _ := ioutil.ReadAll(req.Body)
			response, err := pClient.AuthenticatePrivateChannel(params)

			if err != nil {
				panic(err)
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
	LocalAuthServer.PersistentFlags().StringVar(&appId, "app-id", "", "Pusher App ID")
}

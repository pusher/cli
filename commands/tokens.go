package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"github.com/pusher/pusher-cli/api"
)

var Tokens = &cobra.Command{
	Use:   "tokens",
	Short: "List tokens for a Pusher app",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if appId == "" {
			fmt.Fprintf(os.Stderr, "Please supply --app-id\n")
			os.Exit(1)
			return
		}

		tokens, err := api.GetAllTokensForApp(appId)
		if err != nil {
			fmt.Printf("Failed to retrieve the list of tokens: %s\n", err.Error())
			os.Exit(1)
			return
		}

		if outputAsJson {
			tokensJsonBytes, _ := json.Marshal(tokens)
			fmt.Println(string(tokensJsonBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Key", "Secret"})
			for _, token := range tokens {
				table.Append([]string{token.Key, token.Secret})
			}
			table.Render()
		}
	},
}

func init() {
	Tokens.PersistentFlags().StringVar(&appId, "app-id", "", "Pusher App ID")
	Tokens.PersistentFlags().BoolVar(&outputAsJson, "json", false, "")
}

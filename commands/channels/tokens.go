package channels

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/olekukonko/tablewriter"
	"github.com/pusher/cli/api"
	"github.com/pusher/cli/commands"
	"github.com/spf13/cobra"
)

// Tokens lists the App Key and Secret for a particular app.
var Tokens = &cobra.Command{
	Use:   "tokens",
	Short: "List tokens for a Pusher app",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if commands.AppID == "" {
			fmt.Fprintf(os.Stderr, "Please supply --app-id\n")
			os.Exit(1)
			return
		}

		tokens, err := api.GetAllTokensForApp(commands.AppID)
		if err != nil {
			fmt.Printf("Failed to retrieve the list of tokens: %s\n", err.Error())
			os.Exit(1)
			return
		}

		if commands.OutputAsJSON {
			tokensJSONBytes, _ := json.Marshal(tokens)
			fmt.Println(string(tokensJSONBytes))
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
	Tokens.PersistentFlags().StringVar(&commands.AppID, "app-id", "", "Pusher App ID")
	Tokens.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
}

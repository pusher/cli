package channels

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"encoding/json"

	"github.com/olekukonko/tablewriter"
	"github.com/pusher/cli/api"
	"github.com/pusher/cli/commands"
	"github.com/spf13/cobra"
)

func NewConfigListCommand(pusher api.FunctionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List function configs for an Channels app",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			configs, err := pusher.GetFunctionConfigsForApp(commands.AppID)
			if err != nil {
				return err
			}

			if commands.OutputAsJSON {
				configsJSONBytes, _ := json.Marshal(configs)
				cmd.Println(string(configsJSONBytes))
			} else {
				table := newTable(cmd.OutOrStdout())
				table.SetHeader([]string{"Name", "Desciption", "Type"})
				for _, config := range configs {
					table.Append([]string{config.Name, config.Description, config.ParamType})
				}
				table.Render()
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
	return cmd
}

func NewConfigCreateCommand(functionService api.FunctionService) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a function config for a Channels app",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := functionService.CreateFunctionConfig(commands.AppID, commands.FunctionConfigName, commands.FunctionConfigDescription, commands.FunctionConfigParamType, commands.FunctionConfigContent)
			if err != nil {
				return err
			}

			if commands.OutputAsJSON {
				functionJSONBytes, _ := json.Marshal(config)
				fmt.Fprintln(cmd.OutOrStdout(), string(functionJSONBytes))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "created function config %s\n", config.Name)
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigName, "name", "", "Function config name. Can only contain A-Za-z0-9-_")
	err := cmd.MarkPersistentFlagRequired("name")
	if err != nil {
		return nil, err
	}
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigDescription, "description", "", "Function config description")
	err = cmd.MarkPersistentFlagRequired("description")
	if err != nil {
		return nil, err
	}
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigParamType, "type", "", "Function config type, valid options: param|secret")
	err = cmd.MarkPersistentFlagRequired("type")
	if err != nil {
		return nil, err
	}
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigContent, "content", "", "Function config contents")
	err = cmd.MarkPersistentFlagRequired("content")
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func NewConfigUpdateCommand(functionService api.FunctionService) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a function config for a Channels app",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := functionService.UpdateFunctionConfig(commands.AppID, commands.FunctionConfigName, commands.FunctionConfigDescription, commands.FunctionConfigContent)
			if err != nil {
				return err
			}

			if commands.OutputAsJSON {
				functionJSONBytes, _ := json.Marshal(config)
				fmt.Fprintln(cmd.OutOrStdout(), string(functionJSONBytes))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "updated function config %s\n", config.Name)
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigName, "name", "", "Function config name. Can only contain A-Za-z0-9-_")
	err := cmd.MarkPersistentFlagRequired("name")
	if err != nil {
		return nil, err
	}
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigDescription, "description", "", "Function config description")
	err = cmd.MarkPersistentFlagRequired("description")
	if err != nil {
		return nil, err
	}
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigContent, "content", "", "Function config contents")
	err = cmd.MarkPersistentFlagRequired("content")
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func NewConfigDeleteCommand(functionService api.FunctionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a function config from a Channels app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := functionService.DeleteFunctionConfig(commands.AppID, args[0])
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "deleted function config %s\n", args[0])
			return nil
		},
	}
	return cmd
}

func NewConfigCommand(pusher api.FunctionService) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "configs",
		Short: "Manage function config params for a Channels app",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(NewConfigListCommand(pusher))
	c, err := NewConfigCreateCommand(pusher)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(c)
	c, err = NewConfigUpdateCommand(pusher)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(c)
	cmd.AddCommand(NewConfigDeleteCommand(pusher))
	return cmd, nil
}

func NewFunctionsCommand(pusher api.FunctionService) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "functions",
		Short: "Manage functions for a Channels app",
		Args:  cobra.NoArgs,
	}
	cmd.PersistentFlags().StringVar(&commands.AppID, "app-id", "", "Channels App ID")
	err := cmd.MarkPersistentFlagRequired("app-id")
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(NewFunctionsListCommand(pusher))

	c, err := NewFunctionsCreateCommand(pusher)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(c)
	cmd.AddCommand(NewFunctionDeleteCommand(pusher))
	cmd.AddCommand(NewFunctionGetCommand(pusher))
	c, err = NewFunctionsUpdateCommand(pusher)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(c)
	cmd.AddCommand(NewFunctionGetLogsCommand(pusher))
	c, err = NewFunctionInvokeCommand(pusher)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(c)
	c, err = NewConfigCommand(pusher)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(c)
	return cmd, nil
}

var Functions = &cobra.Command{
	Use:   "functions",
	Short: "Manage functions for a Channels app",
	Args:  cobra.NoArgs,
}

func newTable(out io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(out)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetHeaderLine(false)
	table.SetColumnSeparator("")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	return table
}

func NewFunctionsListCommand(functionService api.FunctionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List functions for an Channels app",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			functions, err := functionService.GetAllFunctionsForApp(commands.AppID)
			if err != nil {
				return err
			}

			if commands.OutputAsJSON {
				functionsJSONBytes, _ := json.Marshal(functions)
				cmd.Println(string(functionsJSONBytes))
			} else {
				table := newTable(cmd.OutOrStdout())
				table.SetHeader([]string{"Name", "Mode", "Events"})
				for _, function := range functions {
					table.Append([]string{function.Name, function.Mode, strings.Join(function.Events, ",")})
				}
				table.Render()
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
	return cmd
}

func cleanMode(m string) string {
	switch strings.ToLower(m) {
	case "sync", "synch", "synchronous":
		return "synchronous"
	case "async", "asynch", "asynchronous":
		return "asynchronous"
	default:
		return m
	}
}

func NewFunctionsCreateCommand(functionService api.FunctionService) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "create <path to code directory>",
		Short: "Create a function for a Channels app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			archive := ZipFolder(args[0])

			mode := cleanMode(commands.FunctionMode)
			if mode == "" {
				mode = "asynchronous"
			}

			function, err := functionService.CreateFunction(commands.AppID, commands.FunctionName, commands.FunctionEvents, archive, mode)
			if err != nil {
				return err
			}

			if commands.OutputAsJSON {
				functionJSONBytes, _ := json.Marshal(function)
				fmt.Fprintln(cmd.OutOrStdout(), string(functionJSONBytes))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "created function %s\n", function.Name)
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
	cmd.PersistentFlags().StringVar(&commands.FunctionName, "name", "", "Function name")
	err := cmd.MarkPersistentFlagRequired("name")
	if err != nil {
		return nil, err
	}
	cmd.PersistentFlags().StringVar(&commands.FunctionMode, "mode", "asynchronous", "Function mode. Either synchronous or asynchronous")
	cmd.PersistentFlags().StringSliceVar(&commands.FunctionEvents, "events", nil, "Channel events that trigger this function")
	return cmd, err
}

func NewFunctionDeleteCommand(functionService api.FunctionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <function_name>",
		Short: "Delete a function from a Channels app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := functionService.DeleteFunction(commands.AppID, args[0])
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "deleted function %s\n", args[0])
			return nil
		},
	}
	return cmd
}

func NewFunctionInvokeCommand(functionService api.FunctionService) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "invoke <function_name>",
		Short: "invoke a function to test it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := functionService.InvokeFunction(commands.AppID, args[0], commands.Data, commands.EventName, commands.ChannelName)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%v\n", result)
			return nil
		},
	}
	cmd.PersistentFlags().StringVar(&commands.Data, "data", "", "Channels event data")
	err := cmd.MarkPersistentFlagRequired("data")
	if err != nil {
		return nil, err
	}
	cmd.PersistentFlags().StringVar(&commands.EventName, "event", "", "Channels event name")
	err = cmd.MarkPersistentFlagRequired("event")
	if err != nil {
		return nil, err
	}
	cmd.PersistentFlags().StringVar(&commands.ChannelName, "channel", "", "Channels name")
	err = cmd.MarkPersistentFlagRequired("channel")
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func NewFunctionGetCommand(functionService api.FunctionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <function_name>",
		Short: "Downloads a function from a Channels app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fn, err := functionService.GetFunction(commands.AppID, args[0])
			if err != nil {
				return err
			}

			if commands.OutputAsJSON {
				functionsJSONBytes, _ := json.Marshal(fn)
				cmd.Println(string(functionsJSONBytes))
			} else {
				zipFileName := fmt.Sprintf("%s.%s.zip", fn.Name, time.Now().Format("2006-01-02-150405"))
				fmt.Fprintf(cmd.OutOrStdout(), "ID: %v\n", fn.ID)
				fmt.Fprintf(cmd.OutOrStdout(), "Name: %v\n", fn.Name)
				fmt.Fprintf(cmd.OutOrStdout(), "Mode: %v\n", fn.Mode)
				fmt.Fprintf(cmd.OutOrStdout(), "Events: %v\n", strings.Join(fn.Events, ","))
				err = os.WriteFile(zipFileName, fn.Body, 0644)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Body: '%v'\n", zipFileName)
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
	return cmd
}

func NewFunctionsUpdateCommand(functionService api.FunctionService) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "update <function_name> [<path to code directory>]",
		Short: "Update a function for a Channels app",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var archive io.Reader = nil
			if len(args) >= 2 {
				archive = ZipFolder(args[1])
			}

			function, err := functionService.UpdateFunction(commands.AppID, args[0], commands.FunctionName, commands.FunctionEvents, archive, cleanMode(commands.FunctionMode))
			if err != nil {
				return err
			}

			if commands.OutputAsJSON {
				functionJSONBytes, _ := json.Marshal(function)
				fmt.Fprintln(cmd.OutOrStdout(), string(functionJSONBytes))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "updated function: %v\n", function.ID)
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
	cmd.PersistentFlags().StringVar(&commands.FunctionName, "name", "", "Function name")
	cmd.PersistentFlags().StringSliceVar(&commands.FunctionEvents, "events", nil, "Channel events that trigger this function")
	cmd.PersistentFlags().StringVar(&commands.FunctionMode, "mode", "", "Function mode. Either synchronous or asynchronous")
	return cmd, nil
}

func NewFunctionGetLogsCommand(functionService api.FunctionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs <function_name>",
		Short: "Get logs of a specific function",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logs, err := functionService.GetFunctionLogs(commands.AppID, args[0])
			if err != nil {
				return err
			}

			if commands.OutputAsJSON {
				JSONBytes, _ := json.Marshal(logs)
				cmd.Println(string(JSONBytes))

				return nil
			}

			for _, l := range logs.Events {
				t := time.Unix(0, l.Timestamp*1000000).Format("2006-01-02 15:04:05")
				cmd.Printf("%s\t%s\n", t, l.Message)
			}

			return nil
		},
	}

	cmd.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")

	return cmd
}

func ZipFolder(baseFolder string) io.Reader {
	r, w := io.Pipe()

	go func() {
		// Create a new zip archive.
		zw := zip.NewWriter(w)

		// Recursively add files to the archive.
		err := addFiles(zw, baseFolder, "")

		// Close the archive and pipeline, reporting any errors.
		if err != nil {
			zw.Close()
		} else {
			err = zw.Close()
		}
		w.CloseWithError(err)
	}()

	return r
}

func addFiles(w *zip.Writer, basePath, baseInZip string) error {
	// Open the Directory.
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(filepath.Join(basePath, file.Name()))
			if err != nil {
				return err
			}

			// Add some files to the archive.
			f, err := w.Create(filepath.Join(baseInZip, file.Name()))
			if err != nil {
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				return err
			}
		} else if file.IsDir() {

			// Recurse.
			newBase := filepath.Join(basePath, file.Name())
			err = addFiles(w, newBase, filepath.Join(baseInZip, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

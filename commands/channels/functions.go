package channels

import (
	"fmt"
	"io"
	"io/fs"
	"strconv"
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

func NewConfigCreateCommand(functionService api.FunctionService) *cobra.Command {
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
				fmt.Fprintf(cmd.OutOrStdout(), "created function config %s", config.Name)
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigName, "name", "", "Function config name. Can only contain A-Za-z0-9-_")
	cmd.MarkPersistentFlagRequired("name")
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigDescription, "description", "", "Function config description")
	cmd.MarkPersistentFlagRequired("description")
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigParamType, "type", "", "Function config type, valid options: param|secret")
	cmd.MarkPersistentFlagRequired("type")
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigContent, "content", "", "Function config contents")
	cmd.MarkPersistentFlagRequired("contents")
	return cmd
}

func NewConfigUpdateCommand(functionService api.FunctionService) *cobra.Command {
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
				fmt.Fprintf(cmd.OutOrStdout(), "updated function config %s", config.Name)
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigName, "name", "", "Function config name. Can only contain A-Za-z0-9-_")
	cmd.MarkPersistentFlagRequired("name")
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigDescription, "description", "", "Function config description")
	cmd.MarkPersistentFlagRequired("description")
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigContent, "content", "", "Function config contents")
	cmd.MarkPersistentFlagRequired("contents")
	return cmd
}

func NewConfigDeleteCommand(functionService api.FunctionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a function config from a Channels app",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := functionService.DeleteFunctionConfig(commands.AppID, commands.FunctionConfigName)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "deleted function config %s\n", commands.FunctionConfigName)
			return nil
		},
	}
	cmd.PersistentFlags().StringVar(&commands.FunctionConfigName, "name", "", "Function config name. Can only contain A-Za-z0-9-_")
	cmd.MarkPersistentFlagRequired("name")
	return cmd
}

func NewConfigCommand(pusher api.FunctionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configs",
		Short: "Manage function config params for a Channels app",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(NewConfigListCommand(pusher))
	cmd.AddCommand(NewConfigCreateCommand(pusher))
	cmd.AddCommand(NewConfigUpdateCommand(pusher))
	cmd.AddCommand(NewConfigDeleteCommand(pusher))
	return cmd
}

func NewFunctionsCommand(pusher api.FunctionService, fs fs.ReadFileFS) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "functions",
		Short: "Manage functions for a Channels app",
		Args:  cobra.NoArgs,
	}
	cmd.PersistentFlags().StringVar(&commands.AppID, "app_id", "", "Channels App ID")
	cmd.MarkPersistentFlagRequired("app_id")
	cmd.AddCommand(NewFunctionsListCommand(pusher))
	cmd.AddCommand(NewFunctionsCreateCommand(pusher, fs))
	cmd.AddCommand(NewFunctionDeleteCommand(pusher))
	cmd.AddCommand(NewFunctionGetCommand(pusher))
	cmd.AddCommand(NewFunctionsUpdateCommand(pusher, fs))
	cmd.AddCommand(NewFunctionGetLogsCommand(pusher))
	cmd.AddCommand(NewConfigCommand(pusher))
	return cmd
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
				table.SetHeader([]string{"ID", "Name", "Events"})
				for _, function := range functions {
					table.Append([]string{strconv.Itoa(function.ID), function.Name, strings.Join(function.Events, ",")})
				}
				table.Render()
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
	return cmd
}

func NewFunctionsCreateCommand(functionService api.FunctionService, fs fs.ReadFileFS) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <path to code file>",
		Short: "Create a function for a Channels app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			code, err := fs.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("could not create function: %s does not exist", args[0])
			}

			function, err := functionService.CreateFunction(commands.AppID, commands.FunctionName, commands.FunctionEvents, string(code))
			if err != nil {
				return err
			}

			if commands.OutputAsJSON {
				functionJSONBytes, _ := json.Marshal(function)
				fmt.Fprintln(cmd.OutOrStdout(), string(functionJSONBytes))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "created function %s with id: %d\n", function.Name, function.ID)
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
	cmd.PersistentFlags().StringVar(&commands.FunctionName, "name", "", "Function name")
	cmd.MarkPersistentFlagRequired("name")
	cmd.PersistentFlags().StringSliceVar(&commands.FunctionEvents, "events", []string{}, "Channel events that trigger this function")
	return cmd
}

func NewFunctionDeleteCommand(functionService api.FunctionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <function_id>",
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

func NewFunctionGetCommand(functionService api.FunctionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <function_id>",
		Short: "Get a function for a Channels app",
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
				fmt.Fprintf(cmd.OutOrStdout(), "ID: %v\n", fn.ID)
				fmt.Fprintf(cmd.OutOrStdout(), "Name: %v\n", fn.Name)
				fmt.Fprintf(cmd.OutOrStdout(), "Events: %v\n", strings.Join(fn.Events, ","))
				fmt.Fprintf(cmd.OutOrStdout(), "Body: %v\n", fn.Body)
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
	return cmd
}

func NewFunctionsUpdateCommand(functionService api.FunctionService, fs fs.ReadFileFS) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <function_id> <path to code file>",
		Short: "Update a function for a Channels app",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[1]
			code, err := fs.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("could not update function: %s does not exist", filePath)
			}

			function, err := functionService.UpdateFunction(commands.AppID, args[0], commands.FunctionName, commands.FunctionEvents, string(code))
			if err != nil {
				return err
			}

			if commands.OutputAsJSON {
				functionJSONBytes, _ := json.Marshal(function)
				fmt.Fprintln(cmd.OutOrStdout(), string(functionJSONBytes))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "updated function: %d\n", function.ID)
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
	cmd.PersistentFlags().StringVar(&commands.FunctionName, "name", "", "Function name")
	cmd.MarkPersistentFlagRequired("name")
	cmd.PersistentFlags().StringSliceVar(&commands.FunctionEvents, "events", []string{}, "Channel events that trigger this function")
	return cmd
}

func NewFunctionGetLogsCommand(functionService api.FunctionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs <function_id>",
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

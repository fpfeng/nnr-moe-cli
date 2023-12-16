package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	nnr_moe_cli "github.com/fpfeng/nnr-moe-cli/cli"
	nnr_moe_core "github.com/fpfeng/nnr-moe-cli/core"
)

func getOutputFunc(mode string) func(interface{}) {
	if mode == "json" {
		return func(response interface{}) {
			j, err := json.Marshal(response)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%v\n", string(j))
		}
	} else {
		return func(interface{}) {}
	}
}

func Dispatcher(result *nnr_moe_cli.CLIParseResult) {
	nnrMoe := nnr_moe_core.New(result.Token)
	outputFunc := getOutputFunc(result.OutputMode)
	switch result.InvokeAction {
	case nnr_moe_cli.InvokeServers:
		outputFunc(nnrMoe.Servers())
	case nnr_moe_cli.InvokeRules:
		outputFunc(nnrMoe.Rules())
	case nnr_moe_cli.InvokeAddRule:
		outputFunc(nnrMoe.AddRule(&result.AddRule))
	case nnr_moe_cli.InvokeEditRule:
		outputFunc(nnrMoe.EditRule(&result.EditedRule))
	case nnr_moe_cli.InvokeDeleteRule:
		outputFunc(nnrMoe.DeleteRule(&result.DeleteRule))
	case nnr_moe_cli.InvokeGetRule:
		outputFunc(nnrMoe.GetRule(&result.GetRule))
	default:
	}
}

func main() {
	cliParseResult := nnr_moe_cli.StartCLI(os.Args)
	Dispatcher(cliParseResult)
}

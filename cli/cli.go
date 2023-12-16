package cli

import (
	"fmt"
	"log"

	core "github.com/fpfeng/nnr-moe-cli/core"
	"github.com/urfave/cli/v2"
)

type TypeInvokeAction string

const (
	InvokeServers    TypeInvokeAction = "servers"
	InvokeRules      TypeInvokeAction = "rules"
	InvokeAddRule    TypeInvokeAction = "add-rule"
	InvokeEditRule   TypeInvokeAction = "rdit-rule"
	InvokeDeleteRule TypeInvokeAction = "delete-rule"
	InvokeGetRule    TypeInvokeAction = "get-rule"
	InvokeGetServer  TypeInvokeAction = "get-server"
)

type CLIParseResult struct {
	Token        string
	OutputMode   string
	InvokeAction TypeInvokeAction
	AddRule      core.RequestAddRule
	DeleteRule   core.RequestRuleRid
	EditedRule   core.RequestEditedRule
	GetRule      core.RequestRuleRid
	GetServerSid string
}

func StartCLI(args []string) *CLIParseResult {
	result := CLIParseResult{
		Token:        "",
		OutputMode:   "json",
		InvokeAction: "",
		AddRule:      core.RequestAddRule{},
		DeleteRule:   core.RequestRuleRid{},
		EditedRule:   core.RequestEditedRule{},
		GetRule:      core.RequestRuleRid{},
	}
	app := &cli.App{
		Usage:       fmt.Sprintf("CLI wrapper of https://nnr.moe v%v", Version),
		Description: HomePage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "token",
				Usage:       "https://nnr.moe/user/setting 的API`密钥`",
				Required:    true,
				Destination: &result.Token,
			},
			&cli.StringFlag{
				Name:        "output",
				Usage:       "api调用成功后打印结果，除`json`外为静默",
				Destination: &result.OutputMode,
				Value:       "json",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "servers",
				Usage: "获取所有可使用节点",
				Action: func(cCtx *cli.Context) error {
					result.InvokeAction = InvokeServers
					return nil
				},
			},
			{
				Name:  "rules",
				Usage: "获取所有规则",
				Action: func(cCtx *cli.Context) error {
					result.InvokeAction = InvokeRules
					return nil
				},
			},
			{
				Name:  "add-rule",
				Usage: "添加规则",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "sid",
						Usage:       "`节点`(源服务器)id",
						Required:    true,
						Destination: &result.AddRule.Sid,
					},
					&cli.StringFlag{
						Name:        "remote",
						Usage:       "目标服务器`域名或IP` (支持DDNS)",
						Required:    true,
						Destination: &result.AddRule.Remote,
					},
					&cli.IntFlag{
						Name:        "rport",
						Usage:       "目标`端口`",
						Required:    true,
						Destination: &result.AddRule.Rport,
					},
					&cli.StringFlag{
						Name:        "type",
						Usage:       "规则`协议`(需节点支持)",
						Required:    true,
						Destination: &result.AddRule.Type,
					},
					&cli.StringFlag{
						Name:        "name",
						Usage:       "规则`名称/备注`",
						Required:    true,
						Destination: &result.AddRule.Name,
					},
				},
				Action: func(cCtx *cli.Context) error {
					result.InvokeAction = InvokeAddRule
					return nil
				},
			},
			{
				Name:  "edit-rule",
				Usage: "编辑规则",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "rid",
						Usage:       "规则`id`",
						Required:    true,
						Destination: &result.EditedRule.Rid,
					},
					&cli.StringFlag{
						Name:        "remote",
						Usage:       "目标服务器`域名或IP` (支持DDNS)",
						Required:    true,
						Destination: &result.EditedRule.Remote,
					},
					&cli.IntFlag{
						Name:        "rport",
						Usage:       "目标`端口`",
						Destination: &result.EditedRule.Rport,
					},
					&cli.StringFlag{
						Name:        "type",
						Usage:       "规则`协议`(需节点支持)",
						Required:    true,
						Destination: &result.EditedRule.Type,
					},
					&cli.StringFlag{
						Name:        "name",
						Usage:       "规则`名称/备注`",
						Required:    true,
						Destination: &result.EditedRule.Name,
					},
				},
				Action: func(cCtx *cli.Context) error {
					result.InvokeAction = InvokeEditRule
					return nil
				},
			},
			{
				Name:  "delete-rule",
				Usage: "删除规则",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "rid",
						Usage:       "规则`id`",
						Required:    true,
						Destination: &result.DeleteRule.Rid,
					},
				},
				Action: func(cCtx *cli.Context) error {
					result.InvokeAction = InvokeDeleteRule
					return nil
				},
			},
			{
				Name:  "get-rule",
				Usage: "查看规则",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "rid",
						Usage:       "规则`id`",
						Required:    true,
						Destination: &result.GetRule.Rid,
					},
				},
				Action: func(cCtx *cli.Context) error {
					result.InvokeAction = InvokeGetRule
					return nil
				},
			},
			{
				Name:  "get-server",
				Usage: "获取单个节点",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "sid",
						Usage:       "节点`id`",
						Required:    true,
						Destination: &result.GetServerSid,
					},
				},
				Action: func(cCtx *cli.Context) error {
					result.InvokeAction = InvokeGetServer
					return nil
				},
			},
		},
	}
	if err := app.Run(args); err != nil {
		log.Fatal(err)
	}
	return &result
}

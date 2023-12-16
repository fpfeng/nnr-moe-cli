package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	nnr_moe_cli "github.com/fpfeng/nnr-moe-cli/cli"
	nnr_moe_core "github.com/fpfeng/nnr-moe-cli/core"
)

const (
	cliTestingRuleName = "cli_testing"
)

var (
	mockOSArgsBase = []string{}
	addRuleSid     = "82ebf39a-b624-463d-a4da-3d644a4749a9" // 广州CM-香港 x1.5
)

func setup() {
	currentProgramPath, _ := os.Executable()
	token := os.Getenv("NNRMOE_TOKEN")
	if len(token) == 0 {
		log.Fatal("Invalid environment variable NNRMOE_TOKEN.")
	}
	mockOSArgsBase = append(mockOSArgsBase, currentProgramPath)
	mockOSArgsBase = append(mockOSArgsBase, "--token")
	mockOSArgsBase = append(mockOSArgsBase, token)

}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func fillMockOSArgs(subArgs []string) []string {
	return append(append([]string{}, mockOSArgsBase...), subArgs...)

}

func captureOutput(f func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	return string(out)
}

func startCLIWithMockArgsThenCaptureResult(subArgs []string) string {
	fullArgs := fillMockOSArgs(subArgs)
	cliParseResult := nnr_moe_cli.StartCLI(fullArgs)
	return captureOutput(func() { Dispatcher(cliParseResult) })
}

func rules() string {
	return startCLIWithMockArgsThenCaptureResult([]string{"rules"})
}

func addRule() string {
	return startCLIWithMockArgsThenCaptureResult([]string{
		"add-rule",
		"--sid",
		addRuleSid,
		"--remote",
		"8.8.8.8",
		"--rport",
		"22",
		"--type",
		"tcp",
		"--name",
		cliTestingRuleName,
	})
}

func checkStatus(stdout string, t *testing.T) {
	if !strings.Contains(stdout, `"Status":1`) {
		t.Error("Invalid status.")
	}
}

func TestServers(t *testing.T) {
	result := startCLIWithMockArgsThenCaptureResult([]string{"servers"})
	checkStatus(result, t)
	if !strings.Contains(result, "香港") {
		t.Error("servers should contains Hongkong.")
	}
}

func TestRules(t *testing.T) {
	result := rules()
	checkStatus(result, t)
}

func TestAddRule(t *testing.T) {
	result := addRule()
	checkStatus(result, t)
	if !strings.Contains(result, "rid") {
		t.Error("Add rule failed.")
	}
}

func TestEditRule(t *testing.T) {
	result := addRule()
	checkStatus(result, t)
	var resp nnr_moe_core.ResponseRuleDetail
	_ = json.Unmarshal([]byte(result), &resp)
	result = startCLIWithMockArgsThenCaptureResult([]string{
		"edit-rule",
		"--rid",
		resp.Data.Rid,
		"--remote",
		"12.34.56.78",
		"--rport",
		"443",
		"--type",
		"udp",
		"--name",
		"cli_testing",
	})
	checkStatus(result, t)
	if !strings.Contains(result, "12.34.56.78") {
		t.Error("Edit rule failed.")
	}
}

func TestDeleteRule(t *testing.T) {
	result := addRule()
	checkStatus(result, t)

	result = rules()
	var resp nnr_moe_core.ResponseRules
	_ = json.Unmarshal([]byte(result), &resp)
	for _, rule := range resp.Data {
		if !strings.Contains(rule.Name, cliTestingRuleName) {
			continue
		}
		// delete all cli testing rule
		result = startCLIWithMockArgsThenCaptureResult([]string{
			"delete-rule",
			"--rid",
			rule.Rid,
		})
		checkStatus(result, t)
	}
}

func TestGetRule(t *testing.T) {
	result := addRule()
	checkStatus(result, t)
	var resp nnr_moe_core.ResponseRuleDetail
	_ = json.Unmarshal([]byte(result), &resp)
	result = startCLIWithMockArgsThenCaptureResult([]string{
		"get-rule",
		"--rid",
		resp.Data.Rid,
	})
	checkStatus(result, t)
}

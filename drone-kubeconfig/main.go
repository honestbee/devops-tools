package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

var build = "0" // build number set at compile-time

// Config to hold Command configuration
type Config struct {
	Repository     string        `json:"repo_name"`
	ServiceAccount string        `json:"token"`
	Namespace      string        `json:"namespace"`
	Timeout        time.Duration `json:"organization"`
	KubeContexts   []string      `json:"kube_contexts"`
	CommandContext context.Context
}

var cfg = new(Config)

func main() {
	app := cli.NewApp()
	app.Name = "drone-kubeconfig"
	app.Version = fmt.Sprintf("0.1.%s", build)
	app.Usage = "create drone secrets for kubernetes service accounts"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "repository, r",
			Usage: "repository name (e.g. octocat/hello-world)",
		},
		cli.StringSliceFlag{
			Name:  "context, c",
			Usage: "map `PREFIX` with equivalent `KUBE_CONTEXT` (e.g STAGING_1A=<kube-content-name>)",
		},
		cli.StringFlag{
			Name:  "service-account, s",
			Usage: "Kubernetes service account for drone",
			Value: "drone",
		},
		cli.StringFlag{
			Name:  "namespace, n",
			Usage: "Kubernetes namespace of service account",
			Value: "kube-system",
		},
		cli.DurationFlag{
			Name:  "timeout, t",
			Usage: "`DURATION` before commands are cancelled",
			Value: time.Second * 5,
		},
	}

	app.Action = run

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	err := checkBinaries()
	if err != nil {
		return err
	}

	cfg = &Config{
		ServiceAccount: c.String("service-account"),
		Namespace:      c.String("namespace"),
		Timeout:        c.Duration("timeout"),
		Repository:     c.String("repository"),
		KubeContexts:   c.StringSlice("context"),
	}

	if cfg.Repository == "" {
		cfg.Repository = c.Args().First()
	}
	if cfg.Repository == "" {
		fmt.Println("missing repository name")
		cli.ShowAppHelpAndExit(c, 1)
	}

	ctxMap, err := parseContexts(cfg.KubeContexts)
	if err != nil {
		return err
	}

	if len(ctxMap) < 1 {
		fmt.Println("missing <PREFIX>=<KUBE_CONTEXT>")
		cli.ShowAppHelpAndExit(c, 1)
	}

	cfg.CommandContext = context.Background()
	for prefix, kubeCtx := range ctxMap {
		fmt.Printf("Setting up %s (with prefix %s) ...\n", kubeCtx, prefix)
		cluster, err := readKubeConfig(cfg, fmt.Sprintf(`{.contexts[?(@.name=="%v")].context.cluster}`, kubeCtx))
		if err != nil {
			return err
		}
		apiServer, err := readKubeConfig(cfg, fmt.Sprintf(`{.clusters[?(@.name=="%v")].cluster.server}`, cluster))
		if err != nil {
			return err
		}
		// fmt.Println(apiServer)

		fmt.Printf("Retrieving token for service account: '%s' - namespace '%s' (with %v timeout) ...\n", cfg.ServiceAccount, cfg.Namespace, cfg.Timeout)
		command := []string{
			"get",
			"sa",
			cfg.ServiceAccount,
		}

		secretName, err := runKubeCommand(cfg, kubeCtx, command, `{.secrets[].name}`)
		if err != nil {
			return err
		}

		command = []string{
			"get",
			"secret",
			secretName,
		}

		tokenb64, err := runKubeCommand(cfg, kubeCtx, command, "{.data.token}")
		if err != nil {
			return err
		}

		token, err := base64.StdEncoding.DecodeString(tokenb64)
		if err != nil {
			return err
		}

		// fmt.Println(string(token))
		fmt.Printf("Adding secrets for repoName '%s' (with %v timeout) ...\n", cfg.Repository, cfg.Timeout)

		_, err = addDroneSecret(cfg, prefix, "API_SERVER", apiServer)
		if err != nil {
			return err
		}
		_, err = addDroneSecret(cfg, prefix, "KUBERNETES_TOKEN", string(token))
		if err != nil {
			return err
		}
	}

	return nil
}

var contextsExp = regexp.MustCompile(`^([\w-]+)=(.*)`)

// parse prefix=context map string array into a golang map
func parseContexts(contexts []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, kvp := range contexts {
		matches := contextsExp.FindStringSubmatch(kvp)
		if len(matches) != 3 {
			return nil, fmt.Errorf("Invalid context mapping: %s", kvp)
		}
		// matches[0] = entire matched string, 1 is first capture group, ...
		result[matches[1]] = matches[2]
	}
	return result, nil
}

func readKubeConfig(cfg *Config, jsonpath string) (string, error) {
	readConfig := []string{
		"config",
		"view",
		"--output",
		fmt.Sprintf("jsonpath=\"%s\"", jsonpath),
	}

	return runCommand(cfg, "kubectl", readConfig)
}

func addDroneSecret(cfg *Config, prefix string, name string, value string) (string, error) {
	command := []string{
		"secret",
		"add",
		"--repository",
		cfg.Repository,
		"--name",
		fmt.Sprintf("%s_%s", prefix, name),
		"--value",
		value,
	}
	// fmt.Println(command)

	return runCommand(cfg, "drone", command)
}

func runKubeCommand(cfg *Config, kubeCtx string, command []string, jsonpath string) (string, error) {
	command = append(command, "--context", kubeCtx, "--namespace", cfg.Namespace)
	command = append(command, "--output", fmt.Sprintf("jsonpath='%s'", jsonpath))

	return runCommand(cfg, "kubectl", command)
}

// Run binary in background with timeout, return unquoted output string
func runCommand(cfg *Config, binary string, params []string) (string, error) {
	ctx, _ := context.WithTimeout(cfg.CommandContext, cfg.Timeout)
	cmd := exec.CommandContext(ctx, binary, params...)
	// fmt.Println(params)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output() //runs the command and returns the output

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Command timed out")
		return "", err
	}

	return trimQuotes(string(out)), err
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		switch {
		case s[0] == '"' && s[len(s)-1] == '"':
			return s[1 : len(s)-1]
		case s[0] == '\'' && s[len(s)-1] == '\'':
			return s[1 : len(s)-1]
		}
	}
	return s
}

func checkBinaries() error {
	bins := []string{
		"kubectl",
		"drone",
	}

	for _, b := range bins {
		_, err := exec.LookPath(b)
		if err != nil {
			return err
		}
	}

	return nil
}

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

func main() {
	app := cli.NewApp()
	app.Name = "drone-kfg"
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
			Usage: "`PREFIX=KUBE_CONTEXT` pairs for retrieving drone secrets",
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
		// cli.StringSliceFlag{
		// 	Name:  "event",
		// 	Usage: "secret limited to these events",
		// },
		// cli.StringSliceFlag{
		// 	Name:  "image",
		// 	Usage: "secret limited to these images",
		// },
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

	serviceAccount := c.String("service-account")
	namespace := c.String("namespace")
	killIn := c.Duration("timeout")

	repoName := c.String("repository")
	if repoName == "" {
		repoName = c.Args().First()
	}
	if repoName == "" {
		fmt.Println("missing repository name")
		cli.ShowAppHelpAndExit(c, 1)
	}

	ctxMap, err := parseContexts(c.StringSlice("context"))
	if err != nil {
		return err
	}

	if len(ctxMap) < 1 {
		fmt.Println("missing PREFIX=KUBE_CONTEXT pairs")
		cli.ShowAppHelpAndExit(c, 1)
	}

	ctx := context.Background()
	for prefix, kubeCtx := range ctxMap {
		fmt.Printf("Setting up %s (with prefix %s) ...\n", kubeCtx, prefix)
		cluster, err := readKubeConfig(ctx, killIn, fmt.Sprintf(`{.contexts[?(@.name=="%v")].context.cluster}`, kubeCtx))
		if err != nil {
			return err
		}
		apiServer, err := readKubeConfig(ctx, killIn, fmt.Sprintf(`{.clusters[?(@.name=="%v")].cluster.server}`, cluster))
		if err != nil {
			return err
		}
		// fmt.Println(apiServer)


		fmt.Printf("Retrieving token for service account: '%s' - namespace '%s' (with %v timeout) ...\n", serviceAccount, namespace, killIn)
		command := []string{
			"get",
			"sa",
			serviceAccount,
		}

		secretName, err := runKubeCommand(ctx, killIn, kubeCtx, namespace, command, `{.secrets[].name}`)
		if err != nil {
			return err
		}

		command = []string{
			"get",
			"secret",
			secretName,
		}

		tokenb64, err := runKubeCommand(ctx, killIn, kubeCtx, namespace, command, "{.data.token}")
		if err != nil {
			return err
		}

		token, err := base64.StdEncoding.DecodeString(tokenb64)
		if err != nil {
			return err
		}

		// fmt.Println(string(token))
		fmt.Printf("Adding secrets for repoName '%s' (with %v timeout) ...\n", repoName, killIn)
		
		_, err = addDroneSecret(ctx, killIn, repoName, prefix, "API_SERVER", apiServer)
		if err != nil {
			return err
		}
		_, err = addDroneSecret(ctx, killIn, repoName, prefix, "KUBERNETES_TOKEN", string(token))
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

func readKubeConfig(ctx context.Context, killIn time.Duration,  jsonpath string) (string, error) {
	readConfig := []string{
		"config",
		"view",
		"--output",
		fmt.Sprintf("jsonpath=\"%s\"", jsonpath),
	}

	return runCommand(ctx, "kubectl", killIn, readConfig)
}

func addDroneSecret(ctx context.Context, killIn time.Duration,  repoName string, prefix string, name string, value string) (string, error) {
	command := []string{
		"secret",
		"add",
		"--repository",
		repoName,
		"--name",
		fmt.Sprintf("%s_%s", prefix, name),
		"--value",
		value,
	}
	// fmt.Println(command)

	return runDroneCommand(ctx, killIn, command)
}

func runKubeCommand(ctx context.Context, killIn time.Duration,  kubeCtx string, namespace string, command []string, jsonpath string) (string, error) {
	command = append(command, "--context", kubeCtx, "--namespace", namespace)
	command = append(command, "--output", fmt.Sprintf("jsonpath='%s'", jsonpath))

	return runCommand(ctx, "kubectl", killIn, command)
}

func runDroneCommand(ctx context.Context, killIn time.Duration, command []string) (string, error) {
	return runCommand(ctx, "drone", killIn, command)
}

// Run binary in background with timeout, return unquoted output string
func runCommand(ctx context.Context, binary string, killIn time.Duration, params []string) (string, error) {
	ctx, _ = context.WithTimeout(ctx, killIn)
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
        if s[0] == '"' && s[len(s)-1] == '"' {
            return s[1 : len(s)-1]
        } else if s[0] == '\'' && s[len(s)-1] == '\'' {
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

	for _, b := range bins{
		_ , err := exec.LookPath(b)
		if err != nil {
			return err
		}
	}

	return nil
}

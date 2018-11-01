package kops

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const manifestPathTpl = "%s/%s/kops_cluster.yml"

type config struct {
	clusterName, kopsStatePath, kopsManifestPath, sshKey string
}

func deleteCluster(conf *config) (string, error) {
	out, err := exec.Command("kops", "delete", "cluster", conf.clusterName,
		"--state", conf.kopsStatePath,
		"--yes").CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}

func updateCluster(conf *config) (string, error) {
	out := getCluster(conf)
	if strings.Contains(out, "cluster not found") == true {
		if err := create(conf); err != nil {
			return out, err
		}
	}

	if err := update(conf); err != nil {
		return string(out), err
	}
	return string(out), nil
}
func getCluster(conf *config) string {
	log.Printf("[%s]Getting cluster with statePath %s", conf.clusterName, conf.kopsStatePath)
	out, _ := exec.Command("kops", "get", "cluster", conf.clusterName,
		"--state", conf.kopsStatePath,
		"-o", "yaml").CombinedOutput()

	return string(out)
}

func create(conf *config) error {
	log.Printf("[%s]Cluster not found, creating kops state...", conf.clusterName)
	out, err := exec.Command("kops", "create",
		"-f", fmt.Sprintf(manifestPathTpl, conf.kopsManifestPath, conf.clusterName),
		"--state", conf.kopsStatePath).CombinedOutput()
	if err != nil {
		log.Printf("%s\n", out)
		return err
	}

	log.Printf("[%s]Cluster state created", conf.clusterName)
	log.Printf("[%s]Creating kops ssh public key...", conf.clusterName)
	out, err = exec.Command("kops", "create", "secret",
		"--name", conf.clusterName,
		"sshpublickey", "admin",
		"-i", conf.sshKey,
		"--state", conf.kopsStatePath).CombinedOutput()
	if err != nil {
		log.Printf("%s\n", out)
		return err
	}

	return nil
}

func update(conf *config) error {

	log.Printf("[%s]Cluster found, updating kops state...", conf.clusterName)
	out, err := exec.Command("kops", "replace", "--force",
		"-f", fmt.Sprintf(manifestPathTpl, conf.kopsManifestPath, conf.clusterName),
		"--state", conf.kopsStatePath).CombinedOutput()
	if err != nil {
		log.Printf("%s\n", out)
		return err
	}

	log.Printf("[%s]Cluster state updated", conf.clusterName)

	log.Printf("[%s]Deploying cloud resources...", conf.clusterName)
	out, err = exec.Command("kops", "update", "cluster", conf.clusterName,
		"--state", conf.kopsStatePath, "--yes").CombinedOutput()
	if err != nil {
		log.Printf("%s\n", out)
		return err
	}
	return nil
}

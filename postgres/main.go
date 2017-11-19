/*TODO
- [ x ] dump function
- [ x ] restore function
- [ x ] help menu
- [   ] loop through list of info and dump
- [   ] loop through list of info and restore
*/
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func dump(dbHost string, dbName string, dbUser string) error {
	cmd := exec.Command("pg_dump", "-U", dbUser, "--format=c", "-f", dbName+".sqlc", "-d", dbName, "-O", "-x", "-h", dbHost)
	cmd.Stdin = strings.NewReader("password: ")
	var out bytes.Buffer
	var stderr bytes.Buffer // https://stackoverflow.com/a/18159705/2490986
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())

	}
	return err
}

func restore(dbHost string, dbName string, dbUser string) error {
	cmd := exec.Command("pg_restore", "-U", dbUser, "-d", dbName, dbName+".sqlc", "-h", dbHost)
	cmd.Stdin = strings.NewReader("password: ")
	var out bytes.Buffer
	var stderr bytes.Buffer // https://stackoverflow.com/a/18159705/2490986
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())

	}
	return err
}

func main() {
	dbUser, dbName, dbHost := os.Getenv("DBUSER"), os.Getenv("DBNAME"), os.Getenv("DBHOST")

	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "dbname",
			Usage:  "`<database name>` to dump",
			EnvVar: "DBNAME",
		},
		cli.StringFlag{
			Name:   "dbhost",
			Usage:  "`<database host>` to connect to",
			EnvVar: "DBHOST",
		},
		cli.StringFlag{
			Name:   "dbuser",
			Usage:  "`<username>` to authenticate with",
			EnvVar: "DBUSER",
		},
	}

	app := cli.NewApp()
	app.Name = "backstore"
	app.Usage = "posgres database backup & restore"
	app.Version = "0.1.0"
	app.Author = "Tuan Nguyen"

	app.Commands = []cli.Command{
		{
			Name:    "dump",
			Aliases: []string{"d"},
			Usage:   "dump database",
			Flags:   flags,
			Action: func(c *cli.Context) error {
				err := dump(dbHost, dbName, dbUser)
				return err
			},
		},
		{
			Name:    "restore",
			Aliases: []string{"r"},
			Usage:   "restore database",
			Action: func(c *cli.Context) error {
				err := restore(dbHost, dbName, dbUser)
				return err
			},
		},
	}

	app.Flags = flags

	app.Run(os.Args)
	// err := restore(dbHost, dbName, dbUser)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

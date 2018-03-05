/*TODO
- [ x ] dump function
- [ x ] restore function
- [ x ] help menu
- [ x ] Password via env
- [ x ] Password via csv file
- [ x ] Loop through list of info and dump
- [ x ] Loop through list of info and restore
- [   ] Check csv file format before read
- [   ] Add debug mode
*/
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"encoding/csv"

	"github.com/urfave/cli"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Use variadic function to overcome missing method-overloading in golang
func dump(dbinfo ...string) error {
	cmd := exec.Command("pg_dump", "--dbname=postgresql://"+dbinfo[2]+":"+dbinfo[3]+"@"+dbinfo[0]+":5432/"+dbinfo[1],
		"--format=c",
		"-f", dbinfo[1]+".sqlc",
		"-O", "-x")
	if dbinfo[3] == "" {
		cmd.Stdin = strings.NewReader("password: ")
	}
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

func restore(dbinfo ...string) error {
	cmd := exec.Command("pg_restore", "--dbname=postgresql://"+dbinfo[2]+":"+dbinfo[3]+"@"+dbinfo[0]+":5432/"+dbinfo[1],
		dbinfo[1]+".sqlc")
	if dbinfo[3] == "" {
		cmd.Stdin = strings.NewReader("password: ")
	}
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

func csvReader(filePath string) [][]string {
	file, err := os.Open(filePath)
	checkErr(err)

	defer file.Close()
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	return records
}

func main() {

	/*
	   The precedence for flag value sources is as follows (highest to lowest):
	   - Command line flag value from user
	   - Environment variable (if specified)
	   - Configuration file (if specified)
	   - Default defined on the flag
	*/

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
		cli.StringFlag{
			Name:   "dbpassword",
			Value:  "",
			Usage:  "`<password>` to authenticate with (optional)",
			EnvVar: "DBPASSWORD",
		},
		cli.StringFlag{
			Name:   "config",
			Value:  "",
			Usage:  "Load config from `FILE`",
			EnvVar: "DBCONFIG",
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
			Action: func(c *cli.Context) (err error) {
				if c.String("config") != "" {
					records := csvReader(c.String("config"))
					for _, record := range records {
						dump(record[0], record[1], record[2], record[3])
					}
				} else {
					dump(c.String("dbhost"), c.String("dbname"), c.String("dbuser"), c.String("dbpassword"))
				}
				return
			},
		},
		{
			Name:    "restore",
			Aliases: []string{"r"},
			Usage:   "restore database",
			Flags:   flags,
			Action: func(c *cli.Context) (err error) {
				if c.String("config") != "" {
					records := csvReader(c.String("config"))
					for _, record := range records {
						restore(record[0], record[1], record[2], record[3])
					}
				} else {
					restore(c.String("dbhost"), c.String("dbname"), c.String("dbuser"), c.String("dbpassword"))
				}

				return
			},
		},
	}

	app.Flags = flags

	app.Run(os.Args)
}

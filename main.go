package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"

	"gopkg.in/ini.v1"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configFilePath := filepath.Join(userHomeDir, ".aws", "config")
	cfg, err := ini.Load(configFilePath)
	if err != nil {
		return err
	}

	profiles := []Profile{}

	for _, section := range cfg.Sections() {
		if section.Name() == "DEFAULT" {
			continue
		}

		p := Profile{
			Name: section.Name()[8:],
		}

		if section.HasKey("sso_account_id") {
			p.AccountNumber = section.Key("sso_account_id").String()
		}

		profiles = append(profiles, p)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	for _, p := range profiles {
		last4 := ""
		if len(p.AccountNumber) > 4 {
			last4 = p.AccountNumber[8:]
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", last4, p.AccountNumber, p.Name)
	}

	w.Flush()

	return nil

}

type Profile struct {
	Name          string
	AccountNumber string
}

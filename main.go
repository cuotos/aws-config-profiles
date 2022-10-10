package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"gopkg.in/ini.v1"
)

var (
	appVersion string = "not provided"
	fullOutput        = false
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	if len(os.Args) == 2 {
		switch arg := os.Args[1]; arg {
		case "-v":
			fmt.Println(appVersion)
			os.Exit(0)
		case "-l":
			fullOutput = true
		default:
			fmt.Printf("unknown flag %s\n", arg)
			os.Exit(1)
		}
	}

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
		if section.Name() == "DEFAULT" || !strings.HasPrefix(section.Name(), "profile ") {
			continue
		}

		p := Profile{
			Name: strings.TrimPrefix(section.Name(), "profile "),
		}

		if section.HasKey("sso_account_id") {
			p.AccountNumber = section.Key("sso_account_id").String()
			p.AccountLastFour = p.AccountNumber[8:]
		}

		if section.HasKey("aws_access_key_id") {
			p.AccessKeyId = section.Key("aws_access_key_id").String()
		}

		profiles = append(profiles, p)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})

	for _, p := range profiles {
		write(w, p, fullOutput)
	}

	w.Flush()

	return nil
}

func write(w io.Writer, profile Profile, long bool) {
	if long {
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", profile.AccountNumber, profile.AccessKeyId, profile.Name)
	} else {
		fmt.Fprintf(w, "%s\t%s\t\n", profile.AccountLastFour, profile.Name)
	}
}

type Profile struct {
	Name            string
	AccountNumber   string
	AccountLastFour string
	AccessKeyId     string
}

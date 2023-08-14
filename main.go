package main

import (
	"flag"
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
	version = "unset"
	commit  = "unset"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	fullOutput := flag.Bool("l", false, "show full account id")
	printVersion := flag.Bool("v", false, "show version")
	flag.Parse()

	if *printVersion {
		fmt.Printf("%s-%s", version, commit)
		os.Exit(0)
	}

	outputProfiles := []Profile{}

	allProfiles, err := getProfilesFromAWSConfig()
	if err != nil {
		return err
	}

	// if user has provided a search string
	if flag.Arg(0) != "" {
		searchString := flag.Arg(0)

		for _, p := range allProfiles {
			if strings.Contains(p.AccountNumber, searchString) || strings.Contains(p.Name, searchString) {
				outputProfiles = append(outputProfiles, p)
			}
		}

	} else {
		outputProfiles = allProfiles
	}

	printOutput(outputProfiles, *fullOutput)

	return nil
}

func printOutput(profiles []Profile, full bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})

	for _, p := range profiles {
		write(w, p, full)
	}

	w.Flush()
}

func getProfilesFromAWSConfig() ([]Profile, error) {

	profiles := []Profile{}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return profiles, err
	}
	configFilePath := filepath.Join(userHomeDir, ".aws", "config")
	cfg, err := ini.Load(configFilePath)
	if err != nil {
		return profiles, err
	}

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

	return profiles, nil
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

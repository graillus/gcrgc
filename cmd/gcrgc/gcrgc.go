package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/graillus/gcrgc/pkg/gcrgc"
)

// ParseArgs parses compmand-line args and returns a Settings instance
func ParseArgs() *gcrgc.Settings {
	settings := gcrgc.Settings{}

	flag.StringVar(&settings.Registry, "registry", "", "Google Cloud Registry name, e.g. \"gcr.io/project-id\". Some repositories can be excluded with -exclude-repository option")
	flag.StringVar(&settings.Date, "date", "", "Delete images older than YYYY-MM-DD")
	flag.BoolVar(&settings.AllRepositories, "all", false, "Include all repositories from the registry. Defaults to false.")
	flag.BoolVar(&settings.UntaggedOnly, "untagged-only", false, "Only remove untagged images. Defaults to false.")
	flag.BoolVar(&settings.DryRun, "dry-run", false, "See images to be deleted without actually deleting them. Defaults to false.")
	flag.Var(&settings.ExcludedRepositories, "exclude-repository", "Repo(s) to be excluded, to be used in addition with the -all option. Can be repeated.")
	flag.Var(&settings.ExcludedTags, "exclude-tag", "Tag(s) to be excluded. Can be repeated.")

	flag.Parse()

	args := flag.Args()
	settings.Repositories = args

	if settings.Registry == "" {
		fmt.Println("The -registry option is missing")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if settings.AllRepositories == false && len(settings.Repositories) == 0 {
		fmt.Println("You must provide at least one repository name, or set the -all option to include all repositories from the registry")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if settings.AllRepositories == false && len(settings.ExcludedRepositories) > 0 {
		fmt.Println("You cannot exclude repositories unless using option -all")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if settings.UntaggedOnly == true && len(settings.ExcludedTags) > 0 {
		fmt.Println("You cannot exclude tags when option -untagged-only is true")
		flag.PrintDefaults()
		os.Exit(1)
	}

	return &settings
}

func main() {
	settings := ParseArgs()

	app := gcrgc.NewApp(settings)
	app.Start()
}

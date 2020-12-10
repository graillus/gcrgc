package main

import (
	"fmt"
	"os"
	"time"

	"github.com/k1LoW/duration"
	flag "github.com/spf13/pflag"

	"github.com/graillus/gcrgc/internal/gcrgc"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("gcrgc [options] <registry>")
	fmt.Println("")
	fmt.Println("The registry argument should have the form gcr.io/<project-id>")
	fmt.Println("")
	fmt.Println("List of possible options:")
	flag.PrintDefaults()
}

// ParseArgs parses compmand-line args and returns a Settings instance
func ParseArgs() *gcrgc.Settings {
	var (
		help            bool
		limitDate       string
		retentionPeriod string
	)

	flag.BoolVar(&help, "help", false, "Print the command usage")

	settings := gcrgc.Settings{}

	flag.StringSliceVar(&settings.Repositories, "repositories", []string{}, "A comma-separated list of repositories to include in the cleanup process.")
	flag.StringVar(&limitDate, "date", "", "Only delete images older than YYYY-MM-DD")
	flag.StringVar(&retentionPeriod, "retention-period", "", "The retention period: only older items will be deleted. Must have the form: `30 days`, `1w`, `24h`")
	flag.BoolVar(&settings.UntaggedOnly, "untagged-only", false, "Only remove untagged images. Defaults to false.")
	flag.BoolVar(&settings.DryRun, "dry-run", false, "Only see the output of what would be deleted but don't actually delete anything. Defaults to false.")
	flag.StringSliceVar(&settings.ExcludedRepositories, "exclude-repositories", []string{}, "A comma-separated list of repositories to be excluded. Not compatible with the --repositories option.")
	flag.StringSliceVar(&settings.ExcludedTags, "exclude-tags", []string{}, "A comma-separated list of tag(s) to be excluded.")
	flag.StringArrayVar(&settings.ExcludedTagPatterns, "exclude-tag-pattern", []string{}, "Tag patterns(s) to be excluded. Repeat the option to provide many.")
	flag.BoolVar(&settings.ExcludeSemVerTags, "exclude-semver-tags", false, "Only remove images not tagged with a SemVer tag. Defaults to false.")

	flag.Parse()

	if help == true {
		printUsage()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Printf("Error: The \"registry\" argument was not provided\n\n")
		printUsage()
		os.Exit(1)
	}
	settings.Registry = args[0]

	if limitDate != "" {
		date, err := time.Parse("2006-01-02", limitDate)
		if err != nil {
			fmt.Printf("Error: Unable to parse the --date flag: invalid date format: %s\n", limitDate)
			os.Exit(1)
		}

		settings.Date = &date
	}

	if retentionPeriod != "" && settings.Date != nil {
		fmt.Println("Error: The --date and --retention-period flags are not compatible. Ony one of them can be provided")
		os.Exit(1)
	}

	if retentionPeriod != "" {
		parsedDuration, err := duration.Parse(retentionPeriod)
		if err != nil {
			fmt.Println("Unable to parse the --retention-period flag. Run with the --help flag for more information")
			os.Exit(1)
		}
		date := time.Now().Add(-parsedDuration)
		settings.Date = &date
	}

	if len(settings.Repositories) == 0 {
		settings.AllRepositories = true
	}

	if len(settings.Repositories) > 0 && len(settings.ExcludedRepositories) > 0 {
		fmt.Println("Error: The --repositories and the --exclude-repositories flags cannot be provided altogether")
		os.Exit(1)
	}

	if settings.UntaggedOnly == true && len(settings.ExcludedTags) > 0 {
		fmt.Println("Error: The --exclude-tags and the --untagged-only flags cannot be provided altogether")
		os.Exit(1)
	}

	if settings.UntaggedOnly == true && len(settings.ExcludedTagPatterns) > 0 {
		fmt.Println("Error: The --exclude-tag-pattern and the --untagged-only flags cannot be provided altogether")
		os.Exit(1)
	}

	if settings.UntaggedOnly == true && settings.ExcludeSemVerTags == true {
		fmt.Println("Error: The --exclude-semver-tags and the --untagged-only flags cannot be provided altogether")
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

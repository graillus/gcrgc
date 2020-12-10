package main

import (
	"fmt"
	"os"
	"time"

	"github.com/k1LoW/duration"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/graillus/gcrgc/internal/gcrgc"
)

func main() {
	settings := ParseArgs()

	app := gcrgc.NewApp(settings)
	app.Start()
}

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
		help bool
		file string
	)

	// Cli config
	flag.BoolVarP(&help, "help", "p", false, "Print the command usage")
	flag.StringVarP(&file, "config", "c", "", "Path to the configuration file.")

	// App settings
	flag.StringSlice("repositories", []string{}, "A comma-separated list of repositories to include in the cleanup process.")
	flag.String("date", "", "Only delete images older than YYYY-MM-DD")
	flag.String("retention-period", "", "The retention period: only older items will be deleted. Must have the form: `30 days`, `1w`, `24h`")
	flag.Bool("untagged-only", false, "Only remove untagged images. Defaults to false.")
	flag.Bool("dry-run", false, "Only see the output of what would be deleted but don't actually delete anything. Defaults to false.")
	flag.StringSlice("exclude-repositories", []string{}, "A comma-separated list of repositories to be excluded. Not compatible with the --repositories option.")
	flag.StringSlice("exclude-tags", []string{}, "A comma-separated list of tag(s) to be excluded.")
	flag.StringSlice("exclude-tag-pattern", []string{}, "Tag patterns(s) to be excluded. Repeat the option to provide many.")
	flag.Bool("exclude-semver-tags", false, "Only remove images not tagged with a SemVer tag. Defaults to false.")

	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	// Print help and exit if the --help flag was provided
	if help == true {
		printUsage()
		os.Exit(0)
	}

	// Parse the config file if the --config flag was provided
	if file != "" {
		viper.SetConfigFile(file)
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Printf("Error: Unable to load config file %s: %s\n", file, err.Error())
			os.Exit(1)
		}
	}

	// Retrieve the registry argument, or exit if missing
	args := flag.Args()
	if len(args) < 1 && file == "" {
		fmt.Printf("Error: The \"registry\" argument is missing\n\n")
		printUsage()
		os.Exit(1)
	}

	// If no file provided, fallback to registry argument
	if file == "" {
		viper.Set("registry", args[0])
	}

	settings := gcrgc.Settings{
		Registry:             viper.GetString("registry"),
		Repositories:         viper.GetStringSlice("repositories"),
		ExcludedRepositories: viper.GetStringSlice("exclude-repositories"),
		UntaggedOnly:         viper.GetBool("untagged-only"),
		ExcludeSemVerTags:    viper.GetBool("exclude-semver-tags"),
		ExcludedTags:         viper.GetStringSlice("exclude-tags"),
		ExcludedTagPatterns:  viper.GetStringSlice("exclude-tag-pattern"),
		Date:                 viper.GetTime("date"),
		DryRun:               viper.GetBool("dry-run"),
	}

	retentionPeriod := viper.GetString("retention-period")
	if retentionPeriod != "" && settings.Date != (time.Time{}) {
		fmt.Println("Error: The \"date\" and \"retention-period\" settings conflict with each other, only one must be set.")
		os.Exit(1)
	}

	if retentionPeriod != "" {
		parsedDuration, err := duration.Parse(retentionPeriod)
		if err != nil {
			fmt.Println("Error: Unable to parse the \"retention-period\" option.")
			os.Exit(1)
		}
		settings.Date = time.Now().Add(-parsedDuration)
	}

	if len(settings.Repositories) == 0 {
		settings.AllRepositories = true
	}

	if len(settings.Repositories) > 0 && len(settings.ExcludedRepositories) > 0 {
		fmt.Println("Error: The \"repositories\" and the \"exclude-repositories\" settings conflict with each other, only one must be set.")
		os.Exit(1)
	}

	if settings.UntaggedOnly == true && len(settings.ExcludedTags) > 0 {
		fmt.Println("Error: The \"exclude-tags\" and the \"untagged-only\" settings conflict with each other, only one must be set.")
		os.Exit(1)
	}

	if settings.UntaggedOnly == true && len(settings.ExcludedTagPatterns) > 0 {
		fmt.Println("Error: The \"exclude-tag-pattern\" and the \"untagged-only\" settings conflict with each other, only one must be set.")
		os.Exit(1)
	}

	if settings.UntaggedOnly == true && settings.ExcludeSemVerTags == true {
		fmt.Println("Error: The \"exclude-semver-tags\" and the \"untagged-only\" settings conflict with each other, only one must be set.")
		os.Exit(1)
	}

	return &settings
}

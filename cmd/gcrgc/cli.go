package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Settings contains app Configuration
type Settings struct {
	Repository     string
	Images         stringList
	Date           string
	UntaggedOnly   bool
	DryRun         bool
	AllImages      bool
	ExcludedImages stringList
	ExcludedTags   stringList
}

type stringList []string

func (s *stringList) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringList) Set(value string) error {
	*s = append(*s, value)

	return nil
}

// ParseArgs parses compmand-line args and returns a Settings instance
func ParseArgs() *Settings {
	settings := Settings{}

	flag.StringVar(&settings.Repository, "repository", "", "Clean all images from a given repository, e.g. \"gcr.io/project-id\". Some images can be excluded with -exclude-image option")
	flag.StringVar(&settings.Date, "date", "", "Delete images older than YYYY-MM-DD")
	flag.BoolVar(&settings.AllImages, "all", false, "Include all images from the repository. Defaults to false.")
	flag.BoolVar(&settings.UntaggedOnly, "untagged-only", false, "Only remove untagged images. Defaults to false.")
	flag.BoolVar(&settings.DryRun, "dry-run", false, "See images to be deleted without actually deleting them. Defaults to false.")
	flag.Var(&settings.ExcludedImages, "exclude-image", "Image(s) to be excluded, to be used in addition with the -all option. Can be repeated.")
	flag.Var(&settings.ExcludedTags, "exclude-tag", "Tag(s) to be excluded. Can be repeated.")

	flag.Parse()

	args := flag.Args()
	settings.Images = args

	if settings.Repository == "" {
		fmt.Println("The -repository option is missing")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if settings.AllImages == false && len(settings.Images) == 0 {
		fmt.Println("You must provide at least one image name, or set the -all option to include all images from the repository")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if settings.AllImages == false && len(settings.ExcludedImages) > 0 {
		fmt.Println("You cannot exclude images unless using option -all")
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

package main

import (
	registry "docker-registry-cleaner/docker-registry-client"
	"fmt"
	"os"
	"regexp"
	"sort"
	"time"

	version "github.com/hashicorp/go-version"
	"github.com/urfave/cli"
)

type registryParams struct {
	URL      string
	username string
	password string
}

func main() {

	app := cli.NewApp()
	app.Name = "Docker Registry Cleaner"
	app.Version = "0.1.0"
	app.Compiled = time.Now()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "url",
			Usage:  "Registry url",
			EnvVar: "URL",
		},
		cli.StringFlag{
			Name:   "username, u",
			Usage:  "Registry username (optional)",
			EnvVar: "USERNAME",
		},
		cli.StringFlag{
			Name:   "password, p",
			Usage:  "Registry password (optional)",
			EnvVar: "PASSWORD",
		},
		cli.StringFlag{
			Name:   "image, i",
			Usage:  "Image name to delete ie 'development/nginx'",
			EnvVar: "IMAGE",
		},
		cli.StringFlag{
			Name:   "imageversion, iv",
			Value:  ".*-SNAPSHOT.*",
			Usage:  "Image Version to delete, this can be a regex \".*-SNAPSHOT.*\"",
			EnvVar: "IMAGE_VERSION",
		},
		cli.IntFlag{
			Name:   "keep, k",
			Value:  3,
			Usage:  "The number of images you want to keep, usefully if you are deleting images by regex",
			EnvVar: "KEEP",
		},
		cli.BoolFlag{
			Name:   "dryrun, d",
			Usage:  "Do not actually delete anything",
			EnvVar: "DRYRUN",
		},
	}

	app.Action = func(c *cli.Context) error {

		r := registryParams{
			URL:      c.String("url"),
			username: c.String("username"),
			password: c.String("password"),
		}

		imageName := c.String("image")
		regxVersion := c.String("imageversion")
		numKeep := c.Int("keep")

		if r.URL == "" {
			return cli.ShowAppHelp(c)
		}

		if imageName == "" {
			return cli.ShowAppHelp(c)
		}

		hub, err := registry.New(r.URL, r.username, r.password)
		if err != nil {
			fmt.Printf("%s", err)
		}

		versionsRaw := []string{}

		tags, err := hub.Tags(imageName)

		// Get list of versions that match the image name
		for _, element := range tags {
			r := regexp.MustCompile(regxVersion)
			matches := r.FindString(element)
			if len(matches) != 0 {
				if matches == element {
					versionsRaw = append(versionsRaw, matches)
				}
			}
		}

		// Prep versions to be sorted
		versions := make([]*version.Version, len(versionsRaw))
		for i, raw := range versionsRaw {
			v, _ := version.NewVersion(raw)
			versions[i] = v
		}

		// After this, the versions are properly sorted from high to low
		sort.Sort(sort.Reverse(version.Collection(versions)))

		if c.Bool("dryrun") {
			fmt.Printf("\nDRY RUN - nothing will be deleted \n")
		}

		fmt.Printf("\nFound %d images that match, keeping the %d latest versions and deleting the rest \n", len(versions), numKeep)

		for i, v := range versions {
			//fmt.Printf("\nsorted version: %s\n", v)
			if i >= numKeep {
				fmt.Printf("\nDelete version: %s\n", v)

				if !c.Bool("dryrun") {
					// Get the manifest digest for the image
					digest, err := hub.ManifestDigest(imageName, v.String())
					fmt.Printf("Deleting Manifest Digest: %s\n", digest)
					// delete manifest
					err = hub.DeleteManifest(imageName, digest)
					if err != nil {
						fmt.Printf("%s", err)
					}
				}
			} else {
				fmt.Printf("\nKeep version: %s\n", v)
			}

		}
		return nil
	}

	app.Run(os.Args)
}

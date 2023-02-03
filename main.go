package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	sdk "github.com/grafana-tools/sdk"
)

func dashboardPush() {
	fmt.Println("Pushing dashboards to", os.Args[2], "\n")
	var (
		filesInDir []os.FileInfo
		rawBoard   []byte
		err        error
	)

	ctx := context.Background()
	client, err := sdk.NewClient(os.Args[2], os.Args[3], sdk.DefaultHTTPClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a client: %s\n", err)
		os.Exit(1)
	}
	filesInDir, err = ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawBoard, err = ioutil.ReadFile(file.Name()); err != nil {
				log.Println(err)
				continue
			}
			_, err := client.SetRawDashboard(ctx, rawBoard)
			if err != nil {
				log.Printf("error importing dashboard from %s", file.Name())
				continue
			}
		}
	}
}

func dashboardPull() {
	fmt.Println("Pulling dashboards from", os.Args[2], "\n")
	var (
		boardLinks []sdk.FoundBoard
		rawBoard   []byte
		meta       sdk.BoardProperties
		err        error
	)

	ctx := context.Background()
	client, err := sdk.NewClient(os.Args[2], os.Args[3], sdk.DefaultHTTPClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a client: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Client created for %s\n", os.Args[1])

	if boardLinks, err = client.SearchDashboards(ctx, "", false); err != nil {
		fmt.Fprint(os.Stderr, "Failed to search Dashboards:", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d dashboards!\n", len(boardLinks))
	// Print found Dashboards
	fmt.Printf("Found Dashboards: %v\n", boardLinks)

	for _, link := range boardLinks {
		fmt.Printf("Downloading %s (ID: %s)\n", link.Title, link.UID)
		if rawBoard, meta, err = client.GetRawDashboardByUID(ctx, link.UID); err != nil {
			fmt.Fprintf(os.Stderr, "%s for %s\n", err, link.UID)
			continue
		}
		if err = ioutil.WriteFile(fmt.Sprintf("%s.json", meta.Slug), rawBoard, os.FileMode(int(0777))); err != nil {
			fmt.Fprintf(os.Stderr, "%s for %s\n", err, meta.Slug)
		}
	}
}

func checkArgs() {
	if len(os.Args) != 4 {
		fmt.Fprint(os.Stderr, "Not all Args are set!\n\nUsage: (grafana-dashboard-sync) [push/pull] {GRAFANA_HOST} {API_KEY}\n")
		os.Exit(0)
	}
}

func main() {
	checkArgs()

	switch os.Args[1] {
	case "push":
		dashboardPush()
	case "pull":
		dashboardPull()
	default:
		fmt.Fprint(os.Stderr, "Usage: (grafana-dashboard-sync) [push/pull] {GRAFANA_HOST} {API_KEY}\n")
		os.Exit(0)
	}
}

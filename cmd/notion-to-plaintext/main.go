package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tmc/notion"
)

var (
	flagVerbose = flag.Bool("v", false, "verbose")
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "please provide block (page) id as parameter")
		os.Exit(1)
	}
	if err := run(flag.Args()[0]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(id string) error {
	opts := []notion.ClientOption{
		notion.WithToken(os.Getenv("NOTION_TOKEN")),
	}
	if *flagVerbose {
		opts = append(opts, notion.WithDebugLogging())
	}
	c, err := notion.NewClient(opts...)
	if err != nil {
		return err
	}
	blockInfo, err := c.GetRecordValues(notion.Record{Table: "block", ID: id})
	if err != nil {
		return err
	}
	if blockInfo[0].Value == nil {
		return fmt.Errorf("issue fetching content, Role=%v", blockInfo[0].Role)
	}
	p, err := c.GetBlock(blockInfo[0].Value.ID)
	if err != nil {
		return err
	}
	r, err := notion.PrintAsVim(p, "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(r))
	return nil
}

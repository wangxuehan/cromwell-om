package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"net/url"
	"os"
	"path"
	"regexp"
)

func work_summary(client *Client, operation string, failed bool, stderr bool, spendTime bool) {
	params := url.Values{}
	//params.Add("includeKey", "executionStatus")
	resp, _ := client.Metadata(operation, params)
	var mtr = MetadataTableResponse{Metadata: resp}
	tabbb := new(TableStruct)
	tabbb.Header = []string{"Call-Name", "Status"}
	if spendTime {
		tabbb.Header = append(tabbb.Header, "TimeConsuming")
	}
	if stderr {
		tabbb.Header = append(tabbb.Header, "Stderr")
	}
	tabbb.Data = [][]string{}
	//fmt.Println(mtr.Metadata.WorkflowRoot)
	for k, v := range mtr.Metadata.Calls {
		//fmt.Println(k)
		wf_name, _ := regexp.Compile("[^.]+\\.")
		real_name := wf_name.ReplaceAllString(k, "")
		for _, v1 := range v {
			//fmt.Println(k, v1.ExecutionStatus, v1.Stderr)
			if failed && v1.ExecutionStatus != "Failed" {
				continue
			}
			if v1.Stderr == "" {
				if v1.SubWorkflowID == "" {
					v1.Stderr = path.Join(mtr.Metadata.WorkflowRoot, "call-"+real_name, "shard-"+fmt.Sprintf("%d", v1.ShardIndex))
				} else {
					v1.Stderr = path.Join(mtr.Metadata.WorkflowRoot, "call-"+real_name,
						"shard-"+fmt.Sprintf("%d", v1.ShardIndex),
						real_name,
						v1.SubWorkflowID,
					)
				}
			}
			add_row := []string{real_name, v1.ExecutionStatus}
			if spendTime {
				add_row = append(add_row, fmt.Sprintf("%v", v1.End.Sub(v1.Start)))
			}
			if stderr {
				add_row = append(add_row, v1.Stderr)
			}
			tabbb.Data = append(tabbb.Data, add_row)
		}
	}
	tabbb.showTable()
}

func main() {
	client := Client{}

	var Version = "v1.0"
	app := &cli.App{
		Name:  "Cromwell-OM",
		Usage: "OM for Cromwell Server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "host",
				Value: "http://127.0.0.1:8000",
				Usage: "Url for your Cromwell Server",
			},
		},
		Before: func(c *cli.Context) error {
			client.Setup(c.String("host"))
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Cromwell-OM version",
				Action: func(c *cli.Context) error {
					fmt.Printf("Version: %s\nAuthor: wangxuehan\n", Version)
					return nil
				},
			},
			{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "Status of workflow, auto return failed info",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "workid", Aliases: []string{"i"}, Required: true, Usage: "WorkFlow ID"},
				},
				Action: func(c *cli.Context) error {
					sr, _ := client.Status(c.String("workid"))
					params := url.Values{}
					params.Add("includeKey", "failures")
					resp, _ := client.Metadata(c.String("operation"), params)
					var mtr = MetadataTableResponse{Metadata: resp}
					fmt.Printf("%s\t%s\n%+v\n", sr.ID, sr.Status, mtr.Metadata.Failures)
					return nil
				},
			},
			{
				Name:    "summary",
				Aliases: []string{"m"},
				Usage:   "Summary through metadata data (table)",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "workid", Aliases: []string{"i"}, Required: true, Usage: "WorkFlow ID"},
					&cli.BoolFlag{Name: "failed", Aliases: []string{"f"}, Usage: "Filter failed tasks"},
					&cli.BoolFlag{Name: "stderr", Aliases: []string{"e"}, Usage: "Show the stderr of the task"},
					&cli.BoolFlag{Name: "time", Aliases: []string{"t"}, Usage: "Show the time consuming of the task"},
				},
				Action: func(c *cli.Context) error {
					work_summary(&client, c.String("workid"), c.Bool("failed"),
						c.Bool("stderr"), c.Bool("time"),
					)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Error %q", err)
	}
}

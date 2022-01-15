package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"net/url"
	"os"
	"path"
	"regexp"
	"time"
)

func subMeta(client *Client, operation string, failed bool, stderr bool, spendTime bool, table *TableStruct,
	indent bool, subflow bool) {

	params := url.Values{}
	//params.Add("includeKey", "executionStatus")
	resp, _ := client.Metadata(operation, params)
	var mtr = MetadataTableResponse{Metadata: resp}

	for k, v := range mtr.Metadata.Calls {
		//fmt.Println(k)
		wf_name, _ := regexp.Compile("[^.]+\\.")
		real_name := wf_name.ReplaceAllString(k, "")

		subNum := 1
		if indent {
			space := ""
			for i := 0; i < subNum; i++ {
				space += " "
			}
			real_name = space + "\\_ " + real_name
		}
		for _, v1 := range v {
			//fmt.Println(k, v1.ExecutionStatus, v1.Stderr)
			if failed && v1.ExecutionStatus != "Failed" {
				continue
			}
			add_row := []string{real_name, v1.ExecutionStatus}

			if v1.Stderr == "" {
				v1.Stderr = path.Join(mtr.Metadata.WorkflowRoot, "call-"+real_name,
					"shard-"+fmt.Sprintf("%d", v1.ShardIndex))
			}

			if spendTime {
				if v1.End.IsZero() {
					v1.End = time.Now()
				}
				add_row = append(add_row, fmt.Sprintf("%v", v1.End.Sub(v1.Start)))
			}
			if stderr {
				add_row = append(add_row, v1.Stderr)
			}
			table.Data = append(table.Data, add_row)
			if v1.SubWorkflowID != "" {
				subNum += 1
				if subflow {
					subMeta(client, v1.SubWorkflowID, failed, stderr, spendTime, table, true, subflow)
				}
			} else {
				subNum -= 1
			}
		}
	}
}

func work_summary(client *Client, operation string, failed bool, stderr bool, spendTime bool, subflow bool) {

	table := new(TableStruct)
	table.Header = []string{"Call-Name", "Status"}
	if spendTime {
		table.Header = append(table.Header, "TimeConsuming")
	}
	if stderr {
		table.Header = append(table.Header, "Stderr")
	}
	table.Data = [][]string{}
	//fmt.Println(mtr.Metadata.WorkflowRoot)
	subMeta(client, operation, failed, stderr, spendTime, table, false, subflow)
	table.showTable()
}

func main() {
	client := Client{}

	var Version = "v1.1"
	app := &cli.App{
		Name:  "Cromwell-OM",
		Usage: "OM for Cromwell Server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "host",
				Value: "http://127.0.0.1:8000",
				Usage: "Url for your Cromwell Server. You can also write to cromwell-om.conf using json format",
			},
		},
		Before: func(c *cli.Context) error {
			conf := readConfig()
			host := conf.Host
			if conf.Host == "" || c.String("host") != "http://127.0.0.1:8000" {
				host = c.String("host")
			}
			client.Setup(host)
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
					resp, _ := client.Metadata(c.String("workid"), params)
					var mtr = MetadataTableResponse{Metadata: resp}
					failure := fmt.Sprintf("%+v", mtr.Metadata.Failures)
					if failure == "[]" {
						failure = ""
					} else {
						failure = "\n" + failure + "\n"
					}
					fmt.Printf("%s\t%s\n%s", sr.ID, sr.Status, failure)
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
					&cli.BoolFlag{Name: "subflow", Aliases: []string{"s"}, Usage: "Show the sub workflow"},
				},
				Action: func(c *cli.Context) error {
					work_summary(&client, c.String("workid"), c.Bool("failed"),
						c.Bool("stderr"), c.Bool("time"), c.Bool("subflow"),
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

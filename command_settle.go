package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type stackEvents []*cloudformation.StackEvent

func (a stackEvents) Len() int           { return len(a) }
func (a stackEvents) Less(i, j int) bool { return a[i].Timestamp.Before(*a[j].Timestamp) }
func (a stackEvents) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func handleSettle(cf *cloudformation.CloudFormation, stackName string, timeout, pollInterval time.Duration) error {
	firstTime := time.Now().Add(0 - (time.Minute))

	seen := make(map[string]bool)

	for {
		if time.Now().After(firstTime.Add(timeout)) {
			return fmt.Errorf("exceeded timeout of %s for applying changes", timeout)
		}

		r, err := cf.DescribeStackEvents(&cloudformation.DescribeStackEventsInput{
			StackName: aws.String(stackName),
		})

		if err != nil {
			return err
		}

		sort.Sort(stackEvents(r.StackEvents))

		for _, e := range r.StackEvents {
			if e.Timestamp.Before(firstTime) {
				continue
			}

			if !seen[*e.EventId] {
				seen[*e.EventId] = true

				fmt.Printf("%s %-20s %s\n", e.Timestamp.Local().Format(time.Kitchen), *e.ResourceStatus, *e.LogicalResourceId)
				if e.ResourceStatusReason != nil {
					fmt.Printf("    %s\n", *e.ResourceStatusReason)
				}
			}
		}

		if len(r.StackEvents) > 0 {
			lastEvent := r.StackEvents[len(r.StackEvents)-1]

			if *lastEvent.LogicalResourceId == stackName {
				switch {
				case strings.HasSuffix(*lastEvent.ResourceStatus, "_COMPLETE"):
					fmt.Printf("Stack has settled in a COMPLETE state\n")
					return nil
				case strings.HasSuffix(*lastEvent.ResourceStatus, "_FAILED"):
					fmt.Printf("Stack has settled in a FAILED state\n")
					return fmt.Errorf("failed: %s", *lastEvent.ResourceStatus)
				}
			}
		}

		time.Sleep(pollInterval)
	}
}

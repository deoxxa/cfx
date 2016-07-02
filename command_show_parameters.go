package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func handleShowParameters(cf *cloudformation.CloudFormation, stackName string) error {
	describeStackInput := cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	}

	describeStackOutput, err := cf.DescribeStacks(&describeStackInput)
	if err != nil {
		return err
	}

	if len(describeStackOutput.Stacks) != 1 {
		return fmt.Errorf("expected to find one stack; instead found %d", len(describeStackOutput.Stacks))
	}

	m := make(map[string]string)
	for _, p := range describeStackOutput.Stacks[0].Parameters {
		m[*p.ParameterKey] = *p.ParameterValue
	}

	d, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(d))

	return nil
}

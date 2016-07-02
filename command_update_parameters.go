package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func handleUpdateParameters(cf *cloudformation.CloudFormation, stackName string, capabilities []string, parameters map[string]string) error {
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

	for k, v := range parameters {
		fmt.Printf("# Changing parameter %q from %q to %q\n", k, m[k], v)

		m[k] = v
	}

	var plist []*cloudformation.Parameter
	for k, v := range m {
		plist = append(plist, &cloudformation.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		})
	}

	updateStackInput := cloudformation.UpdateStackInput{
		StackName:           aws.String(stackName),
		UsePreviousTemplate: aws.Bool(true),
		Capabilities:        aws.StringSlice(capabilities),
		Parameters:          plist,
	}

	if _, err := cf.UpdateStack(&updateStackInput); err != nil {
		return err
	}

	return nil
}

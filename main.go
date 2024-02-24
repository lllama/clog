package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func main() {
	fmt.Println("Running...")

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-2"))
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	client := cloudwatchlogs.NewFromConfig(config)

	paginator := cloudwatchlogs.NewDescribeLogGroupsPaginator(client, &cloudwatchlogs.DescribeLogGroupsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			log.Fatalf("failed to list log groups, %v", err)
		}
		for _, logGroup := range output.LogGroups {
			fmt.Println(*logGroup.LogGroupName)
		}
	}
}

// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/spf13/cobra"
)

var IP string
var Machine string
var Region = "us-east-1"

var sgCmd = &cobra.Command{
	Use:   "sg",
	Short: "SecurityGroup Management",
	Long:  `A SecurityGroup management command tool`,
	Run: func(cmd *cobra.Command, args []string) {
		session, err := session.NewSession()
		if err != nil {
			panic(err)
		}
		svc := ec2.New(session, &aws.Config{Region: aws.String(Region)})
		params := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("tag:Name"),
					Values: []*string{
						aws.String(Machine),
					},
				},
			},
		}
		response, err := svc.DescribeInstances(params)
		if err != nil {
			panic(err)
		}
		for _, inst := range response.Reservations[0].Instances {
			for _, securityGroup := range inst.SecurityGroups {

				ipExists := fetchIpsFromSecurityGroup(svc, securityGroup)
				fmt.Println(ipExists)
			}
		}
	},
}

func fetchIpsFromSecurityGroup(svc *ec2.EC2, securityGroup *ec2.GroupIdentifier) bool {
	paramsSG := &ec2.DescribeSecurityGroupsInput{
		DryRun:   aws.Bool(false),
		GroupIds: []*string{aws.String(*securityGroup.GroupId)},
	}

	resp, err := svc.DescribeSecurityGroups(paramsSG)

	if err != nil {
		panic(err)
	}
	ipExists := false
	for _, ip := range resp.SecurityGroups[0].IpPermissions {
		for _, ipRange := range ip.IpRanges {
			ipExists = strings.Contains(*ipRange.CidrIp, IP)
		}
	}
	return ipExists
}

func init() {
	RootCmd.AddCommand(sgCmd)

	sgCmd.PersistentFlags().StringVarP(&IP, "ip", "i", "", "IPAddress")
	sgCmd.PersistentFlags().StringVarP(&Machine, "machine", "m", "", "Machine")
}

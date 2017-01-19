package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Starts a new AWS session
func getDataFromMachine(machineName string) (*ec2.DescribeInstancesOutput, *ec2.EC2) {
	session, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	svc := ec2.New(session, &aws.Config{Region: aws.String("us-east-1")})
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(machineName),
				},
			},
		},
	}
	response, err := svc.DescribeInstances(params)
	if err != nil {
		panic(err)
	}
	return response, svc
}

func fetchSecurityGroups(instance *ec2.DescribeInstancesOutput) []*ec2.GroupIdentifier {
	if len(instance.Reservations) >= 1 {
		for _, inst := range instance.Reservations[0].Instances {
			return inst.SecurityGroups
		}
	}
	return []*ec2.GroupIdentifier{}
}

// Connects on AWS and insert IP on first SecurityGroup
func InsertIPOnSG(ip string, machineName string) bool {
	machineData, svc := getDataFromMachine(machineName)
	sgs := fetchSecurityGroups(machineData)
	for _, securityGroup := range sgs {
		ipExists := IPExistsOnSecurityGroup(ip, svc, sgs)
		if !ipExists {
			addIPToSG(ip, svc, securityGroup.GroupId)
			return true
		}
	}
	return false
}

// Delete IP from SecurityGroup
func deleteIPFromSG(ip string, machineName string) bool {
	machineData, svc := getDataFromMachine(machineName)
	sgs := fetchSecurityGroups(machineData)
	for _, securityGroup := range sgs {
		ipExists := IPExistsOnSecurityGroup(ip, svc, sgs)
		if ipExists {
			delIPFromSG(ip, svc, securityGroup.GroupId)
			return true
		}
	}
	return false
}

// Verify is IP exists on SecurityGroup
func IPExistsOnSecurityGroup(ip string, svc *ec2.EC2, securityGroup []*ec2.GroupIdentifier) bool {
	paramsSG := &ec2.DescribeSecurityGroupsInput{
		DryRun:   aws.Bool(false),
		GroupIds: []*string{aws.String(*securityGroup[0].GroupId)},
	}
	resp, err := svc.DescribeSecurityGroups(paramsSG)
	fmt.Println(resp)
	if err != nil {
		panic(err)
	}
	ipExists := false
	for _, ipPermission := range resp.SecurityGroups[0].IpPermissions {
		for _, ipRange := range ipPermission.IpRanges {
			ipExists = strings.Contains(*ipRange.CidrIp, ip)
			if ipExists {
				return true
			}
		}
	}
	return ipExists
}

// Authorize the IP Address to SecurityGroup
func addIPToSG(ip string, svc *ec2.EC2, groupID *string) *ec2.AuthorizeSecurityGroupIngressOutput {
	paramsReq := &ec2.AuthorizeSecurityGroupIngressInput{
		DryRun:     aws.Bool(false),
		GroupId:    aws.String(*groupID),
		CidrIp:     aws.String(ip + "/32"),
		IpProtocol: aws.String("TCP"),
		FromPort:   aws.Int64(22),
		ToPort:     aws.Int64(22),
	}
	out, err := svc.AuthorizeSecurityGroupIngress(paramsReq)
	if err != nil {
		panic(err)
	}
	return out
}

// Delete IP Address from SecuritGroup
func delIPFromSG(ip string, svc *ec2.EC2, groupID *string) *ec2.RevokeSecurityGroupIngressOutput {
	paramsReq := &ec2.RevokeSecurityGroupIngressInput{
		DryRun:     aws.Bool(false),
		GroupId:    aws.String(*groupID),
		CidrIp:     aws.String(ip + "/32"),
		IpProtocol: aws.String("TCP"),
		FromPort:   aws.Int64(22),
		ToPort:     aws.Int64(22),
	}
	out, err := svc.RevokeSecurityGroupIngress(paramsReq)
	if err != nil {
		panic(err)
	}
	return out
}

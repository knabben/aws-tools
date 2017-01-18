package cmd

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestFetchIps(t *testing.T) {
	sess, _ := session.NewSession()
	svc := ec2.New(sess)

	sgIdentifier := *ec2.GroupIdentifier{
		GroupId:   aws.String("string"),
		GroupName: aws.String("string"),
	}

	// TODO - invalid indirect of ec2.GroupIdentifier literal (type ec2.GroupIdentifier)
	// fetchIpsFromSecurityGroup(svc, sgIdentifier)
}

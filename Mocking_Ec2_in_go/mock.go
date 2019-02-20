package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)


func funcmock(svc ec2iface.EC2API) {
	resp, err := svc.DescribeAvailabilityZones(nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

type mockEC2Client struct {
	ec2iface.EC2API
}

func (m *mockEC2Client) DescribeAvailabilityZones(*ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error) {
	
	return &ec2.DescribeAvailabilityZonesOutput{}, nil

}
func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := ec2.New(sess)
	funcmock(svc)

	mockSvc := &mockEC2Client{}
	resp, err := mockSvc.DescribeAvailabilityZones(nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type mockEC2Client struct {
	ec2iface.EC2API
}

var Instances []*ec2.Instance

var Reservation []*ec2.Reservation

var SpotReservation []*ec2.Reservation

var SpotInstances []*ec2.Instance

func (m *mockEC2Client) DescribeAvailabilityZones(*ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error) {

	a1 := &ec2.AvailabilityZone{}
	a1.SetRegionName("us-east-2")
	a1.SetState("available")
	a1.SetZoneId("use2-az1")
	a1.SetZoneName("us-east-2a")
	a2 := &ec2.AvailabilityZone{}
	a2.SetRegionName("us-east-2")
	a2.SetState("available")
	a2.SetZoneId("use2-az1")
	a2.SetZoneName("us-east-2a")

	out := &ec2.DescribeAvailabilityZonesOutput{}
	out.AvailabilityZones = append(out.AvailabilityZones, a1, a2)
	//out.AvailabilityZones[0].SetRegionName("us-east-2")
	//out.AvailabilityZones[0].SetState("available")
	//out.AvailabilityZones[0].SetZoneId("use2-az1")
	//out.AvailabilityZones[0].SetZoneName("us-east-2a")

	return out, nil
}

func (m *mockEC2Client) RunInstances(input *ec2.RunInstancesInput){
	var arr []string
	arr = append(arr,"i-c1314eb9ed9fff5bf","i-7ee75d6b234c45c50","i-73500e62d4e941389","i-70080e9695ba64373","i-33603ea965ef0b88c","i-b7b02a394535454da")
	var i int
	i = rand.Intn(5)
	state := &ec2.InstanceState{}
	state.SetName("running")
	I1 := &ec2.Instance{}
	I1.SetImageId("ami-1234abcd")
	I1.SetInstanceId(arr[i])
	I1.SetInstanceType(*input.InstanceType)
	I1.SetState(state)

 	out := &ec2.Reservation{}
 	out.Instances = append(out.Instances,I1)
 	Instances = append(Instances, I1)
 	Reservation = append(Reservation,out)


}
func (m *mockEC2Client) RequestSpotInstances(input *ec2.RequestSpotInstancesInput) (*ec2.RequestSpotInstancesOutput){
	var arr []string
	arr = append(arr,"i-c1314eb9ed9fff5bf","i-7ee75d6b234c45c50","i-73500e62d4e941389","i-70080e9695ba64373","i-33603ea965ef0b88c","i-b7b02a394535454da")



	x := &ec2.LaunchSpecification{}
	x.SetInstanceType(*input.LaunchSpecification.InstanceType)
	spotrequest := &ec2.SpotInstanceRequest{}
	spotrequest.SetProductDescription("Windows")
	spotrequest.SetState("Open")
	spotrequest.SetSpotPrice(*input.SpotPrice)
	spotrequest.SetLaunchSpecification(x)

	count := int(*input.InstanceCount)

	for i := 0; i < count; i++ {
		Id := rand.Intn(5)
		state := &ec2.InstanceState{}
		state.SetName("running")
		I1 := &ec2.Instance{}
		I1.SetImageId("ami-1234abcd")
		I1.SetInstanceId(arr[Id])
		I1.SetInstanceType(*x.InstanceType)
		I1.SetState(state)
		SpotInstances = append(SpotInstances, I1)
	}




	out := &ec2.RequestSpotInstancesOutput{}
	out.SpotInstanceRequests= append(out.SpotInstanceRequests, spotrequest)
	return out
}

func (m *mockEC2Client) DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error){

	out := &ec2.DescribeInstancesOutput{}
	out.SetReservations(Reservation)

	return out, nil
}

func (c *mockEC2Client) StopInstances(input []string) (*ec2.StopInstancesOutput, error){

	var Instance_list []*string
	for _, inst := range Instances{
		Instance_list = append(Instance_list, inst.InstanceId)
	}

	for id := range input{
		for _, inst := range Instance_list{
			if input[id] == *inst{
				fmt.Println((1))
			}
		}


	}
	return nil, nil
}

func printInstances(resp *ec2.DescribeInstancesOutput){

	fmt.Println("  > Number of On demand instances: ", len(resp.Reservations))
	fmt.Println("-----------------------------------------")

	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {

			if *inst.State.Name == "running" {
				fmt.Println("Instance State: ", *inst.State.Name)
				fmt.Println("Instance Name: ", *inst.InstanceId)
				fmt.Println("Instance Type", *inst.InstanceType)
				fmt.Println("-----------------------------------------")
			}
			if *inst.State.Name == "stopped" {
				fmt.Println("Instance State: ", *inst.State.Name)
				fmt.Println("Instance Name: ", *inst.InstanceId)
				fmt.Println("Instance Type", *inst.InstanceType)
			}

		}
	}

}

func printSpotInstances(input []*ec2.Instance) {
	fmt.Println("  > Number of Spot instances: ", len(input))
	fmt.Println("-----------------------------------------")
	for _, inst := range input {

		if *inst.State.Name == "running" {
			fmt.Println("Instance State: ", *inst.State.Name)
			fmt.Println("Instance Name: ", *inst.InstanceId)
			fmt.Println("Instance Type", *inst.InstanceType)
			fmt.Println("-----------------------------------------")
		}


	}
}
func Round(input float64) float64 {
	if input < 0 {
		return math.Ceil(input - 0.5)
	}
	return math.Floor(input + 0.5)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc :=  ec2.New(sess)

	mockSvc := &mockEC2Client{}

	var no_on_demand int
	var no_spot int
	var inst_type string

	fmt.Println("Type of instance to be used: ")
	fmt.Scanln(&inst_type)

	input1 := &ec2.DescribeSpotPriceHistoryInput{AvailabilityZone: aws.String("us-east-1a"),
		InstanceTypes: []*string{
			aws.String(inst_type),
		},ProductDescriptions:[]*string{
			aws.String("Windows")},
	}

	result, err := svc.DescribeSpotPriceHistory(input1)
	if err != nil {
		panic(err)
	}
	fmt.Println("Region: ",*result.SpotPriceHistory[0].AvailabilityZone)
	fmt.Println("Instance Type: ",*result.SpotPriceHistory[0].InstanceType)
	fmt.Println("Product Description: ",*result.SpotPriceHistory[0].ProductDescription)
	fmt.Println("Spot Price per hour in USD: ",*result.SpotPriceHistory[0].SpotPrice)
	fmt.Println("-----------------------------------------")

	fmt.Println()
	fmt.Println()

	fmt.Println("The number of on demand instances to be used: ")
	fmt.Scanln(&no_on_demand)

	fmt.Println("The number of spot instances to be used: ")
	fmt.Scanln(&no_spot)



	fmt.Println()
	fmt.Println()



	fmt.Println()
	fmt.Println()

	cost_spot, err := strconv.ParseFloat(*result.SpotPriceHistory[0].SpotPrice, 64)

	cost_on_demand := 0.032

	input2 := &ec2.RunInstancesInput{InstanceType: aws.String(inst_type)}

	for i := 0; i < no_on_demand; i++ {
		mockSvc.RunInstances(input2)
	}


	resp, err := mockSvc.DescribeInstances(nil)

	if err != nil {
		fmt.Println(err)
	}

	printInstances(resp)
	input := &ec2.RequestSpotInstancesInput{
		InstanceCount: aws.Int64(int64(no_spot)),
		LaunchSpecification: &ec2.RequestSpotLaunchSpecification{
			ImageId:      aws.String("ami-1a2b3c4d"),
			InstanceType: aws.String("t2.small"),
			Placement: &ec2.SpotPlacement{
				AvailabilityZone: aws.String("us-east-1a"),
			},

		},
		SpotPrice: aws.String("0.03"),
		Type:      aws.String("one-time"),
	}



	resp2 := mockSvc.RequestSpotInstances(input)

	start := time.Now()

	price := *resp2.SpotInstanceRequests[0].SpotPrice

	fmt.Println(price)

	printSpotInstances(SpotInstances)



	for {
		y := time.Since(start)/time.Second

		if y == 10{
			break
		}
	}



	elapsed := time.Since(start)
	cost_of_running_spot := cost_spot * float64(elapsed) * float64(no_spot)
	cost_of_running_on_demand := cost_on_demand * float64(elapsed) * float64(no_on_demand)
	fmt.Println("Cost of running spot instances for 10 hrs:", (cost_of_running_spot), "USD")
	fmt.Println("Cost of running on demand instances for 10 hrs:", (cost_of_running_on_demand), "USD")

	//cost_saved := cost_on_demand * float64(elapsed) * float64(no_on_demand + no_spot)

	total_cost := cost_of_running_on_demand + cost_of_running_spot

	fmt.Println("Total cost", total_cost)








}

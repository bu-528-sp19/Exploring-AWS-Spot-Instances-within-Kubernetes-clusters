package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/chronojam/aws-pricing-api/types/schema"
	"strings"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := pricing.New(sess)
	input := &pricing.GetAttributeValuesInput{
		MaxResults:    aws.Int64(100),
		ServiceCode:   aws.String("AmazonEC2"),
		AttributeName: aws.String("InstanceType"),
	}

	result, err := svc.GetAttributeValues(input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result.String())

	ec2 := &schema.AmazonEC2{}

	// Populate this object with new pricing data
	err1 := ec2.Refresh()
	if err1 != nil {
		panic(err1)
	}

	// Get the price of all c4.Large instances,
	// running linux, on shared tenancy
	c4Large := []*schema.AmazonEC2_Product{}
	for _, p := range ec2.Products {
		if p.Attributes.InstanceType == "c4.large" &&
			p.Attributes.OperatingSystem == "Linux" &&
			p.Attributes.Tenancy == "Shared" {
			c4Large = append(c4Large, p)
		}
	}

	// Show the pricing data for each of those.
	for _, p := range c4Large {
		//fmt.Println(p.Sku)
		// Find the correct terms
		for _, term := range ec2.Terms {
			if term.Sku == p.Sku {
				for _, pd := range term.PriceDimensions {
					// I Stripped out the OnDemand/Reserved field, but maybe ill add it back later
					// Only On Demand
					if strings.Contains(pd.Description, "On Demand") {
						fmt.Printf("%s:\n", p.Sku)
						fmt.Printf("\t%s:\n", "PriceDimensions")
						fmt.Printf("\t\t%s\n", pd.Description)
						fmt.Printf("\t\t%s\n", pd.PricePerUnit.USD)
					}
				}
			}
		}
	}

}

package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"os/exec"
	"reflect"
	"sort"
	"strconv"
	"time"
)

func calculateValue(values float64, input_range float64) bool {
	if values >= input_range {
		return true
	} else {
		return false
	}
}

func spot_termination_handler(xs string) bool{
	if xs == "to be terminated"{
		return  true
	}else{
		return false
	}
}


func main() {


	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc :=  ec2.New(sess)
	
	resp, err := svc.DescribeInstances(nil)
	if err != nil{
	
	}
	fmt.Println(resp.Reservations)
	var array1 []string = nil
	var start_nodes int = 2
	var spot_nodes int = 1
	var i int = 0
	var j int = 0
	var k int = 0
	TimeStamp := make(map[int]string)
	start_time := int(time.Now().Unix())
	on_demand_price := 0.0464
	spot_price := 0.0125
	total_cost_on_demand := float64(86400) * on_demand_price * 3/3600
	total_cost_on_spot := float64(86400) * on_demand_price * 2/3600 + float64(86400) * spot_price * 1/3600
	fmt.Println(total_cost_on_demand)
	fmt.Println(total_cost_on_spot)
	fmt.Println("% Cost saved = ",(total_cost_on_demand-total_cost_on_spot)/total_cost_on_demand)
	var budget float64 = 20
	count_false := 0.0
	count_true := 100.0






	for t := range time.NewTicker(15 * time.Second).C {
		if t == time.Now() {
			//
		}
		curr_time := int(time.Now().Unix())
		act_time := curr_time - start_time
		time_passed := curr_time - start_time + 3000
		fmt.Println("Time Passed (in mins):", time_passed/60)
		var describe_instances_output string

		if act_time == 60{
			describe_instances_output = "to be terminated"
		}
		if spot_termination_handler(describe_instances_output){
			fmt.Println("=================================================")
			fmt.Println("Spot instance is going to be terminated soon")
			fmt.Println("=================================================")

		}


		//Network traffic
		resp1, err := http.Get("http://af1b1f0726c5011e9b71e02004e2872a-1377931422.us-east-1.elb.amazonaws.com:9090/api/v1/query?query=sum(rate(request_duration_seconds_count[1m]))")
		reqBody1, err := ioutil.ReadAll(resp1.Body)
		if err != nil {
			//
		}
		value1 := gjson.Get(string(reqBody1), "data.result.0.value.1")
		network_traffic, err := strconv.ParseFloat(value1.String(), 64)
		fmt.Println("Number of requests per second(Network traffic)", network_traffic)

		// Throughput
		resp3, err := http.Get("http://af1b1f0726c5011e9b71e02004e2872a-1377931422.us-east-1.elb.amazonaws.com:9090/api/v1/query?query=sum(rate(request_duration_seconds_count{status_code='200'}[1m]))")
		reqBody3, err := ioutil.ReadAll(resp3.Body)
		if err != nil {
			//
		}
		value3 := gjson.Get(string(reqBody3), "data.result.0.value.1")
		throughput, err := strconv.ParseFloat(value3.String(), 64)
		throughput = (throughput / network_traffic) * 100
		fmt.Println("Percent Throughput:", throughput)

		//Latency
		resp4, err := http.Get("http://af1b1f0726c5011e9b71e02004e2872a-1377931422.us-east-1.elb.amazonaws.com:9090/api/v1/query?query=sum(rate(request_duration_seconds_sum[1m]))/sum(rate(request_duration_seconds_count[1m]))")
		reqBody4, err := ioutil.ReadAll(resp4.Body)
		if err != nil {
			//
		}
		value4 := gjson.Get(string(reqBody4), "data.result.0.value.1")
		latency, err := strconv.ParseFloat(value4.String(), 64)
		fmt.Println("Latency:", latency)

		//ErrorRate
		resp5, err := http.Get("http://af1b1f0726c5011e9b71e02004e2872a-1377931422.us-east-1.elb.amazonaws.com:9090/api/v1/query?query=sum(rate(request_duration_seconds_count{status_code='500'}[1m]))")
		reqBody5, err := ioutil.ReadAll(resp5.Body)
		if err != nil {
			//
		}
		value5 := gjson.Get(string(reqBody5), "data.result.0.value.1")
		error_rate, err := strconv.ParseFloat(value5.String(), 64)
		error_rate = (error_rate / network_traffic) * 100
		fmt.Println("Percent Error Rate:", error_rate)

		//Uptime
		resp1, err := http.Get("http://a149d30fb616a11e9939402df7919f80-1474100923.us-east-1.elb.amazonaws.com:9090/api/v1/query?query=sum(request_duration_seconds_sum{status_code='200'})")
		reqBody1, err := ioutil.ReadAll(resp1.Body)
		if err != nil {
			//
		}
		value1 := gjson.Get(string(reqBody1), "data.result.#.value.1")
		var x1 int = len(value1.String())
		uptime, err := strconv.ParseFloat(value1.String()[2:x1-2], 64)
		
		
		// Downtime
		resp2, err := http.Get("http://a149d30fb616a11e9939402df7919f80-1474100923.us-east-1.elb.amazonaws.com:9090/api/v1/query?query=sum(request_duration_seconds_sum{status_code='500'})")
		reqBody2, err := ioutil.ReadAll(resp2.Body)
		if err != nil {
			//
		}
		value2 := gjson.Get(string(reqBody2), "data.result.#.value.1")
		var x int = len(value2.String())
		downtime, err := strconv.ParseFloat(value2.String()[2:x-2], 64)
		fmt.Println("%Uptime:", uptime/(uptime+downtime)*100)

		var input_traffic float64 = 700
		var input_throughput float64 = 350.00               // User defined node load of each of the requests (SLO)
		var input_latency float64 = 1.5 // User defined response time of each of the requests (SLO)
		var input_errorrate float64 = 0.5

		var result_latency= true
		var result_throughput= true
		current_time := strconv.FormatInt(int64(time.Now().Unix()), 10)
		result_latency = calculateValue(throughput, input_throughput) // To check if aggregated node load lies between the input provided by user
		result_throughput = calculateValue(input_latency, latency)        // To check if aggregated response time lies between the input provided by user

		if time_passed == 300{
			budget += 20
		}

		if latency > input_latency || error_rate > input_errorrate { // Maintaining logs using True and False
			TimeStamp[time_passed] = "false"
			array1 = append(array1, "false")
			count_false += 1
		} else {
			TimeStamp[time_passed] = "true"
			array1 = append(array1, "true")
			count_true += 1
		}
		//fmt.Println(array1)
		//fmt.Println(TimeStamp)
		fmt.Println("%Uptime:", ((count_true)/(count_true+count_false))*100)
		fmt.Println("Number of On Demand:",start_nodes)
		fmt.Println("Number of Spot Instances:",spot_nodes)
		//fmt.Println("Count of False:", count_false)
		//fmt.Println("Count of True", count_true)

		var keys []int
		for k := range TimeStamp {
			keys = append(keys, k)
		}

		sort.Sort(sort.Reverse(sort.IntSlice(keys)))

		if network_traffic > 600 {
			high_load_false := 0
			high_load_true := 0

			fmt.Println("High Network")

			//for k :=0; k < len(keys);k++ { // Reversed the Dictionary
			//	v := TimeStamp[keys[k]]
			//	if v == "false" {
			//		high_load_false += 1
			//	}else {
			//		high_load_true += 1
			//	}
			//}

			//	for k := len(array1)-1 ; k >= 0 ;k-- { // Reversed the Dictionary
			if k < 40 {
				k++
				for i := 0; i < k; i++ { // Reversed the Dictionary
					v := TimeStamp[keys[i]]
					if v == "false" {
						high_load_false += 1
					} else {
						high_load_true += 1
					}
				}
			}

			//fmt.Println("k", k)
			//fmt.Println("false", high_load_false)
			//fmt.Println("true", high_load_true)
			if k == 40 && high_load_false > 25 {
				if budget > 10 {

					fmt.Println("Budget available so spin up a spot instance")
					exec.Command("/bin/sh", "addspot.sh").Output()
					time.Sleep(time.Duration(300) * time.Second)
					fmt.Println("Added spot")
					spot_nodes += 1
					exec.Command("kubectl", "scale", "deployments", "--all", "-n", "sock-shop", "--replicas=4").Output()
					time.Sleep(time.Duration(30) * time.Second)
					fmt.Println("Scaled up")
					budget = budget - 10
					fmt.Println("Budget:", budget)
					fmt.Println("Number of Spot nodes:", spot_nodes)
					fmt.Println("Number of On demand:", start_nodes)
				} else {
					fmt.Println("High Network")
					fmt.Println("Budget available so spin up a on demand instance")
					exec.Command("/bin/sh", "nodeaddition.sh").Output()
					time.Sleep(time.Duration(300) * time.Second)
					fmt.Println("Added on demand ")
					start_nodes += 1
					exec.Command("kubectl", "scale", "deployments", "--all", "-n", "sock-shop", "--replicas=4").Output()
					time.Sleep(time.Duration(30) * time.Second)
					fmt.Println("Scaled up")
					budget = budget - 10
					fmt.Println("Budget:", budget)
					fmt.Println("Number of Spot nodes:", spot_nodes)
					fmt.Println("Number of On demand:", start_nodes)
				}
				k = 0
			}



			//}
		}else if network_traffic < 100 {
			low_load_false := 0
			low_load_true := 0

			if i < 240 {
				i++
				for x := 0; x < i; x++ { // Reversed the Dictionary
					v := TimeStamp[keys[x]]
					if v == "false" {
						low_load_false += 1
					} else {
						low_load_true += 1
					}
				}
			}

			fmt.Println("i", i)
			fmt.Println("false", low_load_false)
			fmt.Println("true", low_load_true)
			if i == 240 {
				if low_load_true >= 200 {
					if budget > 10 {
						if start_nodes > 2 {
							fmt.Println("if number of nodes  > 2 if Budget available scale down ")
							exec.Command("kubectl", "scale", "deployments", "--all", "-n", "sock-shop", "--replicas=2").Output()
							fmt.Println("Sleep for 30 secs")
							time.Sleep(time.Duration(30) * time.Second)
							exec.Command("/bin/sh", "nodesubtraction.sh").Output()
							time.Sleep(time.Duration(200) * time.Second)
							budget = budget - 5
							fmt.Println("Deleted an On demand node")
							fmt.Println("Budget: ", budget)
							fmt.Println("Number of Spot nodes:", spot_nodes)
							fmt.Println("Number of On demand:", start_nodes)
						}else if spot_nodes > 0 {
							fmt.Println("Remove spot")
						}
					}else {
						fmt.Println("No budget available to do anything")
					}
				}
				i = 0
			}
		}else if network_traffic > 100 && network_traffic < 600 {
			mod_load_false := 0
			mod_load_true := 0
			if j < 240 {
				j++
				for i := 0; i < j; i++ { // Reversed the Dictionary
					v := TimeStamp[keys[i]]
					if v == "false" {
						mod_load_false += 1
					} else {
						mod_load_true += 1
					}
				}
			}

			//fmt.Println("j", j)
			//fmt.Println("false", mod_load_false)
			//fmt.Println("true", mod_load_true)
			if k == 240 {
				if mod_load_true > 200 {
					if budget > 10 {
						if start_nodes > 2 {
							fmt.Println("Budget available so swap on demand with a spot instance")
							exec.Command("kubectl", "scale", "deployments", "--all", "-n", "sock-shop", "--replicas=2").Output()
							time.Sleep(time.Duration(30) * time.Second)
							fmt.Println("Scaled down")
							exec.Command("/bin/sh", "nodesubtraction.sh").Output()
							start_nodes -= 1
							time.Sleep(time.Duration(200) * time.Second)
							fmt.Println("Node subtarcted")
							exec.Command("/bin/sh", "addspot.sh").Output()
							time.Sleep(time.Duration(300) * time.Second)
							fmt.Println("Added spot")
							spot_nodes += 1
							exec.Command("kubectl", "scale", "deployments", "--all", "-n", "sock-shop", "--replicas=3").Output()
							time.Sleep(time.Duration(30) * time.Second)
							fmt.Println("Scaled up")
							budget = budget - 10
							fmt.Println("Budget:", budget)
							fmt.Println("Number of Spot nodes:", spot_nodes)
							fmt.Println("Number of On demand:", start_nodes)
						}
					} else {
						fmt.Println("Budget not available so dont do anything")
					}
				}
				j=0
			}else if act_time == 60{
				fmt.Println("Scale down the application , Add on demand node and scale up")
			}
		}

		if time_passed == 210{
			total_cost_on_demand := float64(time_passed) * on_demand_price * 3/3600
			total_cost_on_spot := float64(time_passed) * on_demand_price * 2/3600 + float64(time_passed) * spot_price * 1/3600
			fmt.Println(total_cost_on_demand)
			fmt.Println(total_cost_on_spot)
			fmt.Println("% Cost saved = ",(total_cost_on_demand-total_cost_on_spot)/total_cost_on_demand)
		}

		if count_false > 5{                                // Adding new nodes
			//out1,err1 := exec.Command("cd","nodeaddition.sh").Output()
			//out2,err2 := exec.Command("cd","nodeaddition.sh").Output()
			//out3,err3 := exec.Command("cd","nodeaddition.sh").Output()
			out,err := exec.Command("/bin/sh","nodeaddition.sh").Output()
		
			if err != nil {
				//error := string(err[:])
				fmt.Println("The type of err",reflect.TypeOf(err))
				fmt.Println("error %s", err)
			}
			output := string(out[:])
			fmt.Printf(output)
		
		}


		fmt.Println("---------------------------------------------------------------------------")





		
		// ------------------------------------------------------------------------------------------------------------
		current_time := strconv.FormatInt(int64(time.Now().Unix()), 10)
		var input_node_load float64 = 8.0               // User defined node load of each of the requests (SLO)
		var input_latency float64 = 9.0                 // User defined response time of each of the requests (SLO)
		var input_SLA float64 = 95.00                   // SLA given by the user
		var availability float64 = uptime / (downtime + uptime) //availability

		var result_latency= true
		var result_response= true
		avg_Node_Load := 10.00     // Node load from prometheus
		latency := 5.00            // Latency from prometheus
		Node_Creation_Time := 3.00 // In seconds
		Spot_Creation_Time := 2.00 // In seconds

		result_latency = calculateValue(input_node_load, avg_Node_Load) // To check if aggregated node load lies between the input provided by user
		result_response = calculateValue(input_latency, latency)        // To check if aggregated response time lies between the input provided by user

		if availability < input_SLA || result_response == false || result_latency == false { // Maintaining logs using True and False
			TimeStamp.Set(current_time, false)
		} else {
			TimeStamp.Set(current_time, true)
		}
		

	}

}

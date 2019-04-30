package main

import (
	"fmt"
	"github.com/iancoleman/orderedmap"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	//"os/exec"
	//"reflect"
	"strconv"
	"time"
)


//func calculateValue(values float64, input_range float64) bool {
//	if values >= input_range {
//		return true
//	} else {
//		return false
//	}
//}


var TimeStamp = orderedmap.New()

func main() {

	start_time := int64(time.Now().Unix())
	on_demand_price := 0.0464
	spot_price := 0.0125
	total_cost_on_demand := float64(86400) * on_demand_price * 3/3600
	total_cost_on_spot := float64(86400) * on_demand_price * 2/3600 + float64(86400) * spot_price * 1/3600
	fmt.Println(total_cost_on_demand)
	fmt.Println(total_cost_on_spot)
	fmt.Println("% Cost saved = ",(total_cost_on_demand-total_cost_on_spot)/total_cost_on_demand)


	for t := range time.NewTicker(15 * time.Second).C {
		if t == time.Now(){
			//
		}
		curr_time := int64(time.Now().Unix())
		time_passed := curr_time - start_time
		fmt.Println("Time Passed:", time_passed)

		//Network traffic
		resp1, err := http.Get("http://ad85538dd6aad11e9965402f420e336c-1063311824.us-east-1.elb.amazonaws.com:9090/api/v1/query?query=sum(rate(request_duration_seconds_count[1m]))")
		reqBody1, err := ioutil.ReadAll(resp1.Body)
		if err != nil {
			//
		}
		value1 := gjson.Get(string(reqBody1), "data.result.0.value.1")
		network_traffic ,err :=  strconv.ParseFloat(value1.String(),64)
		fmt.Println("Number of requests per second(Network traffic)", network_traffic)

		// Throughput
		resp3, err := http.Get("http://ad85538dd6aad11e9965402f420e336c-1063311824.us-east-1.elb.amazonaws.com:9090/api/v1/query?query=sum(rate(request_duration_seconds_count{status_code='200'}[1m]))")
		reqBody3, err := ioutil.ReadAll(resp3.Body)
		if err != nil {
			//
		}
		value3 := gjson.Get(string(reqBody3), "data.result.0.value.1")
		throughput ,err :=  strconv.ParseFloat(value3.String(),64)
		fmt.Println("Throughput:",throughput)


		//Latency
		resp4, err := http.Get("http://ad85538dd6aad11e9965402f420e336c-1063311824.us-east-1.elb.amazonaws.com:9090/api/v1/query?query=sum(rate(request_duration_seconds_sum[1m]))/sum(rate(request_duration_seconds_count[1m]))")
		reqBody4, err := ioutil.ReadAll(resp4.Body)
		if err != nil {
			//
		}
		value4 := gjson.Get(string(reqBody4), "data.result.0.value.1")
		latency ,err :=  strconv.ParseFloat(value4.String(),64)
		fmt.Println("Latency:",latency)


		//ErrorRate
		resp5, err := http.Get("http://ad85538dd6aad11e9965402f420e336c-1063311824.us-east-1.elb.amazonaws.com:9090/api/v1/query?query=sum(rate(request_duration_seconds_count{status_code='500'}[1m]))")
		reqBody5, err := ioutil.ReadAll(resp5.Body)
		if err != nil {
			//
		}
		value5 := gjson.Get(string(reqBody5), "data.result.0.value.1")
		error_rate ,err :=  strconv.ParseFloat(value5.String(),64)
		fmt.Println("Error Rate:",error_rate)

		// Uptime
		//resp1, err := http.Get("http://a149d30fb616a11e9939402df7919f80-1474100923.us-east-1.elb.amazonaws.com:9090/api/v1/query?query=sum(request_duration_seconds_sum{status_code='200'})")
		//reqBody1, err := ioutil.ReadAll(resp1.Body)
		//if err != nil {
		//	//
		//}
		//value1 := gjson.Get(string(reqBody1), "data.result.#.value.1")
		//var x1 int = len(value1.String())
		//uptime, err := strconv.ParseFloat(value1.String()[2:x1-2], 64)
		//
		//
		//// Downtime
		//resp2, err := http.Get("http://a149d30fb616a11e9939402df7919f80-1474100923.us-east-1.elb.amazonaws.com:9090/api/v1/query?query=sum(request_duration_seconds_sum{status_code='500'})")
		//reqBody2, err := ioutil.ReadAll(resp2.Body)
		//if err != nil {
		//	//
		//}
		//value2 := gjson.Get(string(reqBody2), "data.result.#.value.1")
		//var x int = len(value2.String())
		//downtime, err := strconv.ParseFloat(value2.String()[2:x-2], 64)
		//fmt.Println("%Uptime:", uptime/(uptime+downtime)*100)

		var budget float64 = 20
		//var input_traffic float64 = 700
		var input_throughput float64 = 350.00               // User defined node load of each of the requests (SLO)
		var input_latency float64 = 1.5                    // User defined response time of each of the requests (SLO)
		var input_errorrate float64 = 5

		//var result_latency= true
		//var result_throughput= true
		//current_time := strconv.FormatInt(int64(time.Now().Unix()), 10)
		//result_latency = calculateValue(throughput, input_throughput) // To check if aggregated node load lies between the input provided by user
		//result_throughput = calculateValue(input_latency, latency)        // To check if aggregated response time lies between the input provided by user

		if time_passed == 300{
			budget = 20
		}

		if throughput < input_throughput || latency >  input_latency || error_rate > input_errorrate { // Maintaining logs using True and False
			TimeStamp.Set(strconv.FormatInt(time_passed, 10), false)
		} else {
			TimeStamp.Set(strconv.FormatInt(time_passed, 10), true)
		}

		keys := TimeStamp.Keys() // Reversing the TimeStamp
		for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
			keys[i], keys[j] = keys[j], keys[i]
		}

		count_false := 0.0
		count_true := 0.0


		for _, k := range keys { // Reversed the Dictionary
			v, _ := TimeStamp.Get(k)
			if v == false {
				count_false += 1
			}else {
				count_true += 1
			}
		}
		fmt.Println(TimeStamp)
		fmt.Println("%Uptime:", ((count_true)/(count_true+count_false))*100)
		fmt.Println("Count of False:", count_false)
		fmt.Println("Count of True", count_true)

		if network_traffic > 600 {

			c_false := 0
			c_true := 0
			for _, k := range keys { // Reversed the Dictionary
				v, _ := TimeStamp.Get(k)
				if c_true + c_false == 20{
					if c_false > 10{
						if budget > 10{
							fmt.Println("Budget available so spin up a spot instance")
							budget = budget - 10
						}else{
							fmt.Println("Budget not available so create an on demand")
						}
						break
					}
				}
				if v == false {
					c_false += 1
				}else {
					c_true += 1
				}
			}
		}else {
			fmt.Println("Check budget add spot")
		}



		if time_passed == 210{
			total_cost_on_demand := float64(time_passed) * on_demand_price * 3/3600
			total_cost_on_spot := float64(time_passed) * on_demand_price * 2/3600 + float64(time_passed) * spot_price * 1/3600
			fmt.Println(total_cost_on_demand)
			fmt.Println(total_cost_on_spot)
			fmt.Println("% Cost saved = ",(total_cost_on_demand-total_cost_on_spot)/total_cost_on_demand)
		}



		//if count_false > 5{                                // Adding new nodes
		//	//out1,err1 := exec.Command("cd","nodeaddition.sh").Output()
		//	//out2,err2 := exec.Command("cd","nodeaddition.sh").Output()
		//	//out3,err3 := exec.Command("cd","nodeaddition.sh").Output()
		//	out,err := exec.Command("/bin/sh","nodeaddition.sh").Output()
		//
		//	if err != nil {
		//		//error := string(err[:])
		//		fmt.Println("The type of err",reflect.TypeOf(err))
		//		fmt.Println("error %s", err)
		//	}
		//	output := string(out[:])
		//	fmt.Printf(output)
		//
		//}


		fmt.Println("---------------------------------------------------------------------------")





		/*
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
		*/
	}

}

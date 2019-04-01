package main

import (
	"fmt"
	"sort"
	"math"
)

func calculate_value(values []float64, input_range []float64) bool{

	var temp float64 = 0
	for _,v := range(values) {
		temp += v
	}
	if input_range[0] <=  (temp/float64(len(values))) && (temp/float64(len(values))) <= input_range[1] {
			return true
	}else {
		return false
	}
}

func calculate_SLA(result bool, uptime int,downtime int) (float64,int ,int){
	if result == true{
		if downtime == 0 {
			return 100.0, uptime + 1, downtime   // If there is no downtime, i.e. downtime ==0
		}else{
			var sla float64 = ((float64(uptime+1)/float64(uptime+1+downtime))*100)  // If there is no downtime, i.e. downtime !=0
			return math.Round(sla*100)/100, uptime + 1, downtime // Rounding to 2 decimal places
		}
	}else{
		var sla float64 = ((float64(uptime)/float64(uptime+1+downtime))*100)
		return math.Round(sla*100)/100, uptime + 1, downtime  // Rounding to 2 decimal places
	}

}

func main() {

	var user_input = map[int]string{  // Input from the user for each of the SLI
		3: "latency",
		2: "response_time",
	}
	var priority = make([]int, 0, len(user_input))  // Storing all the priority in each a list and will sort it,as GO iterates over a map randomly
	for name := range user_input {
		priority = append(priority, name)
	}
	sort.Ints(priority)   //The sorting according to priority of each of the SLI

	var input_latency = []float64{4.0, 8.0} // User defined latency of each of the requests (SLO)
	var input_response_time = []float64{3.0,9.0} // User defined response time of each of the requests (SLO)

	var latency_uptime int = 1   // The uptime of Latency
	var latency_dowtime int = 1  // The downtime of Latency
	var response_time_uptime int = 1 // The uptime of Response Time
	var response_time_downtime int = 1  // The downtime of Response Time


	var latency = []float64{10.0, 2.0, 3.4, 7.0, 5.0} // The latency of each of the requests (SLI)
	var response_time = []float64{9.0,3.0,7.4,5.9,8.1} // The response of each of the requests (SLI)

	for _, value := range priority {

		if user_input[value] == "latency"{
			var result_latency = calculate_value(latency , input_latency) // To check if average of all the latency in last five minutes lies between the input provided by the user
			var sla_latency float64 = 0.0
			sla_latency,latency_uptime,latency_dowtime  = calculate_SLA(result_latency,latency_uptime,latency_dowtime) // Will increase the uptime or downtime of Latency according to the  result_latency
			fmt.Println("The SLA according to Latency : ",sla_latency)

		}else if user_input[value] == "response_time"{
			var result_response = calculate_value(response_time , input_response_time)  // To check if average of all the response time in last five minutes lies between the input provided by the user
			var sla_respomse_time float64 = 0.0
			sla_respomse_time,response_time_uptime,response_time_downtime  = calculate_SLA(result_response,response_time_uptime,response_time_downtime) // Will increase the uptime or downtime of Response Time according to the  result_response_time
			fmt.Println("The SLA according to Response Time : ",sla_respomse_time)
		}

	}


}

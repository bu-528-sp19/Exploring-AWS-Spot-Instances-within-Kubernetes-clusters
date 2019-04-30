package main

import (
	"sort"
)

func calculateValue(values []float64, input_range float64) bool {
	if input_range <= values[1] && values[0] <= input_range {
		return true
	} else {
		return false
	}
}

func main() {
	var user_input = map[int]string{ // Input from the user for each of the SLI
		1: "node_load",
	}
	var priority = make([]int, 0, len(user_input)) // Storing all the priority in each a list and will sort it,as GO iterates over map randomly
	for name := range user_input {
		priority = append(priority, name)
	}
	sort.Ints(priority) //The sorting according to priority of each of the SLI

	var budget float64 = 10                       // Budget for the application critical metrics and will reduce if the any of critical metrics goes down
	var input_node_load = []float64{4.0, 8.0}     // User defined node load of each of the requests (SLO)
	var input_response_time = []float64{3.0, 9.0} // User defined response time of each of the requests (SLO)
	var input_SLA float64 = 99.00                 // SLA given by the user
	var list_of_uptimes = []float64{}             // The list fo all the previous uptime
	var downtime_delta = []float64{}              // The list fo all the difference between two uptime
	var uptime_state bool = true                  // Storing the state of uptime for calculating delta
	var new_uptime float64 = 1000                 // The new uptime value pulled from prometheus in seconds
	var total_time float64 = 10000                // Total time for which the application is running from prometheus
	var downtime = total_time - new_uptime        // Downtime for the application
	var n int = len(list_of_uptimes) - 1          // The length of all the elements in previous uptime
	var aggregated_node_load float64 = 7          // The aggregated value from prometheus for node load
	var aggregated_response_time float64 = 7      // The aggregated value from prometheus for response time

	if total_time != new_uptime && len(list_of_uptimes) == 0 && uptime_state == true { // When downtime happens for the first time
		uptime_state = false
		list_of_uptimes = append(list_of_uptimes, new_uptime)
	} else if list_of_uptimes[n] != new_uptime { // When the uptime begins after a downtime and storing the value of Delta
		uptime_state = true
		if len(list_of_uptimes) == 1 {
			downtime_delta = append(downtime_delta, total_time-list_of_uptimes[n])
		} else {
			var temp float64 = 0
			for i := 0; i < len(list_of_uptimes); i++ { // Adding all the elements in Uptime list
				temp = temp + list_of_uptimes[i]
			}
			for i := 0; i < len(downtime_delta); i++ { // Adding all the elements in Delta list
				temp = temp + downtime_delta[i]
			}
			downtime_delta = append(downtime_delta, total_time-temp) // Storing the value of Delta
		}
	} else if list_of_uptimes[n] == new_uptime && uptime_state == true { // When the downtime time begins and storing the previous uptime
		uptime_state = false
		var temp float64 = 0
		for i := 0; i < len(list_of_uptimes); i++ { // Adding all the elements in Uptime list
			temp = temp + list_of_uptimes[i]
		}
		for i := 0; i < len(downtime_delta); i++ { // Adding all the elements in Delta list
			temp = temp + downtime_delta[i]
		}
		list_of_uptimes = append(list_of_uptimes, total_time-temp) // Storing the value of Uptime
	}
	for _, value := range priority { // For reducing the budget if its less or more than the user defined range
		if user_input[value] == "node_load" { // Highest priority
			var result_latency = calculateValue(input_node_load, aggregated_node_load) // To check if aggregated node load lies between the input provided by user
			if result_latency == false {                                               //  Reducing the budget according to the priority
				budget -= 1
			}
		} else if user_input[value] == "response_time" { // 2nd Highest priority
			var result_response = calculateValue(input_response_time, aggregated_response_time) // To check if aggregated response time lies between the input provided by user
			if result_response == false {                                                       //  Reducing the budget according to the priority
				budget -= 0.5
			}
		}
	}
	if budget <= 0 {
		//If no node is available in cluster  with cpu usage is high:
		//If spot instance available below spot price:
		//Create a new  spot instance  and distribute load
		//Else:
		//Create a on demand node and distribute load
		//Else:
		//Distribute the load in the available nodes in cluster
		//(Wait for 5- 10 minutes before the system becomes stable and the error rate goes down)
	}
}

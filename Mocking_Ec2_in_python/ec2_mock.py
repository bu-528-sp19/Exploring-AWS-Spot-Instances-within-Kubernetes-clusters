
from moto import mock_ec2
import time
import boto3
import random


def create_instances(ami_id,count):
    client = boto3.client('ec2')
    client.run_instances(ImageId=ami_id, MinCount = count, MaxCount = count)

@mock_ec2
def main():
    start =  time.time()
    for i in range(5):
        create_instances('ami-1234abcd', 1)

    client = boto3.client('ec2', region_name='us-east-2')
    response =client.describe_instances()

    id_list = []

    for i in range(5):
        print("Instance ", i+1 )
        id_list.append(response['Reservations'][i]["Instances"][0]["InstanceId"])
        print(" - Instance Type: ", response['Reservations'][i]["Instances"][0]["InstanceType"])
        print(" - Instance ID: ", response['Reservations'][i]["Instances"][0]["InstanceId"])
        print(" - State: ", response['Reservations'][i]["Instances"][0]["State"]["Name"])

    print()
    print("------------------------------------------------------------------")
    print()

    m = 0
    stop_iteration = random.randint(0, 1000000)
    while(m<=1000000):

        if m == stop_iteration:
            stop_time = time.time() - start
            stop_list = []
            rand_inst = random.randint(0, 4)
            stop_list.append(id_list[rand_inst])
            client.stop_instances(InstanceIds=stop_list)
            client = boto3.client('ec2', region_name='us-east-2')
            response = client.describe_instances()

            print("Instance", id_list[rand_inst], " stopped after", round(stop_time,2),"secs.")
            print()
            print("------------------------------------------------------------------")
            print()
            for i in range(5):
                print("Instance ", i + 1)
                print(" - Instance Type: ", response['Reservations'][i]["Instances"][0]["InstanceType"])
                print(" - Instance ID: ", response['Reservations'][i]["Instances"][0]["InstanceId"])
                print(" - State: ", response['Reservations'][i]["Instances"][0]["State"]["Name"])

            print()
            print("------------------------------------------------------------------")
            print()

            create_instances('ami-1234abcd',1)
            client = boto3.client('ec2', region_name='us-east-2')

            response = client.describe_instances()

            max = 0
            instance = []
            for i in range(6):
                print("Instance ", i + 1)
                print(" - Instance Type: ", response['Reservations'][i]["Instances"][0]["InstanceType"])
                print(" - Instance ID: ", response['Reservations'][i]["Instances"][0]["InstanceId"])
                print(" - State: ", response['Reservations'][i]["Instances"][0]["State"]["Name"])

        m = m + 1

main()

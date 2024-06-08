package main

import "fmt"

func test() {
	var output Output
	device, err := readSingleDeviceInfoFromJSON("./device.json")
	if err != nil {
		panic(err)
	}

	output, err = ProcessDevice(device.Auth, device.UserPasskeyHistory, device.DeviceFeatures)
	if err != nil {
		panic(err)
	}

	fmt.Println("PasskeyType : ",output.MatchPassKeyType,"\n","Match Probability : ",output.MatchProbability)
}
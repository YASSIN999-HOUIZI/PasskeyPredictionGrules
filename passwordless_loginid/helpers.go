package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	s "strings"
)

type PassKeyType string

func contains(list []string, element string) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}

// Helper function to check if a string is not in a list
func notContains(list []string, element string) bool {
	return !contains(list, element)
}

func CheckPassKeyType(passkey PassKeyType) bool {
	return !notContains([]string{"local", "device", "cloud", "nopasskey"}, string(passkey))
}

var PassKeys = struct {
	Local     PassKeyType
	Device    PassKeyType
	Cloud     PassKeyType
	NoPasskey PassKeyType
}{
	Local:     "local",
	Device:    "device",
	Cloud:     "cloud",
	NoPasskey: "nopasskey",
}

func readDeviceInfoFromJSON(fileName string) ([]*DeviceData, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var devices []*DeviceData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&devices)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func readSingleDeviceInfoFromJSON(fileName string) (DeviceData, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return DeviceData{}, err
	}
	defer file.Close()

	var device DeviceData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&device)
	if err != nil {
		return DeviceData{}, err
	}
	return device, nil
}

func writeDeviceInfoToCSV(fileName string, devices []*DeviceData) error {
	csvFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)

	fileInfo, err := csvFile.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() == 0 {
		header := []string{"Label", "OSName", "ptype", "prob"}
		err = csvWriter.Write(header)
		if err != nil {
			return err
		}
	}

	for _, device := range devices {
		record := []string{device.Auth.DeviceID, device.Auth.OsName, device.Output.MatchPassKeyType, strconv.FormatFloat(device.Output.MatchProbability, 'f', -1, 64)}
		err = csvWriter.Write(record)
		if err != nil {
			return err
		}
	}

	sep := []string{"================================="}
	err = csvWriter.Write(sep)
	if err != nil {
		return err
	}

	csvWriter.Flush()
	return nil
}

type CloudClient struct {
	Platform      string
	ClientName    string
	ClientVersion string
}

func getPredefinedVersions(platform string) map[string]string {

	lowerPlat := s.ToLower(platform)

	switch lowerPlat {
	case "ios":
		return map[string]string{
			"safari":  "16",
			"chrome":  "16",
			"edge":    "16",
			"firefox": "16",
		}
	case "macos":
		return map[string]string{
			"safari":  "16.1",
			"chrome":  "118",
			"edge":    "122",
			"firefox": "122",
		}
	case "android":
		return map[string]string{
			"chrome":  "108",
			"edge":    "122",
			"samsung": "21",
		}
	default:
		return nil
	}
}

func isVersionMatching(externalData CloudClient) bool {
	predefinedVersions := getPredefinedVersions(externalData.Platform)
	if predefinedVersions == nil {
		fmt.Println("Unknown platform ", externalData.Platform)
		return false
	}

	expectedVersion, exists := predefinedVersions[externalData.ClientName]
	if !exists {
		fmt.Printf("Unknown client name: %s\n", externalData.ClientName)
		return false
	}

	return externalData.ClientVersion == expectedVersion
}

type IntermediatePasskeyPrediction struct {
	Prediction string
}

type IntermediateCloudPrediction struct {
	Prediction string
}

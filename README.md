# Passkey Prediction using Grule

## Prerequistes
This project requires the following:
<br>
Go version 1.18 or later (https://go.dev/doc/install)

## Description
This repository implements a passkey prediction system using a Grule and calculates the probability of the passkey existence based on user history and device information.
- Grule Engine: Grule allows defining rules based on device information properties. These rules determine the possibility of a specific passkey type being used.
- Match Probability Function: The Go function, MatchProbability, calculates the probability of a passkey existence based on the user's specific history and device information. It considers
Device information properties (Device ID, OS name and version, client name and version, Device size).

## Install and Run
1. Open a new shell
2. Clone the repository to your machine
3. Navigate to the project directory
```terminal
cd loginid-decision-engine/passwordless_loginid
```
4. Run the test file 
- Windows
```terminal
.\test.bat
```
- Linux
```terminal
./test.sh
```
5. Run the rule engine as a server
If you prefer to run the rule engine as a server, you can do so via the following commands:
- Windows
```terminal
.\server.bat
```
- Linux
```terminal
./server.sh
```
#### With the Rest client extension installed on vs code, use the Request.http file to test the API 

## Important details related to the integration
The function that runs the rule engine is located in passwordless_loginid/rule_engine and has the following signature: 
- #### func ProcessDevice(devicesinfos DeviceInfo, userPasskeyHistory []*UserPasskeyHistory, deviceFeatures DeviceFeature) (Output,error)

The Output is a structure of two properties: 
- MatchProbability := 0 {least likely} to 100 {most likely} based on the user history
- MatchPasskeyType := the matching passkey type  [clooud | device | local]

An usage example of this function is located in passwordless_loginid/test.go

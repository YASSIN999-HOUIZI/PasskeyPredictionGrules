# Passkey Prediction using Grules

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
2. Copy the repository to your machine
```terminal
git clone https://github.com/YASSIN999-HOUIZI/PasskeyPredictionGrules.git
```
3. Navigate to the project directory and Run the demo
```terminal
cd PasskeyPredictionGrules/passwordless_loginid
go run helpers.go main.go device_fact.go rules.go
```


### Input File Location: passwordless_loginid/devices.json
### Output File Location: passwordless_loginid/device_info.csv

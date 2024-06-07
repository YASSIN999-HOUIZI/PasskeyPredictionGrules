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
cd PasskeyPredictionGrules/passwordless_loginid
```
4. Launch the server
```terminal
./Server.bat
```

### Input File Location: passwordless_loginid/devices.json
### Output File Location: passwordless_loginid/device_info.csv

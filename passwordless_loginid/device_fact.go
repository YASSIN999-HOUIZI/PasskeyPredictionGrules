package main

import (
	"fmt"
	"strconv"
	s "strings"
)

type DeviceInfo struct {
	OsName        string
	OsVersion     string
	ClientName    string
	ClientVersion string
	IpAddress     string
	DeviceSize    string
	DeviceID      string
}

type DeviceFeature struct {
	NavigLanguage string
	IsCMA         bool
	IsUVPPA       bool
	PasskeyType   string
}

type UserPasskeyHistory struct {
	DeviceInfo  DeviceInfo
	PasskeyType string
}

type Output struct {
	MatchPassKeyType string //PassKeyType
	MatchProbability float64
	// MatchPasskeyType float64
}

type DeviceFact struct {
	Auth               DeviceInfo
	UserPasskeyHistory []*UserPasskeyHistory
	DeviceFeatures     DeviceFeature
	Output             Output
}

func (df *DeviceFact) IsPassKeyExisting() bool {
	return df.DeviceFeatures.IsUVPPA
}

func (df *DeviceFact) MatchOS() bool {
	for _, history := range df.UserPasskeyHistory {
		if df.Auth.OsName != history.DeviceInfo.OsName || df.Auth.OsVersion != history.DeviceInfo.OsVersion {
			return false
		}
	}
	return true
}

func (df *DeviceFact) MatchDeviceProperties() bool {
	
	for _, deviceInfo := range df.UserPasskeyHistory {
		// if deviceInfo.DeviceInfo != df.Auth {
		// 	return false
		// }
		if deviceInfo.DeviceInfo.DeviceID != df.Auth.DeviceID || deviceInfo.DeviceInfo.DeviceSize != df.Auth.DeviceSize{
			return false
		}
		if s.ToLower(deviceInfo.PasskeyType) == "local"{
			if deviceInfo.DeviceInfo.ClientName != df.Auth.ClientName || deviceInfo.DeviceInfo.ClientVersion != df.Auth.ClientVersion{
				return false
		}	
		}
		
	}
	fmt.Println(df.Auth.OsName,true)

	return true
}

func (df *DeviceFact) MatchCloudCompatibleBrowser() bool {
	return isVersionMatching(CloudClient{
		Platform:      s.ToLower(df.Auth.OsName),
		ClientName:    s.ToLower(df.Auth.ClientName),
		ClientVersion: df.Auth.ClientVersion,
	})
}

func (df *DeviceFact) IsConditionalMediationAvailable() bool {
	return df.DeviceFeatures.IsCMA
}

func (df *DeviceFact) IsCompatibleDevice() bool {
	return s.ToLower(df.Auth.OsName) != "windows"
}

func (df *DeviceFact) MatchProbability() float64 {
	probability := 0.0

	for _, pastDevice := range df.UserPasskeyHistory {
		weight := 0.0
		if df.Auth.DeviceID == pastDevice.DeviceInfo.DeviceID {
			probability = 100
			return probability
		}

		// Reduce weight for OS version mismatch
		if df.Auth.OsName == pastDevice.DeviceInfo.OsName {
			weight = 1
			versionDiff := compareVersions(df.Auth.OsVersion, pastDevice.DeviceInfo.OsVersion)
			if versionDiff > 0 {
				weight *= 0.8 // upgrade most likely to happen
			} else if versionDiff < 0 {
				weight *= 0.7 // downgrade less likely to happen
			}

			// Reduce weight for client version mismatch
			if df.Auth.ClientName == pastDevice.DeviceInfo.ClientName {
				versionDiff := compareVersions(df.Auth.ClientVersion, pastDevice.DeviceInfo.ClientVersion)
				if versionDiff > 0 {
					weight *= 0.9 // upgrade
				} else if versionDiff < 0 {
					weight *= 0.8 // downgrade
				}
			} else {
				weight *= 0.85
			}
		} else {
			if df.Auth.ClientName == pastDevice.DeviceInfo.ClientName {
				weight = 0.75
			} else {
				weight *= 0.75
			}
		}
		if df.Auth.DeviceSize != pastDevice.DeviceInfo.DeviceSize { // penalty if the device size is different
			weight *= 0.85
		}

		newProbability := weight * 100
		if newProbability > probability { // We keep the greatest probability knowing entries probability are independent
			probability = newProbability
		}
	}

	fmt.Println("test:", probability)

	return min(probability, 90)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func compareVersions(v1 string, v2 string) int {
	v1Parts := s.Split(v1, ".")
	v2Parts := s.Split(v2, ".")

	minLength := len(v1Parts)
	if len(v2Parts) < minLength {
		minLength = len(v2Parts)
	}

	for i := 0; i < minLength; i++ {
		v1Num, err := strconv.Atoi(v1Parts[i])
		if err != nil {
			return 0 // Handle error or use different logic if needed
		}
		v2Num, err := strconv.Atoi(v2Parts[i])
		if err != nil {
			return 0 // Handle error or use different logic if needed
		}
		if v1Num > v2Num {
			return 1 // Upgrade
		} else if v1Num < v2Num {
			return -1 // Downgrade
		}
	}

	// Versions are equal up to the minimum length, consider versions with more parts as upgrades
	return len(v1Parts) - len(v2Parts)
}


func (df *DeviceFact) GetDeviceOrLocal() string {
	if s.ToLower(df.Auth.OsName) == "macos" {
		return "local"
	}
	return "device"
}

func (df *DeviceFact) AssignOutPut(passkeyType string, matchProbability float64) {

	df.Output = Output{
		MatchPassKeyType: passkeyType,
		MatchProbability: matchProbability,
	}
}

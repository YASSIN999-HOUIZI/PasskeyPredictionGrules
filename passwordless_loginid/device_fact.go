package main

import (
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
	MatchPassKeyType string // PassKeyType
	MatchProbability float64
}

type DeviceData struct {
	Auth               DeviceInfo
	UserPasskeyHistory []*UserPasskeyHistory
	DeviceFeatures     DeviceFeature
	Output             Output
}

const LOCAL string = "local"

func (df *DeviceData) IsPassKeyExisting() bool {
	return df.DeviceFeatures.IsUVPPA
}

func (df *DeviceData) MatchOS() bool {
	for _, history := range df.UserPasskeyHistory {
		if df.Auth.OsName != history.DeviceInfo.OsName || df.Auth.OsVersion != history.DeviceInfo.OsVersion {
			return false
		}
	}
	return true
}

func (df *DeviceData) MatchDeviceProperties() bool {

	for _, deviceInfo := range df.UserPasskeyHistory {
		if deviceInfo.DeviceInfo.DeviceID != df.Auth.DeviceID || deviceInfo.DeviceInfo.DeviceSize != df.Auth.DeviceSize {
			return false
		}
		if s.ToLower(deviceInfo.PasskeyType) == LOCAL {
			if deviceInfo.DeviceInfo.ClientName != df.Auth.ClientName || deviceInfo.DeviceInfo.ClientVersion != df.Auth.ClientVersion {
				return false
			}
		}

	}

	return true
}

func (df *DeviceData) MatchCloudCompatibleBrowser() bool {
	return isVersionMatching(CloudClient{
		Platform:      s.ToLower(df.Auth.OsName),
		ClientName:    s.ToLower(df.Auth.ClientName),
		ClientVersion: df.Auth.ClientVersion,
	})
}

func (df *DeviceData) IsConditionalMediationAvailable() bool {
	return df.DeviceFeatures.IsCMA
}

func (df *DeviceData) IsCompatibleDevice() bool {
	return !s.EqualFold(df.Auth.OsName, "windows")
}

func (df *DeviceData) MatchProbability() float64 {
	probability := 0.0

	for _, pastDevice := range df.UserPasskeyHistory {
		weight := 0.0
		if df.Auth.DeviceID == pastDevice.DeviceInfo.DeviceID {
			probability = 100
			return probability
		}
		if df.Auth.OsName == pastDevice.DeviceInfo.OsName {
			weight = 0.9
		} else {
			weight = 0.8
		}
		switch pastDevice.PasskeyType {
		case "cloud":
			if isVersionMatching(CloudClient{
				Platform:      s.ToLower(df.Auth.OsName),
				ClientName:    s.ToLower(df.Auth.ClientName),
				ClientVersion: df.Auth.ClientVersion,
			}) {
				weight *= 0.8
			} else {
				weight *= 0.05
			}
		case LOCAL:
			if df.Auth.ClientName == pastDevice.DeviceInfo.ClientName {
				versionDiff := compareVersions(df.Auth.ClientVersion, pastDevice.DeviceInfo.ClientVersion)
				if versionDiff > 0 {
					weight *= 0.9 // upgrade most likely to happen
				} else if versionDiff < 0 { // downgrade less likely to happen
					weight *= 0.75
				}
			}
		default:
			versionDiff := compareVersions(df.Auth.OsVersion, pastDevice.DeviceInfo.OsVersion)
			if versionDiff > 0 {
				weight *= 0.8 // upgrade
			} else if versionDiff < 0 {
				weight *= 0.65 // downgrade
			}
			if df.Auth.DeviceSize != pastDevice.DeviceInfo.DeviceSize { // penalty if the device size is different
				weight *= 0.65
			}
		}

		newProbability := weight * 100
		if newProbability > probability { // We keep the greatest probability knowing entries probability are independent
			probability = newProbability
		}
	}

	return min(probability, 90)
}

func compareVersions(v1, v2 string) int {
	v1Parts := s.Split(v1, ".")
	v2Parts := s.Split(v2, ".")

	minLength := len(v1Parts)
	if len(v2Parts) < minLength {
		minLength = len(v2Parts)
	}

	for i := 0; i < minLength; i++ {
		v1Num, err := strconv.Atoi(v1Parts[i])
		if err != nil {
			return 0
		}
		v2Num, err := strconv.Atoi(v2Parts[i])
		if err != nil {
			return 0
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

func (df *DeviceData) GetDeviceOrLocal() string {

	if s.EqualFold(df.Auth.OsName, "macos") {
		return LOCAL
	}
	return "device"
}

func (df *DeviceData) AssignOutPut(passkeyType string, matchProbability float64) {

	df.Output = Output{
		MatchPassKeyType: passkeyType,
		MatchProbability: matchProbability,
	}
}

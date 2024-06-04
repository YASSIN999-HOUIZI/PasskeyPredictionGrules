package main

import (
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
	PasskeyType PassKeyType
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
		if deviceInfo.DeviceInfo != df.Auth {
			return false
		}
	}
	return true
}

func (df *DeviceFact) MatchCloudCompatibleBrowser() bool {
	return isVersionMatching(CloudClient{
		Platform: df.Auth.OsName,
		ClientName:    df.Auth.ClientName,
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
	if df.Output.MatchPassKeyType == "" {
		panic("no passkey type")
	}
	return 0.0
}

// func (df *DeviceFact) MatchPasskeyType() float64 {
// 	return 0.0
// }

// func (df *DeviceInfo) IsCloud() bool {
// 	return df.DeviceFeature.PasskeyType == "cloud"
// }

func (df *DeviceFact) AssignOutPut(passkeyType string, matchProbability float64) {

	df.Output = Output{
		MatchPassKeyType: passkeyType,
		MatchProbability: matchProbability,
	}
	// df.Output.MatchPassKeyType = passkeyType //PassKeyType(passkeyType)
}

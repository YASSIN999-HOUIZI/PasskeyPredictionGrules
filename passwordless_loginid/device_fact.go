package main

import (
	"reflect"
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
		if deviceInfo.DeviceInfo != df.Auth {
			return false
		}
	}
	return true
}

func (df *DeviceFact) MatchCloudCompatibleBrowser() bool {
	return isVersionMatching(CloudClient{
		Platform: df.Auth.OsName,
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
	if df.Output.MatchPassKeyType == "" {
		panic("no passkey type")
	}
	if len(df.UserPasskeyHistory) == 0 {
		return 0.0
	}
	// Count the number of times the passkey type appears with matching device info
	perfectMatches := make(map[string]float64)
	partialMatches := make(map[string]float64) // Track partial matches for each type
	passkeyTypeEntriesLength := 0
	for _, entry := range df.UserPasskeyHistory {
		passkeyType := entry.PasskeyType
		entry_device_info := entry.DeviceInfo

		if passkeyType == df.Output.MatchPassKeyType {
			propMatchCount := 0
			passkeyTypeEntriesLength++
			t := reflect.TypeOf(df.Auth)
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				v1 := reflect.ValueOf(df.Auth).FieldByName(f.Name).Interface()
				v2 := reflect.ValueOf(entry_device_info).FieldByName(f.Name).Interface()
				if v1 == v2 {
					propMatchCount++
				}
			}
			// Count full matches
			if propMatchCount == t.NumField() {
				perfectMatches[df.Output.MatchPassKeyType]++
			} else if propMatchCount > 0 {
				partialMatches[df.Output.MatchPassKeyType] += (float64(propMatchCount) * 0.142857)
			}
		}
	}
	// fmt.Println("count ",passkeyTypeEntriesLength)

	if passkeyTypeEntriesLength == 0 {
		return 0.0 * 100
	}
	if _, ok := perfectMatches[df.Output.MatchPassKeyType]; ok {
		return 100
	} else {
		return partialMatches[df.Output.MatchPassKeyType] * (1 / float64(passkeyTypeEntriesLength)) * 100
	}

}

// func (df *DeviceFact) MatchPasskeyType() float64 {
// 	return 0.0
// }

// func (df *DeviceInfo) IsCloud() bool {
// 	return df.DeviceFeature.PasskeyType == "cloud"
// }
func (df *DeviceFact) GetDeviceOrLocal() string{
	if df.Auth.OsName == "macOS"{
		return "local"
	}
	return "device"
}

func (df *DeviceFact) AssignOutPut(passkeyType string, matchProbability float64) {

	df.Output = Output{
		MatchPassKeyType: passkeyType,
		MatchProbability: matchProbability,
	}
	// df.Output.MatchPassKeyType = passkeyType //PassKeyType(passkeyType)
}

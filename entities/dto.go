package entities

type Token struct{
	Jwt string `json:"jwt"`
}

type User struct {
	TwitchId        string `json:"twitchId"`
	DisplayName     string `json:"displayName"`
	NameColor       string `json:"nameColor"`
	ProfileImageUrl string`json:"profileImageUrl"`
}


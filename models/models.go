package models

import "encoding/json"

type RedisStruct struct {
	Link       string
	ShrinkLink string
}
type Input_link struct {
	Link string `form:"link",json:"link"`
}

func MarshalBinary(link RedisStruct) ([]byte, error) {
	return json.Marshal(link)
}
func UnmarshalBinary(data []byte, Link RedisStruct) (RedisStruct, error) {
	return Link, json.Unmarshal(data, &Link)
}

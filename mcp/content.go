package mcp

import (
	"encoding/json"
	"errors"
)

func NewContents(message json.RawMessage) []Content {
	var contents []Content
	err := unmarshalContents(Result(message), &contents)
	if err != nil {
		panic(err)
	}

	return contents
}

func marshalContents(c *[]Content) (Result, error) {
	v := map[string]any{
		"contents": c,
	}
	return json.Marshal(v)
}

func unmarshalContents(result Result, c *[]Content) error {
	v := map[string]any{
		"contents": []map[string]json.RawMessage{},
	}
	err := json.Unmarshal(result, &v)
	if err != nil {
		return err
	}

	contents := v["contents"].([]map[string]json.RawMessage)

	var content Content
	for _, contentJson := range contents {
		mimeType, ok := contentJson["mimeType"]
		if !ok {
			return errors.New("missing mimeType field")
		}

		if contentJson["resource"] != nil {
			data, ok := contentJson["resource"]
			if !ok {
				return errors.New("missing resource field")
			}

			var resource Resource
			err := json.Unmarshal(data, &resource)
			if err != nil {
				return err
			}

			content = &ResourceContent{
				Resource: Resource{},
			}
		} else {
			content = &BinaryContent{
				MimeType: string(mimeType),
				Data:     contentJson["data"],
			}
		}

		*c = append(*c, content)
	}

	return nil
}

type Content interface {
	Type() string
}

var _ Content = (*BinaryContent)(nil)
var _ Content = (*ResourceContent)(nil)

type BinaryContent struct {
	MimeType string
	Data     []byte
}

func (b BinaryContent) Type() string {
	return b.MimeType
}

type ResourceContent struct {
	Resource Resource `json:"resource"`
}

func (e ResourceContent) Type() string {
	return "resource"
}

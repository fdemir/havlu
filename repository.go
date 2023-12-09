package main

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"
	"strings"
)

func GetAll(
	params url.Values,
	collection *[]interface{},
) []interface{} {
	result := []interface{}{}
	limit, _ := strconv.Atoi(params.Get("_limit"))

	if len(params) > 0 {
		count := 0
		for _, item := range *collection {
			hasLimitReached := limit > 0 && count >= limit

			if hasLimitReached {
				break
			}

			shouldAdd := true
			item := item.(map[string]interface{})

			for key, value := range params {

				if strings.HasPrefix(key, "_") {
					continue
				}

				if intValue, err := strconv.Atoi(value[0]); err == nil {
					if item[key] != intValue {
						shouldAdd = false
					}
				} else if boolValue, err := strconv.ParseBool(value[0]); err == nil {
					if item[key] != boolValue {
						shouldAdd = false
					}
				} else {
					if item[key] != value[0] {
						shouldAdd = false
					}
				}
			}

			if shouldAdd {
				result = append(result, item)
			}

			count += 1
		}
	} else {
		result = *collection
	}

	return result
}

func Create(
	payload io.Reader,
	collection *[]interface{},
) error {
	var body map[string]interface{}

	if err := json.NewDecoder(payload).Decode(&body); err != nil {
		return err
	}

	*collection = append(*collection, body)

	return nil
}

func Delete(
	id string,
	collection *[]interface{},
) []interface{} {
	parsedId, _ := strconv.Atoi(id)

	for index, item := range *collection {
		item := item.(map[string]interface{})
		if item["id"] == parsedId {
			*collection = append((*collection)[:index], (*collection)[index+1:]...)
			break
		}
	}

	return *collection
}

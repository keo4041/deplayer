package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func main() {
	//Open the input JSON file
	inputFile, err := os.Open("input.json")

	// Check if there is an error and print out
	if err != nil {
		fmt.Println(err)
	}

	// Defer the closing of the file to the end of the function
	defer inputFile.Close()

	// Read the file
	byteValue, err := ioutil.ReadAll(inputFile)

	// Handle reading error
	if err != nil {
		fmt.Println(err)
	}
	// Unmarshall byteValue into an interface
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	res := transformer(result)
	c, _ := json.MarshalIndent(res, "", "  ")
	fmt.Printf("%v", string(c))

}

func transformer(input map[string]interface{}) map[string]interface{} {
	output := map[string]interface{}{}
	for k, v := range input {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		data := v.(map[string]interface{})

		// STRING Checking
		if dataType, ok := data["S"]; ok {

			stringData := strings.TrimSpace(dataType.(string))
			if stringData != "" {
				t, e := time.Parse(time.RFC3339, stringData)
				if e == nil {
					output[k] = t.Unix()
				} else {
					output[k] = stringData
				}

			} else {
				continue
			}
		}
		// NUMERIC Checking
		if dataType, ok := data["N"]; ok {
			numData := strings.TrimSpace(dataType.(string))
			y, er := strconv.ParseFloat(numData, 64)
			if er == nil {
				output[k] = y
			} else {
				continue
			}
		}

		// BOOL Checking
		if dataType, ok := data["BOOL"]; ok {
			boolData := strings.TrimSpace(dataType.(string))
			if boolData != "" {
				switch strings.ToLower(boolData) {
				case "true", "t", "1":
					output[k] = true
				case "false", "f", "0":
					output[k] = false
				default:
					continue

				}
			} else {
				continue
			}
		}

		// NULL Checking
		if dataType, ok := data["NULL"]; ok {
			nullData := strings.TrimSpace(dataType.(string))
			if nullData != "" {
				switch strings.ToLower(nullData) {
				case "true", "t", "1":
					output[k] = nil
				default:
					fmt.Println("skip 1")
					continue

				}
			} else {
				continue
			}
		}

		// LIST Checking ////////////////////////////////
		if dataType, ok := data["L"]; ok {
			if reflect.TypeOf(dataType).Kind() == reflect.Slice {

				listData := dataType.([]interface{})
				// An slice that contains elements from different data types in not allowed in Golang
				// So I'm converting them from the data type to string before adding them.
				listSlice := make([]string, 0)
				for i := 0; i < len(listData); i++ {
					if reflect.TypeOf(listData[i]).Kind() == reflect.Map {
						listMap := listData[i].(map[string]interface{})
						if dataType, ok := listMap["N"]; ok {
							numericListData := strings.TrimSpace(dataType.(string))
							numericListDataString, er := strconv.ParseFloat(numericListData, 64)
							if er == nil {
								listSlice = append(listSlice, strconv.FormatFloat(numericListDataString, 'f', 0, 64))
							} else {
								continue
							}
						}

						if dataType, ok := listMap["S"]; ok {

							stringListData := strings.TrimSpace(dataType.(string))
							if stringListData != "" {
								stringListDataTimeFormat, e := time.Parse(time.RFC3339, stringListData)
								if e == nil {
									listSlice = append(listSlice, fmt.Sprintf("%d", stringListDataTimeFormat.Unix()))
								} else {
									listSlice = append(listSlice, stringListData)
								}

							} else {
								continue
							}
						}
						if dataType, ok := listMap["NULL"]; ok {
							nullListData := strings.TrimSpace(dataType.(string))
							if nullListData != "" {
								switch strings.ToLower(nullListData) {
								case "true", "t", "1":
									listSlice = append(listSlice, "null")
								default:
									continue

								}
							} else {
								continue
							}
						}
						if dataType, ok := listMap["BOOL"]; ok {
							boolListData := strings.TrimSpace(dataType.(string))
							if boolListData != "" {
								switch strings.ToLower(boolListData) {
								case "true", "t", "1":
									listSlice = append(listSlice, "true")
								case "false", "f", "0":
									listSlice = append(listSlice, "false")
								default:
									fmt.Println("skip 1")
									continue

								}
							} else {
								continue
							}
						}
					}
				}
				if len(listSlice) > 0 {
					output[k] = listSlice
				}

			} else {
				continue
			}

		}

		// MAP Checking
		if dataType, ok := data["M"]; ok {
			if reflect.TypeOf(dataType).Kind() == reflect.Map {
				r := dataType.(map[string]interface{})
				output[k] = transformer(r)
			} else {
				continue
			}
		}

	}
	return output

}

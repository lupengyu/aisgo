package handler

import (
	"errors"
	"fmt"
	"github.com/lupengyu/aisgo/helper"
	"github.com/lupengyu/aisgo/idl"
	"strconv"
	"strings"
)

func Decode(request *idl.DecodeRequest) (response *idl.DecodeResponse, err error) {
	if request == nil {
		return nil, errors.New("request error")
	}
	fmt.Println(request.Data)
	messages := strings.Split(request.Data, ",")
	if len(messages) < 7 {
		return nil, errors.New("data struct error")
	}
	if messages[0] != "!AIVDM" && messages[0] != "!AIVDO" {
		return nil, errors.New("data flag error")
	}
	message1, err := strconv.Atoi(messages[1])
	if err != nil {
		return nil, err
	} else if message1 <= 0 || message1 > 9 {
		return nil, errors.New("data message1 error")
	}
	message2, err := strconv.Atoi(messages[2])
	if err != nil {
		return nil, err
	} else if message2 <= 0 || message2 > message1 {
		return nil, errors.New("data message2 error")
	}
	message3, err := strconv.Atoi(messages[3])
	if err != nil {
		if messages[3] != "" {
			return nil, errors.New("data message3 error")
		}
	} else if message3 < 0 || message3 > 9 {
		return nil, errors.New("data message3 error")
	}

	if messages[4] != "A" && messages[4] != "B" && messages[4] != "" {
		return nil, errors.New("data message4 error")
	}

	bits := ""
	for _, v := range messages[5] {
		value := helper.AsciiToBits[string(v)]
		if value == "" {
			return nil, errors.New("data message5 error")
		}
		bits += value
	}
	if messages[6] != "" {
		check := strings.Split(messages[6], "*")
		if len(check) > 0 && check[0] != "0" {
			for _, v := range check[0] {
				value := helper.AsciiToBits[string(v)]
				if value == "" {
					return nil, errors.New("data message6 error")
				}
				bits += value
			}
		}
	}
	fmt.Println(bits, len(bits))

	resp := &idl.DecodeResponse{
		DecodeType: &idl.DecodeType{
			DataType: helper.BitsToNumbers(bits[0:6]),
			Index:    message2,
			Length:   message1,
			Bits:     bits,
		},
		Status: false,
		Time:   request.Time,
	}
	if request.PreDecodeType != nil {
		if request.PreDecodeType.Length == resp.DecodeType.Length &&
			request.PreDecodeType.Index == resp.DecodeType.Index-1 {
			resp.DecodeType.DataType = request.PreDecodeType.DataType
			resp.DecodeType.Bits = request.PreDecodeType.Bits + resp.DecodeType.Bits
		} else {
			return nil, errors.New("data error")
		}
	}

	fmt.Println("message type:", resp.DecodeType.DataType)

	if resp.DecodeType.DataType == 1 || resp.DecodeType.DataType == 2 || resp.DecodeType.DataType == 3 {
		// 1, 2, 3解码
		if resp.DecodeType.Index != resp.DecodeType.Length || resp.DecodeType.Index != 1 {
			return nil, errors.New("data error")
		}
		if len(resp.DecodeType.Bits) < 168 {
			return nil, errors.New("data error")
		}
		resp.Status = true
		resp.ParameterList = map[string]interface{}{
			"Repeat_Indicator":      helper.BitsToNumbers(resp.DecodeType.Bits[6:8]),
			"MMSI":                  helper.BitsToNumbers(resp.DecodeType.Bits[8:38]),
			"Navigation_Status":     helper.BitsToNumbers(resp.DecodeType.Bits[38:42]),
			"ROT":                   helper.BitsToComplementNumber(resp.DecodeType.Bits[42:50]),
			"SOG":                   float64(helper.BitsToNumbers(resp.DecodeType.Bits[50:60])) * 0.1,
			"Position_Accuracy":     helper.BitsToNumbers(resp.DecodeType.Bits[60:61]),
			"Longitude":             float64(helper.BitsToNumbers(resp.DecodeType.Bits[62:89])) / 600000.0,
			"Latitude":              float64(helper.BitsToNumbers(resp.DecodeType.Bits[90:116])) / 600000.0,
			"COG":                   float64(helper.BitsToNumbers(resp.DecodeType.Bits[116:128])) * 0.1,
			"HDG":                   helper.BitsToNumbers(resp.DecodeType.Bits[128:137]),
			"Time_Stamp":            helper.BitsToNumbers(resp.DecodeType.Bits[137:143]),
			"Reserved_for_regional": helper.BitsToNumbers(resp.DecodeType.Bits[143:145]),
			"RAIM_flag":             helper.BitsToNumbers(resp.DecodeType.Bits[148:149]),
		}
	} else if resp.DecodeType.DataType == 4 || resp.DecodeType.DataType == 11 {
		// 4, 11解码
		if resp.DecodeType.Index != resp.DecodeType.Length || resp.DecodeType.Index != 1 {
			return nil, errors.New("data error")
		}
		if len(resp.DecodeType.Bits) < 168 {
			return nil, errors.New("data error")
		}
		resp.Status = true
		resp.Time = &idl.Data{
			Year:   helper.BitsToNumbers(resp.DecodeType.Bits[38:52]),
			Month:  helper.BitsToNumbers(resp.DecodeType.Bits[52:56]),
			Day:    helper.BitsToNumbers(resp.DecodeType.Bits[56:61]),
			Hour:   helper.BitsToNumbers(resp.DecodeType.Bits[61:66]),
			Minute: helper.BitsToNumbers(resp.DecodeType.Bits[66:72]),
			Second: helper.BitsToNumbers(resp.DecodeType.Bits[72:78]),
		}
		resp.ParameterList = map[string]interface{}{
			"Repeat_Indicator":  helper.BitsToNumbers(resp.DecodeType.Bits[6:8]),
			"MMSI":              helper.BitsToNumbers(resp.DecodeType.Bits[8:38]),
			"UTC_year":          resp.Time.Year,
			"UTC_month":         resp.Time.Month,
			"UTC_day":           resp.Time.Day,
			"UTC_hour":          resp.Time.Hour,
			"UTC_min":           resp.Time.Minute,
			"UTC_second":        resp.Time.Second,
			"Position_Accuracy": helper.BitsToNumbers(resp.DecodeType.Bits[78:79]),
			"Longitude":         float64(helper.BitsToNumbers(resp.DecodeType.Bits[80:107])) / 600000.0,
			"Latitude":          float64(helper.BitsToNumbers(resp.DecodeType.Bits[108:134])) / 600000.0,
		}
	} else if resp.DecodeType.DataType == 5 {
		if resp.DecodeType.Length != 2 {
			return nil, errors.New("data error")
		}
		if resp.DecodeType.Index == resp.DecodeType.Length {
			// 5解码
			fmt.Println(len(resp.DecodeType.Bits))
			if len(resp.DecodeType.Bits) < 424 {
				return nil, errors.New("data error")
			}
			callSign, err := helper.Bits2Ascii(resp.DecodeType.Bits[70:112])
			if err != nil {
				return nil, err
			}
			name, err := helper.Bits2Ascii(resp.DecodeType.Bits[112:232])
			if err != nil {
				return nil, err
			}
			destination, err := helper.Bits2Ascii(resp.DecodeType.Bits[302:422])
			if err != nil {
				return nil, err
			}
			resp.Status = true
			size := resp.DecodeType.Bits[240:270]
			ETA := resp.DecodeType.Bits[274:294]
			resp.ParameterList = map[string]interface{}{
				"Repeat_Indicator": helper.BitsToNumbers(resp.DecodeType.Bits[6:8]),
				"MMSI":             helper.BitsToNumbers(resp.DecodeType.Bits[8:38]),
				"AIS":              helper.BitsToNumbers(resp.DecodeType.Bits[38:40]),
				"IMO":              helper.BitsToNumbers(resp.DecodeType.Bits[40:70]),
				"CallSign":         callSign,
				"Name":             name,
				"Type":             helper.BitsToNumbers(resp.DecodeType.Bits[232:240]),
				"A":                helper.BitsToNumbers(size[0:9]),
				"B":                helper.BitsToNumbers(size[9:18]),
				"C":                helper.BitsToNumbers(size[18:24]),
				"D":                helper.BitsToNumbers(size[24:30]),
				"Length":           helper.BitsToNumbers(size[0:9]) + helper.BitsToNumbers(size[9:18]),
				"Width":            helper.BitsToNumbers(size[18:24]) + helper.BitsToNumbers(size[24:30]),
				"Position":         helper.BitsToNumbers(resp.DecodeType.Bits[270:274]),
				"ETA_min":          helper.BitsToNumbers(ETA[14:20]),
				"ETA_hour":         helper.BitsToNumbers(ETA[9:14]),
				"ETA_day":          helper.BitsToNumbers(ETA[4:9]),
				"ETA_month":        helper.BitsToNumbers(ETA[0:4]),
				"Draught":          float64(helper.BitsToNumbers(resp.DecodeType.Bits[294:302])) * 0.1,
				"Destination":      destination,
			}
		}
	} else {
		return nil, errors.New("data error")
	}

	return resp, nil
}

package handler

import (
	"fmt"
	"github.com/lupengyu/aisgo/idl"
	"testing"
)

func Test_Message1_2_3(t *testing.T) {
	response, err := Decode(&idl.DecodeRequest{
		Data: "!AIVDM,1,1,0,B,16:b4BiP007WveP@taA000IN1P00,0*64",
	})
	fmt.Println(response, err)
	response, err = Decode(&idl.DecodeRequest{
		Data: "!AIVDM,1,1,,A,169FsD001o8ewMhF8Bb997A@05K8,0*26",
	})
	fmt.Println(response, err)
}

func Test_Message_4_11(t *testing.T) {
	response, err := Decode(&idl.DecodeRequest{
		Data: "!AIVDM,1,1,,,;028j:Qv;:cn<OvPlFFl:PQ00000,0*43",
	})
	fmt.Println(response, err)
}

func Test_Message_5(t *testing.T) {
	response, err := Decode(&idl.DecodeRequest{
		Data: "!AIVDM,2,1,2,,56:RPb00000094tP001`PtpLMDtP4T`T4r3;3<0o2hK4672c08n2@H3AC`,0*70",
	})
	fmt.Println(response, err)
	response, err = Decode(&idl.DecodeRequest{
		Data:          "!AIVDM,2,2,2,,888888888888;,2*6C",
		PreDecodeType: response.DecodeType,
	})
	fmt.Println(response, err)
}

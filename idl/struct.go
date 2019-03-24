package idl

type DecodeType struct {
	DataType int
	Index    int
	Length   int
	Bits     string
}

type DecodeRequest struct {
	Data          string
	PreDecodeType *DecodeType
	Time          *Data
}

type DecodeResponse struct {
	DecodeType    *DecodeType
	Status        bool
	ParameterList map[string]interface{}
	Time          *Data
}

type Data struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}

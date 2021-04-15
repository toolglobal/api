package types

import "encoding/json"

// CONTRACT: a zero Result is OK.
type Result struct {
	Code uint32 `json:"Code"`
	Data []byte `json:"Data"`
	Log  string `json:"Log"` // Can be non-deterministic
}

func (res Result) ToJSON() string {
	j, err := json.Marshal(res)
	if err != nil {
		return res.Log
	}
	return string(j)
}

func (res *Result) FromJSON(j string) *Result {
	err := json.Unmarshal([]byte(j), res)
	if err != nil {
		res.Code = CodeType_InternalError
		res.Log = j
	}
	return res
}

func (res Result) IsOK() bool {
	return res.Code == CodeType_OK
}

func (res Result) IsErr() bool {
	return res.Code != CodeType_OK
}

func (res Result) Error() string {
	// return fmt.Sprintf("{code:%v, data:%X, log:%v}", res.Code, res.Data, res.Log)
	return res.ToJSON()
}

func (res Result) String() string {
	// return fmt.Sprintf("{code:%v, data:%X, log:%v}", res.Code, res.Data, res.Log)
	return res.ToJSON()
}

func (res Result) PrependLog(log string) Result {
	return Result{
		Code: res.Code,
		Data: res.Data,
		Log:  log + ";" + res.Log,
	}
}

func (res Result) AppendLog(log string) Result {
	return Result{
		Code: res.Code,
		Data: res.Data,
		Log:  res.Log + ";" + log,
	}
}

func (res Result) SetLog(log string) Result {
	return Result{
		Code: res.Code,
		Data: res.Data,
		Log:  log,
	}
}

func (res Result) SetData(data []byte) Result {
	return Result{
		Code: res.Code,
		Data: data,
		Log:  res.Log,
	}
}

//----------------------------------------

// NOTE: if data == nil and log == "", same as zero Result.
func NewResultOK(data []byte, log string) Result {
	return Result{
		Code: CodeType_OK,
		Data: data,
		Log:  log,
	}
}

func NewError(code uint32, log string) Result {
	return Result{
		Code: code,
		Log:  log,
	}
}

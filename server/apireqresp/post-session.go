// package of Combat API requests and responses
package apireqresp

import (
	"encoding/base64"
	"encoding/json"
)

// the struct for sending by http
type ReqPostSession struct {
	SessionParams string
	ZipFileBase64 string
}

func NewReqPostSession(sessionParams string, ZipFile []byte) *ReqPostSession {
	zipFileBase64 := base64.StdEncoding.EncodeToString(ZipFile)
	return &ReqPostSession{
		SessionParams: sessionParams,
		ZipFileBase64: zipFileBase64,
	}
}

func (t *ReqPostSession) GetJson() ([]byte, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func ParseReqPostSessionFromBytes(bytes []byte) (ReqPostSession, error) {
	res := ReqPostSession{}
	err := json.Unmarshal(bytes, &res)
	return res, err
}

func (t *ReqPostSession) GetDecodedFile() ([]byte, error) {
	return base64.StdEncoding.DecodeString(t.ZipFileBase64)
}

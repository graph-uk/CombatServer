// package of Combat API requests and responses
package apireqresp

import (
	"encoding/base64"
	"encoding/json"
)

// the struct for sending by http
type ReqPostCaseResult struct {
	CaseID        string
	ExitStatus    string
	StdOut        string
	ZipFileBase64 string
}

func NewReqPostCaseResult(caseID, exitStatus, stdOut string, zipFile []byte) *ReqPostCaseResult {
	zipFileBase64 := base64.StdEncoding.EncodeToString(zipFile)
	return &ReqPostCaseResult{
		CaseID:        caseID,
		ExitStatus:    exitStatus,
		StdOut:        stdOut,
		ZipFileBase64: zipFileBase64,
	}
}

func (t *ReqPostCaseResult) GetJson() ([]byte, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func ParseReqPostCaseResultFromBytes(bytes []byte) (ReqPostCaseResult, error) {
	res := ReqPostCaseResult{}
	err := json.Unmarshal(bytes, &res)
	return res, err
}

func (t *ReqPostCaseResult) GetDecodedFile() ([]byte, error) {
	return base64.StdEncoding.DecodeString(t.ZipFileBase64)
}

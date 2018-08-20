// package of Combat API requests and responses
package apireqresp

import (
	"encoding/base64"
	"encoding/json"
)

// the struct for sending by http
type ResGetJob struct {
	CaseID        string
	CaseCMD       string
	ZipFileBase64 string
}

func NewResGetJob(caseID, caseCMD string, ZipFile []byte) *ResGetJob {

	zipFileBase64 := base64.StdEncoding.EncodeToString(ZipFile)

	return &ResGetJob{
		CaseID:        caseID,
		CaseCMD:       caseCMD,
		ZipFileBase64: zipFileBase64,
	}
}

func ParseResGetJobFromBytes(bytes []byte) (ResGetJob, error) {
	res := ResGetJob{}
	err := json.Unmarshal(bytes, &res)
	return res, err
}

func (t *ResGetJob) GetJson() ([]byte, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func (t *ResGetJob) GetDecodedFile() ([]byte, error) {
	return base64.StdEncoding.DecodeString(t.ZipFileBase64)
}

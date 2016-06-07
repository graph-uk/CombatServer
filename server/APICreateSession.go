package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (t *CombatServer) createSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionName := strconv.FormatInt(time.Now().UnixNano(), 10)

	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	sessionParams := r.FormValue("params")
	if sessionParams == "" {
		fmt.Println("cannot extract session params")
		return
	}

	os.MkdirAll("./sessions/"+sessionName, 0777)
	f, err := os.OpenFile("./sessions/"+sessionName+"/archived.zip", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	req, err := t.mdb.DB.Prepare("INSERT INTO Sessions(id,params) VALUES(?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = req.Exec(sessionName, sessionParams)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	io.WriteString(w, sessionName)
	fmt.Println(r.Host + " Create new session: " + sessionName + " " + sessionParams)
}

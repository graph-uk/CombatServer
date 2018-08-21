package server

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/graph-uk/combat-server/server/session"
)

func (t *CombatServer) getSessionStatusTemplate() *string {
	template := `
<!DOCTYPE html>
<html lang="en-US">

<head>
    <title>Session: {{.ID}}</title>
    <link href="/bindata/jquery.bxslider/jquery.bxslider.css" rel="stylesheet">
    <style>
        .rTable {
            display: table;
            width: 100%;
        }
        
        .rTableRow {
            display: table-row;
        }
        
        .rTableStatusGreen {
            background-color: green;
        }
        
        .rTableStatusRed {
            background-color: red;
        }
        
        .rTableStatusProgress {
            background-color: #99ff99;
        }
        
        .rTableHeading {
            background-color: #ddd;
            display: table-header-group;
        }
        
        .rTableStatusCell {
            display: table-cell;
            padding: 3px 10px;
            border: 1px solid #999999;
            width: 10px
        }
        
        .rTableCell,
        .rTableHead {
            display: table-cell;
            padding: 3px 10px;
            border: 1px solid #999999;
        }
        
        .rTableHeading {
            display: table-header-group;
            background-color: #ddd;
            font-weight: bold;
        }
        
        .rTableFoot {
            display: table-footer-group;
            font-weight: bold;
            background-color: #ddd;
        }
        
        .rTableBody {
            display: table-row-group;
        }
        
        .input-button {
            width: 100%;
            text-align: left;
            background-color: #FFffff;
            border-radius: 10px;
            -moz-border-radius: 10px;
            -webkit-border-radius: 10px;
            border: 1px solid #ccc;
            font-weight: bolder;
            color: #000;
            margin: 0 auto;
            padding: 5px;
        }
        
        .spoil {}
        
        .smallfont {}
        
        .alt2 {}
		
    </style>
	<script src="/bindata/jquery-2.1.4.min.js"></script>

    <script type="text/javascript">
        function Spoil(tryID) {
            if (document.getElementById(tryID).style.display != '') {
                document.getElementById(tryID).style.display = '';
            } else {
                document.getElementById(tryID).style.display = 'none';
            }
			
			if (typeof window.inited_slieders == 'undefined') {
   				window.inited_slieders = new Array();
			}
        
			a = window.inited_slieders.indexOf(tryID);
			if (a == -1) {
				window.inited_slieders.push(tryID);
				$('.slider2'+tryID).bxSlider({
			    	slideWidth: 650,
			    	minSlides: 1,
			    	maxSlides: 1
			  	});
			}

		}
    </script>
	
</head>

<body>
    <h2>Session: {{.ID}}</h2>
    <div class="rTable">
        <div class="rTableHeading">
            <div class="rTableRow">
                <div class="rTableCell">State</div>
                <div class="rTableCell">Details</div>
            </div>
        </div>
        <div class="rTableBody">
			{{range .Cases}}
				<div class="rTableRow">
					{{if eq .InProgress true}}
					    <div class="rTableStatusCell rTableStatusProgress"></div>
					{{else}}
						{{if eq .Finished true}}
							{{if eq .Passed true}}
								<div class="rTableStatusCell rTableStatusGreen"></div>
							{{else}}
								<div class="rTableStatusCell rTableStatusRed"></div>
							{{end}}
						{{else}}
							<div class="rTableStatusCell"></div>
						{{end}}	
					{{end}}
					
					
			        <div class="rTableCell">
						{{.CMDLine}}
						{{range .Tries}}
							{{$tryID := .ID}}
							<div class="smallfont">
								<input type="button" value="Try" ; class="input-button" onclick="Spoil('{{.ID}}')">
		                    </div>
		                    <div class="alt2">
		                        <div id="{{.ID}}" style="display:none">
									<div class="rTableRow">
										{{if ne (len .Screens) 0}}
											<div class="rTableCell" style="width: 650px">
												<div class="slider2{{.ID}}" style="float: left;">
												{{range .Screens}}
													<div class="slide">
														<span><a href="/tries/{{html $tryID}}/out/{{.ID}}.html">PageSource</a></span><br>
														<span>URL: {{.URL}}</span>
														<img src="/tries/{{html $tryID}}/out/{{.ID}}.png">
													</div>
												{{end}}
												</div>		
											</div>
										{{end}}
										<div class="rTableCell" style="vertical-align:top;">
											<span>
												{{range .STDOut}}
													{{.}}<br>
				                            	{{end}}
											</span>
										</div>
									</div>
		                        </div>
		                    </div>
						{{end}}
					</div>
				</div>
	        {{end}}
    	</div>
    </div>
							

	<script src="/bindata/jquery.bxslider/jquery.bxslider.min.js"></script>
	
	

</body>

</html>	
	`
	return &template
}

func (t *CombatServer) pageSessionStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		path := strings.Split(r.URL.Path, "/")
		sessionID := strings.TrimSpace(path[len(path)-1])
		if strings.TrimSpace(sessionID) == "" {
			//			w.Write([]byte("Session ID is not specified. Please, provide session ID like: /sessions/11203487203498"))
			t.showSessionsPage(w)
			return
		}

		// Get session status
		Page_session, err := session.NewAssignedSession(sessionID, &t.mdb, t.startPath)
		if err != nil {
			w.Write([]byte("Error: " + err.Error() + "<br>\n"))
			return
		}

		PS_session, err := Page_session.GetSessionPageStruct()
		if err != nil {
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		// Create a template.
		pageBuffer := new(bytes.Buffer)

		tt, err := template.New("SessionPage").Parse(*t.getSessionStatusTemplate())
		if err != nil {
			w.Write([]byte("Error: " + err.Error() + "<br>\n"))
			return
		}

		err = tt.Execute(pageBuffer, PS_session)
		if err != nil {
			w.Write([]byte("Error: " + err.Error() + "<br>\n"))
			return
		}

		w.Write(pageBuffer.Bytes())
	}
}

func (t *CombatServer) getSessionsStatusTemplate() *string {
	template := `
<!DOCTYPE html>
<html lang="en-US">

<head>
    <title>{{.ProjectName}}</title>
    <style>
        .rTable {
            display: table;
            width: 100%;
        }
        
        .rTableRow {
            display: table-row;
        }
        
        .rTableStatusGreen {
            background-color: green;
        }
        
        .rTableStatusRed {
            background-color: red;
        }
        
        .rTableStatusProgress {
            background-color: #99ff99;
        }
        
        .rTableHeading {
            background-color: #ddd;
            display: table-header-group;
        }
        
        .rTableStatusCell {
            display: table-cell;
            padding: 3px 10px;
            border: 1px solid #999999;
            width: 10px
        }
        
        .rTableCell,
        .rTableHead {
            display: table-cell;
            padding: 3px 10px;
            border: 1px solid #999999;
        }
        
        .rTableHeading {
            display: table-header-group;
            background-color: #ddd;
            font-weight: bold;
        }
        
        .rTableFoot {
            display: table-footer-group;
            font-weight: bold;
            background-color: #ddd;
        }
        
        .rTableBody {
            display: table-row-group;
        }
        
        .input-button {
            width: 100%;
            text-align: left;
            background-color: #FFffff;
            border-radius: 10px;
            -moz-border-radius: 10px;
            -webkit-border-radius: 10px;
            border: 1px solid #ccc;
            font-weight: bolder;
            color: #000;
            margin: 0 auto;
            padding: 5px;
        }
        
        .spoil {}
        
        .smallfont {}
        
        .alt2 {}
		
    </style>
</head>

<body>
    <h2>{{.ProjectName}}</h2>
    <div class="rTable">
        <div class="rTableHeading">
            <div class="rTableRow">
                <div class="rTableCell">Sessions</div>
            </div>
        </div>
        <div class="rTableBody">
			{{range .Sessions}}
				<div class="rTableRow">
			        <div class="rTableCell">
						<a href="/sessions/{{.ID}}">{{.ID}}</a>
					</div>
				</div>
	        {{end}}
    	</div>
    </div>					
</body>

</html>	
	`
	return &template
}

type PS_testSessions struct {
	ProjectName string
	Sessions    []*PS_session
}

type PS_session struct { // testCase struct for
	ID string
}

func (t *CombatServer) getSessionsPageStruct() (*PS_testSessions, error) {
	var result PS_testSessions
	result.ProjectName = t.config.ProjectName

	req, err := t.mdb.DB.DB().Prepare(`SELECT id FROM Sessions ORDER BY id DESC`)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	rows, err := req.Query()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	for rows.Next() {
		var curSession PS_session
		err := rows.Scan(&curSession.ID)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		result.Sessions = append(result.Sessions, &curSession)
	}
	rows.Close()

	return &result, nil
}

func (t *CombatServer) showSessionsPage(w http.ResponseWriter) {
	pageStruct, err := t.getSessionsPageStruct()
	if err != nil {
		w.Write([]byte("Something wrong. See more in log."))
		fmt.Println(err.Error())
		return
	}

	// Create a template.
	pageBuffer := new(bytes.Buffer)

	tt, err := template.New("SessionPage").Parse(*t.getSessionsStatusTemplate())
	if err != nil {
		w.Write([]byte("Error: " + err.Error() + "<br>\n"))
		return
	}

	err = tt.Execute(pageBuffer, pageStruct)
	if err != nil {
		w.Write([]byte("Error: " + err.Error() + "<br>\n"))
		return
	}

	w.Write(pageBuffer.Bytes())
}

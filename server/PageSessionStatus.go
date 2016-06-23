package server

import (
	"net/http"
	//	"strconv"
	//"encoding/json"
	"bytes"
	"html/template"
	"strings"

	"github.com/graph-uk/combat-server/server/session"
)

//type SessionStatus struct {
//	Finished           bool
//	//TotalCasesCount    int
//	//FinishedCasesCount int
//	//FailReports        []string
//}

func (t *CombatServer) getSessionStatusTemplate() *string {
	template := `
<!DOCTYPE html>
<html lang="en-US">

<head>
    <title>Session: {{.ID}}</title>
    <link rel="stylesheet" href="html.css">
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


	<link rel="stylesheet" href="/tries/unslider/dist/css/unslider.css">

    <script type="text/javascript">
        function Spoil(tryID) {
            if (document.getElementById(tryID).style.display != '') {
                document.getElementById(tryID).style.display = '';
            } else {
                document.getElementById(tryID).style.display = 'none';
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
						<div class="smallfont"><input type="button" value="Try" ; class="input-button" onclick="Spoil('{{.ID}}')" />
	                    </div>
	                    <div class="alt2">
	                        <div id="{{.ID}}" style="display: none;">
								<div class="ScreenSlider">
									<ul>
										{{range .Screens}}
											<li><img src="/tries/{{html $tryID}}/out/{{.}}.png"></li>
										{{end}}
									</ul>
								</div>

								{{range .STDOut}}
									{{.}}<br>
	                            {{end}}
	                        </div>
	                    </div>
					{{end}}
				</div>
			</div>
        {{end}}
    </div>
    </div>
	<script src="//code.jquery.com/jquery-2.1.4.min.js"></script>
	<script src="/tries/unslider/src/js/unslider.js"></script>
		<script>
		jQuery(document).ready(function($) {
			$('.ScreenSlider').unslider();
		});
	</script>
</body>

</html>	
	`
	//	template2 := `hello {{.UserName}}!`
	//	return &template2

	return &template
}

func (t *CombatServer) pageSessionStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		path := strings.Split(r.URL.Path, "/")
		sessionID := strings.TrimSpace(path[len(path)-1])
		if strings.TrimSpace(sessionID) == "" {
			w.Write([]byte("Session ID is not specified. Please, provide session ID like: /sessions/11203487203498"))
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
			w.Write([]byte("Error: " + err.Error() + "<br>\n"))
			return
		}

		// Create a template.
		pageBuffer := new(bytes.Buffer)

		tt, err := template.New("SessionPage").Parse(*t.getSessionStatusTemplate())
		if err != nil {
			w.Write([]byte("Error: " + err.Error() + "<br>\n"))
			return
		}
		//println(PS_session.ID)

		err = tt.Execute(pageBuffer, PS_session)
		if err != nil {
			w.Write([]byte("Error: " + err.Error() + "<br>\n"))
			return
		}

		w.Write(pageBuffer.Bytes())
	}
}

const combat = window.combat = window.combat || {};
const $app = document.querySelector('#app');


combat.createTag = (tagName, attrs) => {
	const tag = document.createElement(tagName);

	if (typeof attrs === 'object') {
		Object.keys(attrs).forEach(key => {
			const value = attrs[key];

			if (key !== 'children') {
				return value !== '' && tag.setAttribute(key, value);
			}

			if (typeof value === 'string' || typeof value === 'number') {
				tag.innerHTML = value;
			} else if (typeof value === 'object') {
				if (value.constructor === Array) {
					tag.append(...value);
				} else {
					tag.append(value);
				}
			}
		});
	}

	return tag;
}

if (typeof window.combatLogs === 'object') {
	console.log(window.combatLogs);
	setTimeout(() => combat.renderTable($app, window.combatLogs), 0);
}

// window.onload = function() {
//     var xhttp = new XMLHttpRequest();
//     xhttp.onreadystatechange = function() {
//         if (this.readyState === 4 && this.status === 200) {
//             var notificationParameters = JSON.parse(this.responseText);
//             setInitialStatus(notificationParameters);
//         }
//     };
//     xhttp.open("GET", "/api/v1/config");
//     xhttp.send();
// };

// window.onload= () =>{
//     http.get('/api/v1/config', (resp) => {
//         let data = '';
//         // A chunk of data has been recieved.
//         resp.on('data', (chunk) => {
//             data += chunk;
//         });
//         // The whole response has been received. Print out the result.
//         resp.on('end', () => {
//            setInitialStatus(data)
//         });
//
//     }).on("error", (err) => {
//         console.log("Error: " + err.message);
//     });
// }

document.addEventListener("DOMContentLoaded", (event) =>{
    var notificationParameters = getNotificationParameters()
});


function getNotificationParameters() {
    fetch(`/api/v1/config`, {
        method: "GET"
    })
        .then(resp => (setInitialStatus(JSON.parse(resp))))
        .catch(()=>console.log("GET failed"))
}

function setInitialStatus(notificationParameters) {
    var notificationButton = document.getElementById("notification");
	var status = Boolean(notificationParameters.NotificationEnabled);
   	if (status){
        notificationButton.setAttribute("class", "notificationEnabled");
        notificationButton.innerHTML= "Disable notification";

    }
   	else{
        notificationButton.setAttribute("class", "notificationDisabled")
        notificationButton.innerHTML= "Enable notification";

    }
}

// function changeNotificationStatus(notificationButton) {
//     var xhttp = new XMLHttpRequest();
//     xhttp.onreadystatechange = function() {
//         if (this.readyState === 4 && this.status === 200) {
//             var notificationParameters = JSON.parse(this.responseText);
//             if (Boolean(notificationParameters[1])===true){
//                 xhttp.open("PUT", "/api/v1/config")
//                 xhttp.send( JSON.stringify({"MuteTimestamp":"null", "MutedEnable":"null"}))
//                 notificationButton.setAttribute("class", "notificationEnabled")            }
//             else{
//                 if (Boolean(notificationParameters[1])===true){
//                     xhttp.open("PUT", "/api/v1/config")
//                     xhttp.send( JSON.stringify({"MuteTimestamp":"null", "MutedEnable":"true"}))
//                     notificationButton.setAttribute("class", "notificationDisabled")
//                 }
//             }
//         }
//     };
//     xhttp.open("GET", "/api/v1/config");
//     xhttp.send();
// }

// function changeNotificationStatus(notificationButton) {
//     xhttp.onreadystatechange = function() {
//         if (this.readyState === 4 && this.status === 200) {
//             var notificationParameters = JSON.parse(this.responseText);
//             if (Boolean(notificationParameters[1])===true){
//                 xhttp.open("PUT", "/api/v1/config")
//                 xhttp.send( JSON.stringify({"MuteTimestamp":"null", "MutedEnable":"null"}))
//                 notificationButton.setAttribute("class", "notificationEnabled")            }
//             else{
//                 if (Boolean(notificationParameters[1])===true){
//                     xhttp.open("PUT", "/api/v1/config")
//                     xhttp.send( JSON.stringify({"MuteTimestamp":"null", "MutedEnable":"true"}))
//                     notificationButton.setAttribute("class", "notificationDisabled")
//                 }
//             }
//         }
//     };
//     xhttp.open("GET", "/api/v1/config");
//     xhttp.send();
// }
var notificationButton = document.getElementById("notification");

function changeNotificationStatus() {
    var notificationParameters = getNotificationParameters();
    var data = notificationParameters.NotificationEnabled? {"MuteTimestamp":"null", "NotificationEnabled":"null"}:{"MuteTimestamp":"null", "NotificationEnabled":"true"}
    fetch(`/api/v1/config`, {
        method: "PUT",
        body: JSON.stringify(data)
    })
        .then(() => {
        	if(notificationParameters.NotificationEnabled){

			 notificationButton.setAttribute("class", "notificationDisabled");
			 notificationButton.innerHTML= "Enable notification";

        	}
        	else {
        		notificationButton.setAttribute("class", "notificationEnabled");
                notificationButton.innerHTML= "Disable notification";

            }})
        .catch(()=>console.log("GET failed"))
}

notificationButton.addEventListener("click", function () {
    changeNotificationStatus();
});

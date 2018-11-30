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





function setInitialStatus(notificationParameters) {
    console.log("setting status")
	console.log(notificationParameters)
    var notificationButton = document.getElementById("notification");
    var status = notificationParameters.NotificationEnabled;
    console.log(status)
    if (status){
        notificationButton.setAttribute("class", "notificationEnabled");
        notificationButton.innerHTML= "Disable notification";

    }
    else{
        notificationButton.setAttribute("class", "notificationDisabled")
        notificationButton.innerHTML= "Enable notification";

    }
}

document.addEventListener("DOMContentLoaded", (event) =>{
	console.log("here")
    fetch(`/api/v1/config`, {
        method: "GET"
    })
        .then((resp) => (resp.json().then((data) => {
            setInitialStatus(data)
        })))
        .catch((error)=>console.log("GET failed " + error))
});


function getNotificationParameters() {
   fetch(`/api/v1/config`, {
        method: "GET"
    })
       .then((resp) => (resp.json().then((data) => {
           changeNotificationStatus(data)
       })))
        .catch((error)=>console.log("GET failed " + error))
}

var notificationButton = document.getElementById("notification");

function changeNotificationStatus(notificationParameters) {
        var data = notificationParameters.NotificationEnabled? {"MuteTimestamp":null, "NotificationEnabled":false}:{"MuteTimestamp":null, "NotificationEnabled":true}
    fetch(`/api/v1/config`, {
        method: "PUT",
		headers: {'Content-Type': 'application/json'},
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
    getNotificationParameters()
});

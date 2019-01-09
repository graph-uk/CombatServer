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


var theTimer=null;

function setInitialStatus(notificationParameters) {
    var notificationButton = document.getElementById("notification");
    var status = notificationParameters.NotificationEnabled;
    if (status){
        console.log("Setting initial status: notification enabled");
        notificationButton.setAttribute("class", "notificationEnabled");
        notificationButton.innerHTML= "Disable notification";
    }
    else{
        clearInterval(theTimer);
        showTimer(notificationParameters);
        console.log("Setting initial status: notification disabled");
        notificationButton.setAttribute("class", "notificationDisabled");
        notificationButton.innerHTML= "Enable notification";

    }
}

document.addEventListener("DOMContentLoaded", (event) =>{
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
           changeNotificationStatus(data);
       })))
        .catch((error)=>console.log("GET failed " + error))
}


var notificationButton = document.getElementById("notification");

function changeNotificationStatus(notificationParameters) {
    let data = notificationParameters.NotificationEnabled? {"MuteTimestamp":null, "NotificationEnabled":false}:{"MuteTimestamp":null, "NotificationEnabled":true}
    fetch(`/api/v1/config`, {
        method: "PUT",
		headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(data)
    })
        .then(() => {
            fetch(`/api/v1/config`, {
                method: "GET"
            })
                .then((resp) => (resp.json().then((data) => {
                    if(!data.NotificationEnabled){
                        console.log("Changing status: notification disabled");
                        notificationButton.setAttribute("class", "notificationDisabled");
                        notificationButton.innerHTML= "Enable notification";
                        clearInterval(theTimer);
                        showTimer(data);
                    }
                    else {
                        notificationButton.setAttribute("class", "notificationEnabled");
                        notificationButton.innerHTML= "Disable notification";
                        var display = document.querySelector('#time');
                        console.log("Changing status: notification enabled");
                        clearInterval(theTimer);
                        display.innerHTML = "You have enabled notification.";

                    }
                })))
                .catch((error)=>console.log("GET failed " + error))

        })
        .catch(()=>console.log("GET failed"))
}

notificationButton.addEventListener("click", function () {
    getNotificationParameters()
});


var endTime, now, hours, minutes, seconds, duration;

function showTimer(notificationParameters){
    console.log("Showing the timer");
        let display = document.querySelector('#time');
        theTimer = setInterval(function() {
            endTime = (Date.parse(notificationParameters.MuteTimestamp)) + notificationParameters.MuteDurationMinutes * 60 * 1000;
            console.log(notificationParameters.MuteDurationMinutes + " is the duration in minutes");
            console.log(endTime + " end time");
            now = new Date().getTime();
            console.log(now + " now");
            duration = (endTime - now) / 1000;
            console.log(duration + " should be the duration");

            hours = parseInt(duration / 3600, 10);
            minutes = Math.floor((duration %= 3600) / 60);
            seconds = parseInt(duration % 60, 10);

            hours = hours < 10 ? "0" + hours : hours;
            minutes = minutes < 10 ? "0" + minutes : minutes;
            seconds = seconds < 10 ? "0" + seconds : seconds;
            display.innerHTML = "Muted for " + notificationParameters.MuteDurationMinutes + " minutes. \n Time left : " + hours + "h " + minutes + "m " + seconds + "s";
            console.log(duration + "is the duration");
            if(duration <=0){
                clearInterval(theTimer);
                notificationButton.setAttribute("class", "notificationEnabled");
                notificationButton.innerHTML= "Disable notification";
                display.innerHTML = "Notification was enabled.";
            }
        }, 1000);
}






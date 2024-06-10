pcSignal = {}
pcSender = {}
pcRecievers = {}
meetingId = "" 
userId = ""


// use http://localhost:4200/call;meetingId=07927fc8-af0a-11ea-b338-064f26a5f90a;userId=alice;peerId=bob
// and http://localhost:4200/call;meetingId=07927fc8-af0a-11ea-b338-064f26a5f90a;userId=bob;peerId=alice
// start the call



function initPeers() {



    meetingId = document.getElementById('mid').value;


    userId = document.getElementById('uid').value;


    if (meetingId == ""){

        alert("feed meeting ID!")

        return

    }

    if (userId == ""){

        alert("feed user ID!")

        return

    }


/*
    pcSender = new RTCPeerConnection({
    iceServers: [
        {
        urls: 'stun:stun.l.google.com:19302'
        }
    ]
    })
    pcReciever = new RTCPeerConnection({
    iceServers: [
        {
        urls: 'stun:stun.l.google.com:19302'
        }
    ]
    })
*/


    pcSender = new RTCPeerConnection({
        iceServers: [
            {
            urls: 'stun:localhost:3478'
            }
        ]
    })


    pcSender.onicecandidate = async function(event) {
        if (event.candidate === null) {

            console.log("sender ice")

            let resp = await axios.post("/peers/room/sdp/m/" + meetingId + "/c/"+ userId + "/s/" + true,
            {
                "sdp" : btoa(JSON.stringify(pcSender.localDescription))
            })

            pcSender.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(resp.data.Sdp))))
            

            initSignalUser(userId)

        }
    }



}


function initSignalUser(userId){



    let uinfo = {
        "command": "ADD",
        "user_id": userId
    }

    pcSignal = new WebSocket("ws://localhost:8082/signal")


    pcSignal.onopen = function (event) {


        pcSignal.send(JSON.stringify(uinfo))

    }


    pcSignal.onmessage = function (event) {


        let data = event.data 

        console.log("received: ")
        console.log(data)

        let signal_data = JSON.parse(data)


        if(signal_data.command == "ADDUSER") {


            addReceiver(signal_data.user_id)

        }




    }



}

function startCall() {


  // sender part of the call
    navigator.mediaDevices.getUserMedia({video: true, audio: true})
        .then(function(stream){
            var senderVideo  = document.getElementById('senderVideo');
            senderVideo.srcObject = stream;
            var tracks = stream.getTracks();
            for (var i = 0; i < tracks.length; i++) {
                pcSender.addTrack(stream.getTracks()[i]);
            }
            pcSender.createOffer().then(function(d){pcSender.setLocalDescription(d)})
        })

  // you can use event listner so that you inform he is connected!
    pcSender.addEventListener('connectionstatechange', function (event) {
        if (pcSender.connectionState === 'connected') {
            console.log("pc sender connected")
        }
    })



}


function addReceiver(addUserId){



    pcRecievers[addUserId] = new RTCPeerConnection({
        iceServers: [
            {
            urls: 'stun:localhost:3478'
            }
        ]
    })

    

    pcRecievers[addUserId].onicecandidate = async function(event) {
        if (event.candidate === null) {

            console.log("receiver ice")

            let resp = await axios.post("/peers/room/sdp/m/" + meetingId + "/c/"+ addUserId + "/s/" + false, 
            {
                "sdp" : btoa(JSON.stringify(pcRecievers[addUserId].localDescription))
            })
            
            pcRecievers[addUserId].setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(resp.data.Sdp))))

        }
    }




    pcRecievers[addUserId].addTransceiver("video", 
        {"direction": "recvonly"}
    )



    pcRecievers[addUserId].createOffer()
        .then(function(d) {
            pcRecievers[addUserId].setLocalDescription(d)
        })





    pcRecievers[addUserId].ontrack = function (event) {

        var receiver_id = "receiverVideo-" + addUserId

        var receivers = document.getElementById('peer-receive')

        receivers.innerHTML += `
        
        <div class="layer2">
            <video autoplay id="${receiver_id}" width="160" height="120" controls muted></video>
        </div>

        `

        var receiverVideo = document.getElementById(receiver_id)
        receiverVideo.srcObject = event.streams[0]
        receiverVideo.autoplay = true
        receiverVideo.controls = true
    }


}
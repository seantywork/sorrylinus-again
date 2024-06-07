
pcSender = {}
pcReciever = {}
meetingId = "" 
userId = ""
peerId = ""

// use http://localhost:4200/call;meetingId=07927fc8-af0a-11ea-b338-064f26a5f90a;userId=alice;peerId=bob
// and http://localhost:4200/call;meetingId=07927fc8-af0a-11ea-b338-064f26a5f90a;userId=bob;peerId=alice
// start the call

function initPeers() {


    meetingId = document.getElementById('mid').value;


    userId = document.getElementById('uid').value;


    peerId = document.getElementById('pid').value;

    if (meetingId == ""){

        alert("feed meeting ID!")

        return

    }

    if (userId == ""){

        alert("feed user ID!")

        return

    }

    if (peerId == ""){

        alert("feed peer ID!")

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
        pcReciever = new RTCPeerConnection({
        iceServers: [
            {
            urls: 'stun:localhost:3478'
            }
        ]
        })




    pcSender.onicecandidate = async function(event) {
        if (event.candidate === null) {

            console.log("sender ice")

            let resp = await axios.post("/peers/sdp/m/" + meetingId + "/c/"+ userId + "/p/" + peerId + "/s/" + true,
            {
                "sdp" : btoa(JSON.stringify(pcSender.localDescription))
            })

            pcSender.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(resp.data.Sdp))))
            

        }
    }

    pcReciever.onicecandidate = async function(event) {
        if (event.candidate === null) {

            console.log("receiver ice")

            let resp = await axios.post("/peers/sdp/m/" + meetingId + "/c/"+ userId + "/p/" + peerId + "/s/" + false, 
            {
                "sdp" : btoa(JSON.stringify(pcReciever.localDescription))
            })
            
            pcReciever.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(resp.data.Sdp))))

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


    pcReciever.addTransceiver("video", 
        {"direction": "recvonly"}
    )

    pcReciever.createOffer()
        .then(function(d) {
            pcReciever.setLocalDescription(d)
        })

    pcReciever.ontrack = function (event) {
        var receiverVideo = document.getElementById('receiverVideo')
        receiverVideo.srcObject = event.streams[0]
        receiverVideo.autoplay = true
        receiverVideo.controls = true
    }

}

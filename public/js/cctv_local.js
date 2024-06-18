/*
let pc = new RTCPeerConnection({
    iceServers: [
      {
        urls: 'stun:stun.l.google.com:19302'
      }
    ]
  })

*/

TURN_SERVER_ADDRESS = ""

CLIENT_REQ = {
    "data":""
}

pc = {}

log = function(msg) {
    document.getElementById('div').innerHTML += msg + '<br>'
}




async function startCCTVOffer(){


    let offer_val = document.getElementById('localSessionDescription').value


    let offerjson = JSON.parse(JSON.stringify(CLIENT_REQ))

    offerjson.data = offer_val

    let result = await axios.post("/cctv/offer", offerjson)

    if (result.data.status != "success") {

        alert("failed to start cctv offer")
    }


    document.getElementById('remoteSessionDescription').value = result.data.reply



}


function startSession() {
    let sd = document.getElementById('remoteSessionDescription').value
    if (sd === '') {
        return alert('Session Description must not be empty')
    }

    try {
        pc.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(sd))))
    } catch (e) {
        alert(e)
    }
}


async function getTurnServerAddressAndInit(){

    let result = await axios.get("/cctv/turn/address")

    if(result.data.status != "success"){

        alert("failed to get turn server address")

        return
    }


    TURN_SERVER_ADDRESS = result.data.reply 

    console.log("turnServerAddr: " + TURN_SERVER_ADDRESS)


    pc = new RTCPeerConnection({
        iceServers: [
            {
                urls: TURN_SERVER_ADDRESS
            }
        ]
    })


    pc.ontrack = function (event) {
        let el = document.createElement(event.track.kind)
        el.srcObject = event.streams[0]
        el.autoplay = true
        el.controls = true
    
        document.getElementById('remoteVideos').appendChild(el)
    }
    
    pc.oniceconnectionstatechange = function(e) {log(pc.iceConnectionState)}
    pc.onicecandidate = function(event){
        if (event.candidate === null) {
            document.getElementById('localSessionDescription').value = btoa(JSON.stringify(pc.localDescription))
        }
    }
    
    // Offer to receive 1 audio, and 2 video tracks
    pc.addTransceiver('audio', {'direction': 'recvonly'})
    pc.addTransceiver('video', {'direction': 'recvonly'})
    pc.createOffer().then(function(d){ pc.setLocalDescription(d)}).catch(log)

}


getTurnServerAddressAndInit()


pc = {}
ws = {}

PEERS_SIGNAL_ADDRESS = ""

TURN_SERVER_ADDRESS= {}


function initPeers(){


    navigator.mediaDevices.getUserMedia({ video: true, audio: true })
        .then(function(stream){

            pc = new RTCPeerConnection({
//                        iceServers: [
//                            {
//                                urls: TURN_SERVER_ADDRESS.addr,
//                                username: TURN_SERVER_ADDRESS.id,
//                                credential: TURN_SERVER_ADDRESS.pw
//                            }
//                        ]
                    })
                
            document.getElementById('localVideo').srcObject = stream
            stream.getTracks().forEach(function(track) {pc.addTrack(track, stream)})
        
            if (location.protocol !== 'https:') {

                ws = new WebSocket("ws://" + PEERS_SIGNAL_ADDRESS)
        
            } else {
        
                ws = new WebSocket("wss://" + PEERS_SIGNAL_ADDRESS)
        
        
            }
            
            ws.onclose = function(evt) {
                window.alert("Websocket has closed")
            }
        
            ws.onmessage = function(evt) {
                let msg = JSON.parse(evt.data)
        
                if (!msg) {
                    return console.log('failed to parse msg')
                }
        
        
                switch (msg.command) {
                    case 'offer':
                    let offer = JSON.parse(atob(msg.data))
                    if (!offer) {
                        return console.log('failed to parse answer')
                    }
        
                    pc.setRemoteDescription(offer)
                    pc.createAnswer().then(function(answer) {
                        pc.setLocalDescription(answer)
                        ws.send(JSON.stringify({command: 'answer', data: btoa(JSON.stringify(answer))}))
                    })
                    return
        
                    case 'candidate':
                    let candidate = JSON.parse(atob(msg.data))
                    if (!candidate) {
                        return console.log('failed to parse candidate')
                    }
        
                    pc.addIceCandidate(candidate)
                }
            }

            ws.onerror = function(evt) {
                console.log("ERROR: " + evt.data)
            }

            pc.ontrack = function (event) {
                if (event.track.kind === 'audio') {
                    return
                }
        
                let el = document.createElement(event.track.kind)
                el.srcObject = event.streams[0]
                el.autoplay = true
                el.controls = true
                document.getElementById('remoteVideos').appendChild(el)
        
                event.track.onmute = function(event) {
                    el.play()
                }
        
                event.streams[0].onremovetrack = function({track}) {
                    if (el.parentNode) {
                        el.parentNode.removeChild(el)
                    }
                }
            }
        
        
            pc.onicecandidate = function(e){
                
                if (!e.candidate) {
        
                    console.log("not a candidate")
        
                    return
                }
        
                console.log("sending candidate info to server")
        
                ws.send(JSON.stringify({command: 'candidate', data: btoa(JSON.stringify(e.candidate))}))
            }
                
        
        
            console.log("opened peer connection ready")

        })
        .catch(function(e){

            alert(e)
        })



}

async function init(){
    let result = await axios.get("/api/peers/signal/address")

    if(result.data.status != "success"){

        alert("failed to get peers signal address")

        return
    }


    PEERS_SIGNAL_ADDRESS = result.data.reply 

    console.log("peersSignalAddr: " + PEERS_SIGNAL_ADDRESS)


    console.log("opened channel for peer signal")

    result = await axios.get("/api/turn/address")

    if(result.data.status != "success"){

        alert("failed to get turn server address")

        return
    }


    let addrs = JSON.parse(result.data.reply) 

    TURN_SERVER_ADDRESS = addrs[0]

    console.log("turnServerAddr: " + TURN_SERVER_ADDRESS.addr)
}

init()
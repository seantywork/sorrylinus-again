

pc = {}
ws = {}

PEERS_SIGNAL_ADDRESS = ""

function initSignal(){


    navigator.mediaDevices.getUserMedia({ video: true, audio: true })
        .then(function(stream) {
            pc = new RTCPeerConnection()
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

            document.getElementById('localVideo').srcObject = stream
            stream.getTracks().forEach(function(track) {pc.addTrack(track, stream)})

            if (location.protocol !== 'https:') {

                ws = new WebSocket("ws://" + PEERS_SIGNAL_ADDRESS)

            } else {

                ws = new WebSocket("wss://" + PEERS_SIGNAL_ADDRESS)


            }

            pc.onicecandidate = function(e){
                
                if (!e.candidate) {
                    return
                }

                ws.send(JSON.stringify({command: 'candidate', data: btoa(JSON.stringify(e.candidate))}))
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

    }).catch(window.alert)

}

async function getSignalAddr(){


    let result = await axios.get("/peers/signal/address")

    if(result.data.status != "success"){

        alert("failed to get peers signal address")

        return
    }


    PEERS_SIGNAL_ADDRESS = result.data.reply 

    console.log("peersSignalAddr: " + PEERS_SIGNAL_ADDRESS)
    

}


getSignalAddr()
pc = {}

TURN_SERVER_ADDRESS = {}


CLIENT_REQ = {
    "data":""
}



async function initCCTV(){


    pc = new RTCPeerConnection({
//        iceServers: [
//            {
//                urls: TURN_SERVER_ADDRESS.addr,
//                username: TURN_SERVER_ADDRESS.id,
//                credential: TURN_SERVER_ADDRESS.pw
//            }
//        ]
    })

    pc.oniceconnectionstatechange = function(e) {console.log(pc.iceConnectionState)}

    pc.onicecandidate = async function(event){

        if (event.candidate === null){


            let req = {
                data: JSON.stringify(pc.localDescription)
            }

            let resp = await axios.post("/api/cctv/create", req)

            if (resp.data.status != "success") {

                alert("failed to start cctv offer")
            }
            try {
                pc.setRemoteDescription(new RTCSessionDescription(JSON.parse(resp.data.reply)))
            } catch (e){
                alert(e)
            }

        }


    }

    pc.ontrack = function (event) {

        var el = document.createElement(event.track.kind)
        el.srcObject = event.streams[0]
        el.autoplay = true
        el.controls = true

        document.getElementById('rtmpFeed').appendChild(el)
    }

    pc.addTransceiver('video')
    pc.addTransceiver('audio')
    
    let offer = await pc.createOffer()

    pc.setLocalDescription(offer)

    console.log("init success")

}


function closeCCTV(){



}



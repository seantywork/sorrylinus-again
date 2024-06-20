pc = {}

TURN_SERVER_ADDRESS = ""


CLIENT_REQ = {
    "data":""
}



async function initCCTV(){

    let result = await axios.get("/api/cctv/turn/address")

    if(result.data.status != "success"){

        alert("failed to get turn server address")

        return
    }

    let addrs = JSON.parse(result.data.reply) 

    TURN_SERVER_ADDRESS = addrs[0].addr

    pc = new RTCPeerConnection({
        iceServers: [
            {
                urls: TURN_SERVER_ADDRESS,
                username: addrs[0].id,
                credential: addrs[0].pw
            }
        ]
    })

    pc.oniceconnectionstatechange = function(e) {console.log(pc.iceConnectionState)}

    pc.onicecandidate = async function(event){

        if (event.candidate === null){


            let req = {
                data: JSON.stringify(pc.localDescription)
            }

            let resp = await axios.post("/api/cctv/create", req)

            if (result.data.status != "success") {

                alert("failed to start cctv offer")
            }
            try {
                console.log(resp.data)
                pc.setRemoteDescription(new RTCSessionDescription(resp.data))
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



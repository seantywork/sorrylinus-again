pc = {}



async function initCCTV(){

    pc = new RTCPeerConnection()
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

    console.log(offer)

    let req = {
        data: JSON.stringify(offer)
    }

    let resp = await axios.post("/api/cctv/create", req)

    pc.setRemoteDescription(resp.data)

    console.log("init success")

}


function closeCCTV(){



}

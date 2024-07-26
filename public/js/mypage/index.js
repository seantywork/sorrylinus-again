pc = {}

ws = {}

CH_SOLI_ALIVE = 0

SOLI_DISCOVERY_SENT = 1

TURN_SERVER_ADDRESS = {}

SOLI_SIGNAL_ADDRESS = ""

STREAMING_KEY = ""

CLIENT_REQ = {
    "data":""
}

SERVER_RE = {
    "status": "",
    "reply": ""
}

USER_LIST_STRUCT = {

    users: []

}

USER_CREATE = {

    id: "",
    passphrase: "",
    duration_seconds:""
}

SOLI_USER = {

    user: "",
    passphrase: ""

}

SOLI_ENTER = {

    user: "",
    key: ""

}

CCTV_STRUCT = {

    "streaming_key":"",
    "description": ""

}


async function openSoli(){


    let u_id = document.getElementById("soli-open-user-id").value 

    if(u_id == ""){
  
        alert("no soli id")
    
        return
    
    }


    let u_pw = document.getElementById("soli-open-user-pw").value 


    if(u_pw == ""){
  
        alert("no soli pw")
    
        return
    
    }

    document.getElementById("soli-open-user-id").value = ""

    document.getElementById("soli-open-user-pw").value = ""

    let su = JSON.parse(JSON.stringify(SOLI_USER))

    su.user = u_id
    su.passphrase = u_pw

    let req = {
        data: JSON.stringify(su)
    }

    let resp = await fetch(`/api/sorrylinus/open`, {
        body: JSON.stringify(req),
        method: "POST"
    })

    let result = await resp.json()

    if (result.status != "success"){

        alert("failed to open: " + result.reply)
    
        return
    }

    let enterInfo = JSON.parse(result.reply)

    console.log(enterInfo)


    if (location.protocol !== 'https:') {

        ws = new WebSocket("ws://" + SOLI_SIGNAL_ADDRESS)

    } else {

        ws = new WebSocket("wss://" + SOLI_SIGNAL_ADDRESS)


    }


    ws.onopen = function(evt){


        ws.send(JSON.stringify({command: 'auth', data: result.reply}))

        console.log("sent enter auth")

    }

    ws.onclose = function(evt) {
        alert("Soli websocket has closed")
    }

    ws.onmessage = function(evt){

        let msg = JSON.parse(evt.data)

        if (!msg) {
            console.log('failed to parse msg')
            console.log("evt--------")
            console.log(msg)
            return 
        }

        if (msg.status != "success"){

            alert("failed: " + msg.data)

            return

        } else {

            if (CH_SOLI_ALIVE == 0){
                CH_SOLI_ALIVE = 1
            }

            if (SOLI_DISCOVERY_SENT == 0){

                document.getElementById("discovery-reader").innerText = msg.data

                SOLI_DISCOVERY_SENT = 1

            } else {

                document.getElementById("soli-action-result").innerText = msg.data

            }

        }


    }

}



async function sendSoliQuery(){


    let u_query = document.getElementById("soli-action-data").value 


    if(u_query == ""){
  
        alert("no soli query")
    
        return
    
    }

    document.getElementById("soli-action-data").value = ""

    ws.send(JSON.stringify({command: 'roundtrip', data: u_query }))

}



async function initSoli(){

    let options = {
        method: "GET"
    }
    let result = await fetch("/api/sorrylinus/signal/address", options)

    let data = await result.json()

    if(data.status != "success"){

        alert("failed to get soli signal address")

        return
    }


    SOLI_SIGNAL_ADDRESS = data.reply 

    console.log("soliSignalAddr: " + SOLI_SIGNAL_ADDRESS)

}

async function sendDiscovery(){

    if(CH_SOLI_ALIVE != 1){

        alert("soli not opened, yet")

        return
    }

    SOLI_DISCOVERY_SENT = 0

    ws.send(JSON.stringify({command: 'roundtrip', data: "discovery:"}))

}

async function initCCTV(){

    if(CH_SOLI_ALIVE != 1){

        alert("soli not opened, yet")

        return
    }

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

            let options = {
                method: "POST",
                headers: {
                  "Content-Type": "application/json" 
                },
                body: JSON.stringify(req) 
            }

            let resp = await fetch("/api/cctv/open", options)

            let data = await resp.json()

            if (data.status != "success") {

                alert("failed to start cctv offer")
            }
            try {
            
                cs = JSON.parse(data.reply)

                console.log(cs)

                let remoteDesc = JSON.parse(cs.description)
                
                pc.setRemoteDescription(new RTCSessionDescription(remoteDesc))
            
                STREAMING_KEY = cs.streaming_key

                console.log("streaming address: " + cs.location)

                await delayMs(15000)
                
                ws.send(JSON.stringify({command: 'roundtrip', data:  "cctv-stream:" + cs.location}))
            
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

        document.getElementById('cctv-reader').appendChild(el)

    }

    pc.addTransceiver('video')
    pc.addTransceiver('audio')
    
    let offer = await pc.createOffer()

    pc.setLocalDescription(offer)

    console.log("init success")

}



async function testCCTV(){

    alert("test cctv")

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

            let options = {
                method: "POST",
                headers: {
                  "Content-Type": "application/json" 
                },
                body: JSON.stringify(req) 
            }

            let resp = await fetch("/api/cctv/open", options)

            let data = await resp.json()

            if (data.status != "success") {

                alert("failed to start cctv offer")
            }
            try {
            
                cs = JSON.parse(data.reply)

                console.log(cs)

                let remoteDesc = JSON.parse(cs.description)
                
                pc.setRemoteDescription(new RTCSessionDescription(remoteDesc))
            
                STREAMING_KEY = cs.streaming_key

                alert("streaming addr:" + cs.location) 
            
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

        document.getElementById('cctv-reader').appendChild(el)

    }

    pc.addTransceiver('video')
    pc.addTransceiver('audio')
    
    let offer = await pc.createOffer()

    pc.setLocalDescription(offer)

    console.log("init success")

}



async function closeCCTV(){

    let req = {
        data: STREAMING_KEY
    }

    let options = {
        method: "POST",
        headers: {
          "Content-Type": "application/json" 
        },
        body: JSON.stringify(req) 
    }

    let resp = await fetch("/api/cctv/close", options)

    let data = await resp.json()

    if (data.status != "success") {

        alert("failed to close cctv")
    } else {

        alert("successfully closed cctv:" + data.reply)
    }



}


async function listUsers(){

    let resp = await fetch(`/api/auth/user/list`, {
        method: "GET"
    })

    let result = await resp.json()

    if(result.status != "success"){
  
        alert("failed to get sample list")
    
        return
    
    }

    let userReader = document.getElementById("user-reader")

    let userList = JSON.parse(result.reply)


    if (userList.users == null){
  
        userReader.innerHTML = `
            <pre> :(     No users, yet </pre>
        `  
    
      } else {
    
        for(let i = 0; i < userList.users.length; i ++){
    
            userReader.innerHTML += `
            <p> ${userList.users[i]} </p> 
            <input type="button" onclick="deleteUser('${userList.users[i]}')" value="delete">
            <br>
            `
         
        }
      }
}


async function createUser(){

    let u_id = document.getElementById("create-user-id").value 

    if(u_id == ""){
  
        alert("no user id")
    
        return
    
    }


    let u_pw = document.getElementById("create-user-pw").value 


    if(u_pw == ""){
  
        alert("no user pw")
    
        return
    
    }


    let u_dur = document.getElementById("create-user-duration").value 


    if(u_dur == ""){
  
        alert("no user duration")
    
        return
    
    }

    let uc = JSON.parse(JSON.stringify(USER_CREATE))

    uc.id = u_id
    uc.passphrase = u_pw
    uc.duration_seconds = parseInt(u_dur, 10)

    let req = {
        data: JSON.stringify(uc)
    }

    let resp = await fetch(`/api/auth/user/add`, {
        body: JSON.stringify(req),
        method: "POST"
    })

    let result = await resp.json()

    if(result.status != "success"){

        alert("failed to create user")

        return
    }

    alert("successfully created user: " + result.reply)

    await listUsers()

}

async function deleteUser(userId){

    let req = {
        data: userId
    }


    let resp = await fetch(`/api/auth/user/remove`, {
        body: JSON.stringify(req),
        method: "POST"
    })

    let result = await resp.json()

    if(result.status != "success"){

        alert("failed to delete user")

        return
    }

    alert("successfully deleted user: " + result.reply)


    await listUsers()

}


async function flushLog(){

    let options = {
        method: "GET"
    }
    let result = await fetch("/api/manage/log/flush", options)

    let data = await result.json()

    if(data.status != "success"){

        alert("failed to flush log")

        return
    }

    alert("successfully flushed log: "+ data.reply)

}



(async function (){

    await listUsers()

    await initSoli()

})()



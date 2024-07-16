pc = {}

TURN_SERVER_ADDRESS = {}


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

CCTV_STRUCT = {

    "streaming_key":"",
    "description": ""

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

    alert("Not implemented, yet")

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
            <pre> ${userList.users[i]} </pre> 
            <button onclick="deleteUser('${userList.users[i]}')">Delete</button>
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




(async function (){

    await listUsers()

})()



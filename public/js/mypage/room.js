
PEERS_CREATE = {

    room_name: "",
    users: []

}

async function getRoomList(){

    let resp = await fetch("/api/peers/entry", {
      method: "GET"
    })
  
    let result = await resp.json()
  
    if(result.status != "success"){
  
      alert("failed to get sample list")
  
      return
  
    }
  
  
    let roomReader = document.getElementById("room-reader")
  
    let roomEntry = JSON.parse(result.reply)
  
    if (roomEntry.room_name == null){
  
      roomReader.innerHTML = `
          <pre> :(   No room exists, so far </pre>
      `
  
  
    } else {
    
        roomReader.innerHTML = ""
  
        for(let i = 0; i < roomEntry.room_name.length; i ++){
    
            roomReader.innerHTML += `
            <a class="tui-button" href="/room/${roomEntry.room_name[i]}">
                ${roomEntry.room_name[i]}
            </a>
            <input class="tui-button" type="button" onclick="deleteRoom('${roomEntry.room_name[i]}')" value="delete">
            <br>
            `
        
        }
    }
  
  
  }


async function createRoom(){



    let roomName = document.getElementById("create-room-name").value 

    if(roomName == ""){
  
        alert("no room name")
    
        return
    
    }


    let roomUsers = document.getElementById("create-room-users").value 


    if(roomUsers == ""){
  
        alert("no room users")
    
        return
    
    }


    roomUsers = roomUsers.replace(" ", "")

    let roomUserList = roomUsers.split(",")

    console.log(roomUserList)

    let p_create = JSON.parse(JSON.stringify(PEERS_CREATE))

    p_create.room_name = roomName

    for(let i = 0 ; i < roomUserList.length; i ++){

        if(roomUserList[i] != ""){
            p_create.users.push(roomUserList[i])
        }
    }

    let req = {
        data: JSON.stringify(p_create)
    }

    let resp = await fetch(`/api/peers/create`, {
        body: JSON.stringify(req),
        method: "POST"
    })

    let result = await resp.json()

    if(result.status != "success"){

        alert("failed to create room")

        return
    }

    alert("successfully created room: " + result.reply)

    await getRoomList()

}


async function deleteRoom(roomName){

    let req = {
        data: roomName
    }


    let resp = await fetch(`/api/peers/delete`, {
        body: JSON.stringify(req),
        method: "POST"
    })

    let result = await resp.json()

    if(result.status != "success"){

    alert("failed to delete room")

    return
    }

    alert("successfully deleted room: " + result.reply)


    await getRoomList()


}



(async function() {

    await getRoomList()
 
 })()
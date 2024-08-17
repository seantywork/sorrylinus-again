var editor_content = {}

var ENTRY_STRUCT = {

    entry:[
        {
            "title":"",
            "id":"",
            "type":""
        }
    ]

}

var PEERS_ENTRY_STRUCT = {

    room_name:[]

}



var CONTENT_LIST = []

var ROOM_LIST = []

var CONTENT_PAGE_PTR = 0
var ROOM_PAGE_PTR = 0

var PAGE_MAX = 5


async function getContentList(){

  let resp = await fetch("/api/content/entry", {
    method: "GET"
  })

  let result = await resp.json()

  if(result.status != "success"){

    alert("failed to get sample list")

    return

  }

  CONTENT_LIST = JSON.parse(result.reply)

  renderContentList()

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
  
 
    ROOM_LIST = JSON.parse(result.reply)


    renderRoomList()
  
  
  }




function renderContentList(){


  let contentReader = document.getElementById("content-reader")

  let contentEntry = JSON.parse(JSON.stringify(CONTENT_LIST))

  contentReader.innerHTML = ""

  if (contentEntry.entry == null){

    contentReader.innerHTML = `

        <pre> :(    Nothing to see here, yet </pre>


    `

  } else {

    let sortedEntry = getNewDateSortedList("desc", "timestamp", contentEntry.entry)

    let pageStart = CONTENT_PAGE_PTR * PAGE_MAX
    let pageEnd = pageStart + PAGE_MAX

    for(let i = 0; i < sortedEntry.length; i ++){

      if(pageStart <= i && i < pageEnd){
       
        contentReader.innerHTML += `
        <div style="display: block;">
          <a class="tui-button" href="/content/${sortedEntry[i].type}/${sortedEntry[i].id}">
            ${sortedEntry[i].title} 
          </a> [${sortedEntry[i].author}:${sortedEntry[i].timestamp}]
        </div> 
        <br>
        `
      } else {

        contentReader.innerHTML += `
        <div style="display: none;">
          <a class="tui-button" href="/content/${sortedEntry[i].type}/${sortedEntry[i].id}">
            ${sortedEntry[i].title} 
          </a> [${sortedEntry[i].author}:${sortedEntry[i].timestamp}] 
        </div>
        <br>
        ` 
      }
     
    }
  }
}


function renderRoomList(){



  let roomReader = document.getElementById("room-reader")
  
  let roomEntry = JSON.parse(JSON.stringify(ROOM_LIST))

  roomReader.innerHTML = ""

  if (roomEntry.room_name == null){

    roomReader.innerHTML = `

        <pre> :(     You're not invited, yet </pre>

    `


  } else {

    let pageStart = ROOM_PAGE_PTR * PAGE_MAX
    let pageEnd = pageStart + PAGE_MAX

    for(let i = 0; i < roomEntry.room_name.length; i ++){
   
      if(pageStart <= i && i < pageEnd){
        roomReader.innerHTML += `
        <div style="display: block;">
          <a class="tui-button" href="/room/${roomEntry.room_name[i]}">
              ${roomEntry.room_name[i]}
          </a>
        </div>
        <br>
        `

      } else {

        roomReader.innerHTML += `
        <div style="display: none;">
          <a class="tui-button" href="/room/${roomEntry.room_name[i]}">
              ${roomEntry.room_name[i]}
          </a>
        </div>
        <br>
        `
      }

  
    }

  }

}

function contentNext(){


  let contentLength = CONTENT_LIST.entry.length 

  let tmpPagePtr = (CONTENT_PAGE_PTR + 1) * PAGE_MAX

  if(tmpPagePtr >= contentLength){

    alert("goto: content: first page")

    CONTENT_PAGE_PTR = 0

  } else {

    CONTENT_PAGE_PTR += 1
  }

  renderContentList()


}

function roomNext(){

  let roomLength = ROOM_LIST.room_name.length

  let tmpPagePtr = (ROOM_PAGE_PTR + 1) * PAGE_MAX

  if(tmpPagePtr >= roomLength){

    alert("goto: room: first page")

    ROOM_PAGE_PTR = 0

  } else {

    ROOM_PAGE_PTR += 1

  }

  renderRoomList()

}


(async function() {

    await getContentList()

    await getRoomList()
 
 })()
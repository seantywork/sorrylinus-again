
async function getVideoList(){

    let resp = await fetch("/api/content/entry", {
      method: "GET"
    })
  
    let result = await resp.json()
  
    if(result.status != "success"){
  
      alert("failed to get sample list")
  
      return
  
    }
  
  
    let contentReader = document.getElementById("video-reader")
  
    let contentEntry = JSON.parse(result.reply)
  
    let videoCount = 0

    if (contentEntry.entry == null){
    
        contentReader.innerHTML = `
            <pre> :(    Nothing to see here, yet </pre>
        `
  
        return
    }

    contentReader.innerHTML = ""
  
    for(let i = 0; i < contentEntry.entry.length; i ++){

        if(contentEntry.entry[i].type != "video"){
            continue
        }

        contentReader.innerHTML += `
        <a href="/content/${contentEntry.entry[i].type}/${contentEntry.entry[i].id}">
            ${contentEntry.entry[i].title}
        </a>
        <button onclick="deleteVideo('${contentEntry.entry[i].id}')">Delete</button>
        <br>
        `
        videoCount += 1
    }

    if (videoCount == 0){
        contentReader.innerHTML = `
        <pre> :(    Nothing to see here, yet </pre>
    `
 
    }
  
  
  }

async function deleteVideo(videoId){

  let req = {
    data: videoId
  }


  let resp = await fetch(`/api/video/delete`, {
      body: JSON.stringify(req),
      method: "POST"
  })

  let result = await resp.json()

  if(result.status != "success"){

    alert("failed to delete video")

    return
  }

  alert("successfully deleted video: " + result.reply)


  await getVideoList()



}
  
(async function() {

    await getVideoList()

})()
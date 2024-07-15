function getVideoContent(){

    let contentId = CONTENT_ID

    if (contentId == "") {

        alert("no content id provided")

        return
    }

    videoEl = document.getElementById('video-reader')

  
    let addHtml = `
        <source src="/api/video/c/${contentId}" type="video/mp4">
    `
  
    videoEl.innerHTML = addHtml
  
  
  }
  


(async function(){

    await getVideoContent()

})()
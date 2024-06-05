


function playFiles(){

    videoEl = document.getElementById('play-files')

    watchID = document.getElementById('wid').value;


    if (watchID == ""){

        alert("feed watch ID!")

        return

    }

    let addHtml = `
        <source src="/files/d/${watchID}" type="video/mp4">
    `

    videoEl.innerHTML = addHtml


}
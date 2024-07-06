


function handleSubmit(event) {
    const form = event.currentTarget;
    const url = new URL(form.action);
    const formData = new FormData(form);
    const searchParams = new URLSearchParams(formData);
  
    const fetchOptions = {
      method: form.method,
    };
  
    if (form.method.toLowerCase() === 'post') {
      if (form.enctype === 'multipart/form-data') {
        fetchOptions.body = formData;
      } else {
        fetchOptions.body = searchParams;
      }
    } else {
      url.search = searchParams;
    }
  
    fetch(url, fetchOptions);
  
    event.preventDefault();
}

function playFiles(){

  videoEl = document.getElementById('play-files')

  watchID = document.getElementById('wid').value;


  if (watchID == ""){

      alert("feed watch ID!")

      return

  }

  let addHtml = `
      <source src="/video/watch/c/${watchID}" type="video/mp4">
  `

  videoEl.innerHTML = addHtml


}

const form = document.querySelector('form');
form.addEventListener('submit', handleSubmit);

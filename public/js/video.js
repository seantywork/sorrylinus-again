
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

const form = document.querySelector('form');
form.addEventListener('submit', handleSubmit);

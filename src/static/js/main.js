function ready(cb) {
  // in case the document is already rendered
  if (document.readyState != "loading") {
    cb()
  } else {
    document.addEventListener("DOMContentLoaded", cb)
  }
}

// TODO: https://developer.mozilla.org/en-US/docs/Web/API/History_API
// set the filter as a query string and auto-filter on page load.

ready(function() {
  let tags = document.querySelectorAll(".hashtag")
  tags.forEach(tag => {
    tag.addEventListener("click", e => {
      console.log(e.target.innerText)
      e.preventDefault()
    })
  })
})

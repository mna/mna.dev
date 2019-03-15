let allCards = []

function ready(cb) {
  // in case the document is already rendered
  if (document.readyState != "loading") {
    cb()
  } else {
    document.addEventListener("DOMContentLoaded", cb)
  }
}

function debounced(delay, fn) {
  let timerId
  return function (...args) {
    if (timerId) {
      clearTimeout(timerId)
    }
    timerId = setTimeout(() => {
      fn(...args)
      timerId = null
    }, delay)
  }
}

// TODO: https://developer.mozilla.org/en-US/docs/Web/API/History_API
// set the filter as a query string and auto-filter on page load.

function filterCards(e) {
  let text = e.target.value.trim()
  let words = text.match(/\S+/g) || []

  allCards.forEach(card => {
    let content = card.innerText

    if (words.every(w => content.includes(w))) {
      card.classList.remove("hidden")
    } else {
      card.classList.add("hidden")
    }
  })
}

ready(function() {
  let tags = document.querySelectorAll(".hashtag")
  tags.forEach(tag => {
    tag.addEventListener("click", e => {
      console.log(e.target.innerText)
      e.preventDefault()
    })
  })

  allCards = document.querySelectorAll(".card")

  let search = document.getElementById("search")
  search.addEventListener("input", debounced(200, filterCards))
})

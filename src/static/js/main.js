;(function() {

let allCards = []
let haystackCards = []
let allTags = []
let selectedTags = new Set()
let searchBox = null

// watch for document ready and call cb when it is.
function ready(cb) {
  // in case the document is already rendered
  if (document.readyState != "loading") {
    cb()
  } else {
    document.addEventListener("DOMContentLoaded", cb)
  }
}

// returns a function that calls fn only after delay milliseconds
// have passed without other calls.
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

// resets the search box.
function clearSearch() {
  searchBox.value = ""
  // TODO: clear query string or whatever is used to link to a search
}

// show only cards from the haystack that match the searched words.
function filterCardsBySearch(e) {
  e.preventDefault()

  let text = e.target.value.trim().toLowerCase()
  let words = text.match(/\S+/g) || []

  // search only in the haystack, i.e. in the currently selected tag(s)
  haystackCards.forEach(card => {
    let content = card.innerText
    let tweet = card.querySelector(".twitter-tweet")
    if (tweet) {
      let tweetText = tweet.shadowRoot && tweet.shadowRoot.querySelector(".Tweet-text")
      if (!tweetText) {
        tweetText = tweet.querySelector("p")
      }
      if (tweetText) {
        content = tweetText.innerText
      }
    }

    content = content.toLowerCase()
    if (words.every(w => content.includes(w))) {
      card.classList.remove("hidden")
    } else {
      card.classList.add("hidden")
    }
  })
}

function isSuperset(set, subset) {
  for (var elem of subset) {
    if (!set.has(elem)) {
      return false
    }
  }
  return true
}

function tagClicked(e) {
  e.preventDefault()

  let selectedTag = e.target.innerText
  if (selectedTags.has(selectedTag)) {
    selectedTags.delete(selectedTag)
  } else {
    selectedTags.add(selectedTag)
  }

  let keys = Array.from(selectedTags.keys())
  if (keys.length === 0 ) {
    history.pushState(null, "", "/")
  } else {
    history.pushState(null, "", "?tags=" + encodeURIComponent(keys.join(",")))
  }

  applySelectedTags()
}

function applySelectedTags() {
  // filtering on a tag clears the search and resets the haystack
  clearSearch()
  haystackCards.length = 0

  allTags.forEach(tag => {
    let tagText = tag.innerText
    if (selectedTags.has(tagText)) {
      tag.classList.add("is-active")
    } else {
      tag.classList.remove("is-active")
    }
  })

  allCards.forEach(card => {
    let cardTags = new Set(Array.from(card.querySelectorAll(".hashtag")).map(ht => ht.innerText))

    let tweet = card.querySelector(".twitter-tweet")
    if (tweet) {
      // tweets have "#twitter" and "#micro" tags implied
      cardTags = new Set(["#twitter", "#micro"])
    }

    if (isSuperset(cardTags, selectedTags)) {
      card.classList.remove("hidden")
      haystackCards.push(card)
    } else {
      card.classList.add("hidden")
    }
  })
}

function tagsFromQueryString() {
  let params = new URL(window.location).searchParams
  let tags = params.get("tags")
  console.log(tags)

  if (tags) {
    selectedTags = new Set(tags.split(","))
  } else {
    selectedTags.clear()
  }
  applySelectedTags()
}

// bootstrap execution, called when document is ready.
ready(function() {
  // grab list of all cards, and on load the haystack is the set of all cards
  allCards = Array.from(document.querySelectorAll(".card"))
  haystackCards = allCards.slice()

  allTags = Array.from(document.querySelectorAll(".hashtag"))
  allTags.forEach(tag => {
    tag.addEventListener("click", tagClicked)
  })

  searchBox = document.getElementById("search")
  searchBox.addEventListener("input", debounced(200, filterCardsBySearch))

  window.onpopstate = tagsFromQueryString

  tagsFromQueryString()
})

})();

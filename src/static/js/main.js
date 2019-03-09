document.addEventListener('DOMContentLoaded', initMenu)

function initMenu() {
  const menus = Array.prototype.slice.call(document.querySelectorAll('.navbar-burger'), 0)

  menus.forEach(el => el.addEventListener('click', menuClick))
}

function menuClick() {
  // Get the target from the "data-target" attribute
  const targetId = this.dataset.target
  const target = document.getElementById(targetId)

  // Toggle the "is-active" class on both the "navbar-burger" and the "navbar-menu"
  this.classList.toggle('is-active')
  target.classList.toggle('is-active')
}

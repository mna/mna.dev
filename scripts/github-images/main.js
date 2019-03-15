let OctoCatCrawler = require("octocat-images")
let crawler = new OctoCatCrawler()
let args = process.argv.slice(2)

if (args.length != 1) {
  throw `want 1 argument, the output image directory, got ${args.length}`
}

// List all octocat images
let imgDir = args[0]
crawler.list((err, octocats) => {
  if (err) {
    throw err
  }
  for (let octocat of octocats) {
    console.log(octocat.number, octocat.name, octocat.url)
    octocat.save(imgDir)
  }
})

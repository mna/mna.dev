# Launching for less than 6$/month

I recently launched my latest side-project, [fiends.io](https://www.fiends.io) - a search engine for heavy metal reviews. Yeah, it's *really* niche, basically a niche within a niche (metalheads who care about album reviews!?). I made it for me first, hopefully a few others will find value in it too. But that's not the point of this post, I want to talk about how *little* up-front money was required to launch it (and how little it takes to keep it running month after month) to help dismiss the idea or alleviate the fear that it takes significant funds to bring such a project to life.

I'm not going to discuss marketing nor development budgets, but everything that covers the infrastructure and operations to run - as in my case - a web-based system.

## Minimal but thorough

I implemented fiends.io based on [my "Small is Beautiful" philosophy](https://mna.dev/posts/small-is-beautiful.html) that I outlined in a blog post last year. Being a search engine, a good chunk of the complexity is hidden from view - the website is simple and straightforward, but behind the scenes, there's a lot going on. I have a bot to visit the pages, handle network errors, retries, throttling requests to respect the robots.txt crawl delays, etc., a message queue to asynchronously process the html, extract the artist, album, rating, index the reviews, optimize the images, etc. and the web server. Everything requires only three core components in the infrastructure: the Lua VM (I implemented everything in [Lua](https://www.lua.org/), a little language that is *fast enough* for many tasks, easy to understand and really *fun* to work with), the almighty [Postgresql database](https://www.postgresql.org/) and a Linux OS.

As of this writing, it visited close to 100 000 distinct URLs and indexed over 40 000 album reviews. It is robust, reliable and quite efficient. You could say it is a small-but-complete, minimal-but-thorough modern system.

On the operations side, I collect a lot of metrics (new URLs, worker messages processed, web requests, requests duration, database connections, etc.) and the logs of all components of the system. I have dashboards with alerts triggered when various metrics cross a threshold. I have daily, weekly and monthly database backups stored outside the server.

I care a lot about users' privacy, so although I do have web analytics, I didn't go for Google Analytics. Oh, and I have email addresses with my custom domain.

All this for a few pennies above 5$/month.

## Breakdown of costs

* **The server: 5$/month**. At the top of the list, the better part of the budget is eaten by the Virtual Private Server (VPS). My custom deployment script handles multiple providers, it could be on [Linode](https://www.linode.com/) or [Digital Ocean](https://www.digitalocean.com/) or a number of other providers, their pricing is similar and the smallest server is sufficient for a surprisingly reasonable number of concurrent users, even with all the backend processes doing their work.

* **Web Analytics: 0$/month**. Even though I passed on Google Analytics, it turns out [Cloudflare started offering a privacy-oriented analytics platform](https://www.cloudflare.com/web-analytics/), for free. It is a bit light on features, but sufficient for my needs.

* **Metrics and Logging: 0$/month**. Many observability platforms offer free plans, but often those are limited either in the retention window (e.g. single or few days for logs) or in the features (e.g. no alerting available for free). [New Relic](https://newrelic.com/pricing) has a great free plan that includes 100 GB per month of data ingestion and one full access user, perfect for solo founders or very small teams. I can set up alerts, query my logs from a month ago, create dashboards, add my custom metrics (StatsD-based), etc.

* **Storage for Backups: <1$/month**. The database is already relatively big, given the nature of the data. It's a multi-GB backup. I have a systemd timer on the server that generates the backup every day and uploads it to the external storage. I went with [Backblaze B2 Cloud Storage](https://www.backblaze.com/b2/cloud-storage-pricing.html), which comes with 10 GB free, is AWS S3-compatible and is among the cheapest provider for common operations. The S3 compatibility is nice as you can benefit from the multitude of existing S3 tools - I use [the `s3cmd` command-line tool](https://github.com/s3tools/s3cmd) to upload the backups or download one for a restore. My backup strategy is very simple and leverages the built-in B2 lifecycle feature - I upload all my daily backups under the same name, and using the lifecycle settings keep only the last 7 versions. Similarly, on the first day of each week the same backup is also stored as the weekly backup (again, the lifecycle for the weekly one takes care of keeping the last few versions) and same thing on the first day of each month. My current usage is at around 25 cents per month.

* **Custom email addresses: 0$/month**. I use the free service [Forward Email](https://forwardemail.net) to get addresses with my custom domain name. They offer paid plans for more advanced features, but I didn't need them, the free plan is fine for my use-case.

I discovered that last service for email addresses on that directory of Things-as-a-Service that offer free plans for developers: https://free-for.dev, that's a great place to look for when launching on a budget.

This means that something like 200$ can be more than enough to run such a service for a year, giving you room to upgrade to bigger servers at some point in your growth (e.g. 10$/month after 3 months, 15$/month after 6, you're still well within that budget).

One last thing while we're touching on scaling: I don't think you should be optimizing for nines at such an early stage, both in terms of system complexity and financial costs. What you should absolutely do, though, is to have plans (multiple!) for when/if this is needed. I probably have another article's worth of things to say about this - if there's some interest, I'll write a follow-up on that subject.

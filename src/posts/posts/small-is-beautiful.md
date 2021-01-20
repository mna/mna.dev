# Small is Beautiful (The Developer's Edition)

{{template "post-meta.html" .Config}}

Modern systems - often *web apps*, with multiple other components outside the web server itself, so *web systems* - are incredibly complex. Microservices, serverless, distributed, container orchestration, eventually consistent, no single point of failure, etc., the concepts, technologies and parties involved pile up. They quickly become an inextricable maze (I won't say mess, as some are very well designed, but a complex maze nonetheless) that no single person on the engineering team has a _complete and thorough_ understanding of. There are so many configurations, services, third-party tools, proxies, moving parts that when something breaks, it can take people from *multiple teams* with different access privileges just to investigate what went wrong where, nevermind to fix it.

This post is about an exercise in restraint. In minimalism and simplicity. This is about building and designing _human-scale_ web systems.

One where _everyone_ on the engineering team can quickly gain complete and thorough understanding. One where knowledge of only a handful of technologies is required to navigate and learn the system by yourself, should you have to. Now, [simple is not easy][sine], and one must be careful of the [stupid simple][ss] trap. The goal is to have _useful_ simple, _reliable_ simple, _fast enough_, not _fastest possible_. Oh, and _fun_!

Who would that be for? Well for starters, many systems never need to operate at a huge scale, be it because it is a company-internal system, a geographically-constrained system (e.g. government agency or nonprofit organization for its citizens), or some kind of niche product, having a limited potential audience. Even for those "shoot-for-the-stars" startups, most will never get big enough to worry about scaling problems, and many will have to pivot instead of scale - something a small, simple system will help to do. And finally, it's good to remember that the [Facebook][fb], [Twitter][tw], [Shopify][shop] of this world all started with PHP or Ruby on Rails monoliths and - interestingly, and only much later on - took quite different approaches in scaling their codebase - approaches decided based on actual needs *at this point in time* that they couldn't have guessed at the start.

## The Small Architecture

What would that small approach to web systems look like, and what would it offer? There are many lenses through which such a system can be defined, analyzed and evaluated. I'll focus on three: infrastructure, architecture and operations. Note that I will only discuss back-end concerns, not front-end.

### Infrastructure

That small architecture relies on only three core components - a programming language, a storage engine (database), and of course an operating system. That's it, and that's enough to go a long way, as we'll see. Although what matters most in this post are the _concepts_, not the specific _technologies_, I will still discuss a concrete implementation of those concepts, using:

* the Lua programming language
* the PostgreSQL database
* the Linux operating system (the Fedora distribution, in my personal case)

You could easily imagine a system where you applied similar concepts using different technologies, however I strongly suggest you pick:

* a small, simple and _fast enough_ programming language (I think Go, Python or Javascript would be fine alternatives)
* a robust, proven _relational_ database (mysql/mariadb would be a fine alternative)
* your preferred OS

Of course those choices end up being a bit subjective, although things like the team's background, the local (or global, even) market, available support in the targeted deployment data center, special requirements for the product itself that make one language or database significantly better suited than another, etc. all come into play.

What should already be obvious is the small amount of knowledge required to understand the system: being familiar enough with the language and database, and a good working knowledge of the operating system (especially some command-line skills and an understanding of how long-running services/daemons work, e.g. systemd in my specific case). All components live on a single server. That's the whole infrastructure.

### A Short Detour: Lua

I can hardly imagine any useful programming language that embodies simplicity, minimalism and efficiency as well as Lua does. And to me, it also plays a huge part in the _fun_ factor. Here is a language that requires very little formalities, where you can translate your ideas to code cleanly, that has a good module system, a nice REPL to test things out, an amazing C API when needed, and a surprisingly vast ecosystem of libraries and tools (the fact that it is almost trivial to expose any C library to Lua certainly helps).

I should write down a proper love letter at some point, but for this post, I'll just mention those few points:

* It has 2 main gotchas (every language has some) that you will read about in almost every Lua post out there: variables are declared as _global_ by default, and array indices start at 1, not 0. The first is trivial to turn into a non-issue, just use [luacheck][] and configure your editor to run it on save, it will warn you if you didn't use `local` to declare a variable. For the array index, most of the times you will use _iterators_ in `for` loops so it is not a problem, otherwise you'll just have to get used to it. After a while, it's switching to other languages that get tricky.

* It supports multiple return values, and its error-handling idiom is to `return nil, 'error'` (i.e. `nil` as first returned value, and the error message next), so that calling such a function in an `assert(call_fn())` call raises that error (because nil is falsy and fails the assertion, like `false`, and the error message becomes the assertion failure message).

* It has a single `table` compound type, it can be used as array and dictionary/object-like. Its literal syntax make it popular to use as "declarative", configuration-like format, which is useful because it is still just Lua code, so you can build a table declaratively using the literal syntax, or imperatively with code that generates the table.

* It has a powerful `metatable` support, enabling among other things prototype inheritance (a bit like Javascript).

* It has built-in support for coroutines, which are first-class values (_collaborative multithreading_ via `resume` and `yield`).

* As mentioned, it has great support for _iterators_ so you can use `for` loops over anything you want, any way you want.

A fancy hello world could look like that (which you can run with `$ lua <file>`):

```lua
local msg = 'Hello, world!'
local n, inc = 0, true
for c in string.gmatch(msg, '.') do
  if inc then
    n = n + 1
  else
    n = n - 1
  end
  if c == ' ' then
    inc = false
  end
  print(string.rep(' ', 2*n) .. c)
end
```

### Architecture

A simple web system still needs to be secure and requires a number of architectural components to reliably and efficiently do its work.

* HTTPS, with the advent of [Let's Encrypt][le] and other free TLS certificate providers, this is a given.
* SQL injection, cross-site scripting (XSS), cross-site request forgery (CSRF) and other common attacks protection.
* Secure authentication, authorization and session handling.
* Support for reliable asynchronous processing (e.g. a message queue), so that e.g. a request handler can just successfully enqueue a message and expect it to be processed at some point, as some processing-heavy tasks are better done outside of a handler.
* Support for scheduled processing of jobs (e.g. cron-like tasks), as most systems have some kind of recurring jobs that need to run.
* Support for notifications (e.g. publish-subscribe) so that one part of the system can trigger a notification and other part(s) can subscribe to those.
* Database connection pooling with configurable limits to prevent overloading the system.
* Timeouts and limits basically everywhere for I/O operations.
* Logging and metrics capturing, although this will be addressed more in the Operations section.

It turns out that all of this can be achieved in a reasonable and efficient way with just the 3 core components of our infrastructure. I started a Lua *web systems framework* that encapsulates those ideas in [tulip][]. There are already a number of Lua web frameworks, but they all rely on external hosts (e.g. OpenResty, which is a modified nginx for Lua, or Apache with `mod_lua`) and don't support the latest plain Lua version (often targeting LuaJIT, which is an optimized Lua version snapshotted around 5.2 instead of the latest 5.4). Tulip is based on a plain Lua stack ([cqueues][] and [lua-http][lhttp]), PostgreSQL and a POSIX operating system - in other words, it is a concrete implementation of the infrastructure and architecture I describe here.

Many of the points in the list above depend on the programming language and its available libraries, but regarding specifically the message queue, cron jobs and notifications, with postgresql:

* **message queue** can be done a number of different ways, but in Tulip, there are `pending`, `active` and `dead` tables; the messages get enqueued in the pending one under a given queue name, when dequeued by a worker for processing they are atomically moved in the active table, and are either marked as "done" by the worker (which deletes the message) or moved back to the pending table if they expire their "time-to-live" (which, interestingly, leverages the scheduled jobs feature to run a "garbage collector" of such messages). If a message has been processed a given number of times without success, it gets moved to the dead table instead of back to the pending one, so it can be investigated manually (and this should be paired with an alert).
* **scheduled (cron-like) jobs** are implemented using the postgresql extension [pg\_cron][cron]. It supports the same scheduling options as the `cron` unix command-line, calling the provided SQL command when the right time comes. To run actual programs (i.e. not just a SQL command), Tulip leverages the message queue feature - scheduling a job triggers sending a message to a queue, and the worker then gets the job to process. Of course there can be some delays depending on the worker's load, but under normal circumstances the worker gets the message in a timely manner (workers lag should be part of the monitoring).
* **pub-sub-like notifications** are natively supported by postgresql with the [notify/listen][pubsub] feature.

For the web server part, and the code in general, even though the small architecture results in a monolith, it doesn't mean it should all be tied up together. Tulip itself is based on very loosely-coupled packages, using a [well-defined, simple and small "package contract"][contract]. It is easily extendable and in fact I recommend building the various parts of the application as distinct, small and isolated packages too.

Some common workflows such as password reset (i.e. the "forgot password?" feature) and email verification - using a signed random token sent via email - can all be done with this small architecture too and are in fact provided as part of the `account` package in tulip. It is based on a generic `token` package that handles creation of such random tokens, with a maximum age and optional "expire on first use" behaviour. This token package also handles the session IDs. Once again leveraging the scheduled jobs, there is a command that runs at regular intervals to remove expired tokens from the table. An authorization middleware is available to restrict some URLs only to e.g. authenticated or verified users, or to members of a specific group.

### Operations

I'll also address in this section the **required components for development**, as this is also a huge potential for complexity. In fact let's start with that. What is absolutely needed is:

* a version control system
* a host for the repository and a way to track issues, review pull requests/patches
* some form of [CI/CD][cicd] (continuous integration, delivery and deployment)
* a way to run a local environment

I won't get into the company chat, documentation, etc. as those are not development-specific concerns. The text editor or the "general programming environment" should be left to the preference of the developer - I don't know if it has to be said, but do not force a text editor/IDE choice to developers. They know better than you (whoever _you_ are in this scenario) what they like and with which tools they are most comfortable and productive.

I think the first point is obvious, `git` is the clear choice here and it's fair to assume most developers are at least familiar with it, as such it doesn't add a huge complexity tax (even if the tool is not simple). Source code hosting, issue and pull requests/patches tracking has many popular options - Github, Gitlab, Bitbucket, etc. - but my personal choice would be one that resonates well with the minimal and simple mindset, [Sourcehut][sr]. Incidentally, that platform also has great support for builds to use as CI/CD, so the same choice ticks two boxes (most source forges have some support for this too, of course). It should run tests and any other static checks or linting tools required. The nice thing with sourcehut's builds is that the build definition can be part of the repository, so by checking it out locally, everyone has access to the complete lifecycle of the system, including the CI/CD definition. We'll talk about deployment a bit more in a minute.

It must be possible to run such a system locally, this has many benefits including a much shorter feedback loop and easier debugging. Too often this requires multiple convoluted steps, pages of documentation, many tools to install, etc. and ends up lagging behind anyway, in a semi-working state or with outdated documentation.

Given the tiny infrastructure involved here, it is trivial to run the system locally, assuming you run the same (or compatible enough) operating system: install postgresql, create a database and configure the application to connect to it. However, I personally prefer to have everything related to a project live inside that project's directory - so in this case, to have all postgresql data and configuration in that directory, just like I don't install my Lua dependencies system-wide, but locally in a `lua_modules` subdirectory (using a small script I wrote, [llrocks][]). As I'm familiar with Docker, I don't mind using it to run the postgresql instance in this situation (and I actually use `docker-compose` even if it's just a single database service, because it provides a more declarative way to encode the configuration and I prefer its command-line UI vs docker itself). But I know many talented and experienced developers who dislike or don't know docker too well and wouldn't want to be forced to have that dependency - unlike `git`, I'd argue `docker` does add a more significant complexity tax even if you'd practically only need to run `docker-compose up`, `docker-compose start` and `docker-compose stop`, so it's nice to have both options.

Now regarding **operations-related** things proper, we need at least:

* CI/CD as already discussed, but more specifically a deployment mechanism
* Data backup and restore
* Observability, monitoring and alerting
* Good support to analyze and investigate incidents
* Scaling scenarios

Some modern tools offer "infrastructure as code" where you declaratively configure the infrastructure you want, and the tool turns it into reality. It is really great, but it is also quite complex, as the way to get the infrastructure up is a bit of a "magical black box" and it can be [tricky to debug][terra] when it fails. For such a small and simple infrastructure, I don't think it's worth adding this kind of complexity. Instead, a straightforward command-line script - ideally written in the same programming language as the rest of the system - that deploys the application using basic, imperative commands in an easily readable sequence of steps not only makes it easy to deploy manually or automatically (e.g. as part of CI/CD), but makes it easy to understand the requirements to run the system, where it logs stuff, where configuration is stored, etc.

I have started work on [such a script][deploy] for a tulip-based system. It is not meant to be general - it is a starting point that should be adapted for each system's needs and should be stored in the system's repository, but it gives a good idea of how clear and simple it can be. It is nice to use, and although in its current form it uses [Digital Ocean][do] as Virtual Private Server (VPS) provider, it could easily be changed to another, as most expose their features through an API anyway. It supports those distinct steps:

```
The deploy script is the combination of:
1. create an infrastructure (optional)
2. install the database on that infrastructure (optional)
3. deploy the new code to that infrastructure (optional)
4. (re)start the application's services
5. activate this deployment (optional)
```

And its usage looks like this:

```
$ deploy www.example.com

By default, this is sufficient to run the most common case, which
is to deploy the latest tagged version to the existing infrastructure
associated with that sub-domain.

$ deploy --create 'region:size[:image]' --ssh-keys 'list,of,names' --tags 'list,of,tags' www.example.com
$ deploy --create ... --with-db 'db-backup-id' www.example.com

$ deploy --without-code www.example.com

Only restarts the services on the existing infrastructure, no code nor database is
deployed (so, execute step 4 only).
```

Once the base OS image is created (which may take about 10 minutes or so), a new server can be setup in a few seconds. The command is designed in a way that allows creating arbitrary test/staging deployments (using different subdomains, e.g. if your system lives at `example.com`, you can reserve `www.example.com` for production, and deploy staging to `staging.example.com`). Eventually, it will support private deployments (where you would need a secret key to access the server, if you want your test environments to be completely safe from view), "deploying" (restoring from) specific database backups and running multiple services (e.g. the web server and any number of message queue workers). Taking those regular database backups and storing them securely outside the server is also something that should be part of the installation.

Because the infrastructure is so simple, a [postgresql database backup][backup] is the only thing needed to get back on track after an incident, everything else is in the source code repository.

Critical to a good operations story is the observability of the system. I mentioned in the architecture section the importance of logging and capturing metrics, this is where this information is put to good use. Having a single place to look for system health monitoring (often in the form of a dashboard) is a must, and so is good alerting. There are lots of books and posts on this subject, I won't go into much details, but this is one area where it might make sense to use a third-party SaaS platform to provide clear dashboards, flexible alerting and powerful insights when issues arise. The tulip framework outputs structured logs in a standard way (typically to stdout, but if running in systemd, it would likely be logging to journald), and supports [statsd-compatible metrics][statsd]. Most monitoring platform have highly flexible log and metrics collectors/extractors/scrapers/converters, making it relatively easy to get the data into their system. If you really want to handle this without a third-party platform, then you may want to look at [Prometheus][prom], [Grafana][graf] and [Grafana Loki][loki], but I feel that's more complexity than handing the data over to e.g. [DataDog][dd].

One might argue that a minimal and simple approach would be to use our postgresql database for this too, leveraging its text-search capability for logs querying, and creating a time-series-like table for metrics. There are two problems with this - first, I think that powerful visualization of the system's health is important, so the storage is only one part of this, and two, the whole point of this is to be able to be alerted and to query logs if e.g. the database itself is down.

Investigating incidents is simplified by the straightforward infrastructure. There may be problems with the server's disk space, load, networking, etc. If the server is fine, there may be issues with the postgresql database, or the various application services that make up the system. I think it is worth having a script - again, ideally implemented with the same programming language - that assists in identifying issues. I have implemented this in the past, the idea is to take a very deliberate, step-by-step validation of an as-exhaustive-as-possible list of components, e.g.:

* check that localhost (the computer used to run the script) has networking access
* check that connecting to the remote server via ssh works
* check disk space, load (maybe output the `top` command), etc. on the remote server
* ping the domain to ensure DNS resolves correctly
* make basic http/https/http 1.1/http 2 requests to the domain
* check the status of all services involved (e.g. in my case `systemctl status`, which also prints the last few lines of the relevant logs)
* connect to the database and check load/process list (e.g. information from `pg_stat_activity`, `pg_stat_database`, maybe run [pg\_top][pgtop])
* check the status page of your VPS provider

It is worth taking the time to polish this tool and its output, as it automates (and encodes directly in the repository, so that this information can be learned) many steps that would be required to check in case of incidents, when stress levels are high, and can quickly pinpoint the source of an issue. If some services have failed for some reason, the deploy script mentioned earlier can be used to restart all services on an existing infrastructure, and if things are truly bad beyond repair, a deployment on a new infrastructure can be done in seconds (maybe from a database backup, maybe from an older version of the code).

The last point I mentioned in operations is the scaling scenarios. It can obviously scale vertically, depending on what your VPS provider supports (that is, more RAM, more disk space, more/faster CPUs). However, Lua supports concurrency but not parallelism, so it won't take advantage of multiple CPUs in a given process (but it may help if you have multiple processes in your system, like a web server and a few message queue workers, in addition to the database engine).

Depending on where the bottleneck is, the various parts of the infrastructure (database, worker process, web server process) can be installed to their own separate machine. Multiple web server processes can be executed on a machine behind a reverse proxy to take advantage of all CPUs. There are a number of ways that simple and minimal infrastructure and architecture can scale without adding too much extra components and complexity - the key here is to adjust the tooling accordingly, so the actual commands stay simple and the scripts keep on encoding that knowledge. Most importantly, investigate, replicate, isolate, benchmark, analyze and _understand_ the bottlenecks before applying a solution.

## In Conclusion

I believe there is tremendous value in keeping things _simple_ and having everyone in the team being able to have a complete mental model of the whole system. I think complexity gets out of hand because we are collectively very bad at judging the cost (across all relevant lenses, not just the development aspect, but deployment, observability, etc.) of adding components in a system. It often looks like a small, harmless (actually, helpful!) and reasonable thing to do when all changes are taken in isolation. We are also often too quick to point to "new tool!" to solve a problem when what really solves it might be "new concept!", and there may be simpler ways to integrate that concept in our system. It is probably some orders of magnitude harder to _remove_ complexity than it is _adding_ it.

This small architecture is a starting point, not an end goal - the end goal should be what you build with it and how it improves people's lives! So adjust and adapt it to your needs - it's worth repeating that the _concepts_, not the specific _technologies_, are what matters in this post. As a general rule, strive to use as much as possible of what you control (the database, your code, the OS - especially if it is an open-source one) and as little as required of what you don't (third-party SaaS, VPS provider, etc. - making it easy to swap them out as no exclusive or uncommon feature is core to your system). Finally, a warning: be prepared to (respectfully) fight hard to keep things simple - as almost everyone else (often unknowingly) will fight *against* it.

[sine]: https://www.infoq.com/presentations/Simple-Made-Easy/
[ss]: https://andrewskurka.com/stupid-light-not-always-right-or-better/
[fb]: https://softwareengineeringdaily.com/2019/07/15/facebook-php-with-keith-adams/
[tw]: https://blog.twitter.com/engineering/en_us/a/2013/new-tweets-per-second-record-and-how.html
[shop]: https://shopify.engineering/deconstructing-monolith-designing-software-maximizes-developer-productivity
[luacheck]: https://github.com/mpeterv/luacheck
[le]: https://letsencrypt.org/
[tulip]: https://git.sr.ht/~mna/tulip
[cqueues]: https://github.com/wahern/cqueues
[lhttp]: https://github.com/daurnimator/lua-http
[cron]: https://github.com/citusdata/pg_cron
[pubsub]: https://www.postgresql.org/docs/13/sql-notify.html
[contract]: https://man.sr.ht/~mna/tulip/#architecture
[cicd]: https://www.redhat.com/en/topics/devops/what-is-ci-cd
[sr]: https://sourcehut.org/
[llrocks]: https://git.sr.ht/~mna/llrocks
[terra]: https://github.com/hashicorp/terraform-plugin-sdk/issues/88
[deploy]: https://git.sr.ht/~mna/tulip-cli/tree/main/item/exp/deploy.lua
[do]: https://www.digitalocean.com/
[backup]: https://www.postgresql.org/docs/current/backup-dump.html
[prom]: https://prometheus.io/
[graf]: https://grafana.com/grafana/
[loki]: https://grafana.com/oss/loki/
[dd]: https://www.datadoghq.com/
[statsd]: https://github.com/statsd/statsd/blob/master/docs/metric_types.md
[pgtop]: https://pg_top.gitlab.io/

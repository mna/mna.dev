# Linux on the Dell XPS 13 (2020)

{{template "post-meta.html" .Config}}

I recently [tweeted][tweet] about how I was tempted to buy the new Dell
XPS 13 laptop (the 2020 edition, model 9300). [Most][review1]
[reviews][review2] [praised][review3] [it][review4] as one of the top laptops
in its class, so shortly thereafter I went ahead and purchased one in late April,
and received it in May. Would it live up to the hype, especially with Linux on
it (because this is the model that comes with Windows preinstalled - not the Linux-based
"developer edition")? I was pretty excited to find out.

## First impressions

The first thing that struck me was how _tiny_ and incredibly _light_ that thing is!
I placed it beside my 13" Macbook Pro 2015 and it is _significantly_ smaller (and lighter).
However, it feels very comfortable to use, enough resting space besides the (reasonably large)
trackpad - which is itself great, not quite as _awesome_ as the macbook's, but next best
thing, I'd say -, and an _amazing_ keyboard. I just love the feel of those key presses.

The screen is also wonderful, even coming from years of staring at a retina
slab. Despite having chosen the "lower-res" version (i.e. the 1920 x 1200, not
the 4K which comes with touch-screen, which I wasn't interested in), the matte
finish makes this perfect to work outside in the sun, and the brightness range
is so wide that I barely ever use the top-half of it.

I launched Windows and tested that every hardware component worked as expected,
and saved a recovery image in case I have to reset it for some issue while on
the warranty.

## Installing Linux

At this point I was ready to wipe it off and install Linux. About a year and a
half ago, my last (!) blog post was about [how I installed Arch Linux on my
Macbook Pro][arch], where I hailed the MbP as my favorite hardware. A few
things happened since then:

* The macbook's battery started losing its capacity - quite a bit - to a point
where I have to fully charge it twice to barely make it through a day.
* I had some audio issues where the sound would cut off for a few seconds at a time.
* It's a 2015 model and starts showing its age - I'm also worried the screen might
give up at some point, like my MbP 2012 did eventually.
* It's been a great workhorse for about 5 years, but I have no interest in the later
MbP models, and if given the choice I prefer to work on Linux.

So I still use it for some work contracts where I have everything I need setup and
working already, but otherwise I switched to my Dell. So that explains the laptop
change but what about the Linux part? Which distro would get the nod?

### A slight detour

While waiting for the XPS to arrive, I tried various things on a (very old) Dell
laptop I had lying around (an Inspiron, I think). As I mentioned in my original
[tweet][], I first intended to use [Gentoo Linux][gentoo] but I quickly got tired
of compiling every component, even though it would've been much quicker on the
latest hardware, on second thought it's not the level of tweaking I want to get into.

I also tried [Manjaro][manjaro] and a couple others, which were okay until I settled
for a while on a minimal [Void Linux][void] setup, running the [Sway][] window manager
on Wayland. I had some fun configuring this and making it work to my liking, but
it also reminded me of something I may have forgotten because I hadn't tweaked my Arch
setup in a long time: it takes a lot of work, googling, trial-and-error, github-issue-reading,
and disparate package-installing to get this to a good place. You have to install
icons, a notification daemon, d-bus things to get Firefox or Google Chrome to use the
system's notification, tons of fonts to have something that looks good, etc.

Years ago, when I switched from Ubuntu to Arch, that kind of setup was what I was
_looking for_. Now that I've gained a better understanding of how it all works, I'm not
that interested in going through all those motions.

### So... Fedora

In the end, I decided to install a vanilla [Fedora][]. It is widely used, has great support
and maintenance, is kind of the "official distro" of the [GNOME][] desktop environment,
and is up-to-date enough for my taste (e.g. as of early June 2020, I run the 5.6.16 Linux kernel).

However, I'm not a fan of traditional desktop environments, I prefer the keyboard-driven window
managers, and like to run most of my windows full-screen. What pushed me to make the switch
to GNOME in full confidence was [that article about PaperWM][paper], which resonated a lot with me.
I'm now super at home in GNOME with PaperWM to manage my windows, and the small [switcher][] extension
to launch or switch to an app via the keyboard.

## XPS + Linux = 🖤

Installing Fedora Linux is very straightforward. I booted from a USB key and launched the
graphical installer, going through the options in typical fashion, all went well and was
as fast and easy as it comes.

### Power management

As I finished installing Fedora and proceeded with a reboot, my heart sank a
little as I experienced my first glitch. The computer did not restart, and in
fact - as I could tell because the keyboard backlight would still light up when
I pressed a key - it did not fully shut down.

What followed was hours-long googling, finding loads of pages mentioning
similar issues and recommending various acpi-related kernel settings to try on
the bootloader's (GRUB) kernel command line.  I removed the `quiet` and `rhgb`
settings in order to see systemd's output, and inspected the `dmesg` log on
startup, but couldn't find anything helpful.

Between each (failed) attempt, I had to shutdown and then press and hold power
for a few seconds to do a hard poweroff. I was not happy and quite concerned as
nothing I tried solved the problem, and many attempts required at least two
reboots as the documentation mentioned it would not take effect before
restarting twice. Building up hope and ending up with deception time and again.

In the end - I don't remember why exactly - I tried something I thought would
have no effect on the shutdown procedure, I disabled the UEFI Secure Boot.
Everything worked perfectly from this point on, no special kernel setting
required. As I understand it (note that I haven't read a ton about it), Secure
Boot prevents an unsigned OS from booting on your computer (e.g.  from a
malicious USB drive). That's a fairly low risk threat for my use of this
computer, and besides that's the only way I've found to make shutdown work, so
the choice was easy. If someone knows of a way to make that work while leaving
Secure Boot activated, please reach out to me and I'll gladly update the post.

Besides this, it's been pretty smooth sailing. I was happy to see that Fedora
was giving me software updates for the Dell firmware too, without me having to
manually add a Dell repository or anything.

### Battery life

I haven't closely monitored battery life under heavy use yet, but I've done a few
days where I've clearly unplugged it first thing in the morning, used it regularly
throughout the day, and got through the day without issue (up until late-ish at night).
So I think it can handle between 8 to 12 hours depending on use, though the time
remaining estimate wildly varies, hitting 17-20 hours at time.

Also I just noticed that Bluetooth was enabled and I typically have no use for it,
so I just disabled it. That should help a bit too. I installed the `powertop`
package and enabled the systemd service to run `powertop --auto-tune` at startup.
Maybe using and tweaking `tlp` could get better results, but I'm happy with
how well it fares as is.

Closing the lid properly goes to sleep (to `s2idle` level) and resumes correctly
and promptly on lid open, wifi and all. I left it unplugged a few times overnight
on sleep, and I believe it dropped something like 10%, maybe 1% per hour between
when I left it and when I got back to it (again, not closely monitored).

Charging lights up a slim bar under the touchpad, which turns off when the
charging is complete (which seems pretty fast, though that's just an impression).

### Special keys

All the special keys (function keys) work as expected - volume keys, keyboard backlight,
screen brightness, print screen, cursor navigation and insert/delete. The power key
immediately puts it to sleep. I haven't tested the dual monitor key as I don't have
a USB-C to HDMI adapter yet, but I'd be surprised if that specific one did not work.

As mentioned earlier, the screen brightness is very impressive. The screen is crisp and
beautiful at well below 50% (even below 25%) for my old-ish eyes when inside (using my
default dark-themed terminal configuration) and I have a light-themed terminal configuration
for when I'm outside, which makes it very usable at around 50% even in sunny conditions.
Of course, thanks to the matte finish, I don't find myself staring at a mirror - that
screen is perfect to work in many conditions.

The keyboard backlight has three levels, and I believe it has a light sensor that makes
it fade out if there is enough light, or stay on otherwise, but I may be wrong. I usually
leave it off unless working in the dark outside late at night.

### Trackpad

I hinted about it earlier - the trackpad is the kind I prefer, a full surface without special
zones for left/right click - click anywhere with one finger for left, anywhere with two fingers
for right. The two-finger scroll works very well (I set it up to the "natural" scroll as on
the mac, which I prefer, and disable the "tap to click", which I hate).

Surprisingly for me, some multi-touch gestures are supported, as a 3-finger swipe left-right
will switch to the next/previous window, and up-down to switch between workspaces. I don't use
that all that much as I'm used to keyboard navigation, but that was still nice to see.

### Fonts and cursor

Even though this is not the 4K screen, this is still a high resolution display, and as a result
the fonts were a bit small for my mid-forties self. Maybe it would be acceptable for some folks,
but for me I simply enabled the "Large Text" accessibility option, and set my cursor size to
"Medium".

On Firefox, I've set a default zoom level of 110% and use an extension to set specific per-site
settings when that isn't enough.

### Ports

The XPS has 2 USB-C ports and a microSD slot. I don't think I have a microSD card around,
so I haven't tested it, but I plugged my Apple USB Keyboard and a USB key in, and both
worked perfectly. The key properly gets detected and mounted and displays an appropriate
notification. Note that my laptop came in with a USB(-A) to USB-C adapter.

### Audio & video

This is interesting, as it may be another area where there are some glitches,
although it depends on the application used. My first test for the camera was
by using GNOME's `Cheese` application, and the image was black-and-white with
lots of "snow" and some red leds were blinking on each side of the camera's
white light, so clearly something's wrong here.

However, I also tried VLC's Video Capture and Zoom (which can be installed as a
flatpak), and both work properly with a clear (if not amazing) image. I did
take Zoom to a longer test and at some point it froze for a few seconds, before
working again flawlessly for the rest of the call. It did not seem like a web
issue, so again, there might be something not quite ok with the camera, but not
enough to say it doesn't work.

Audio works without issue (tested via Zoom and with the Sound Recorder
application). I also tested Zoom with my Apple headphones with integrated mic
and it works too.

### Wifi

The wifi worked without issue and without special configuration, from the
USB-based installer to the bare-metal setup. I also tested connecting to my
iPhone's Personal Hotspot and it quickly detected the network and connected
without problem.

I also have a self-hosted Wireguard VPN with Pi-Hole ad-blocker on a Digital
Ocean instance, and I configured a keyboard shortcut to enable or disable it
quickly (with integrated notification). All that works perfectly.

## Conclusion

So besides the initial scare for the shutdown process (and keeping an eye on
the camera glitch), it's been pretty smooth sailing and I ended up with a
powerful, fast, ultra-light and tiny machine with a setup I really love and that
is super fun to work on.

Hopefully that setup can last me at least 5 years as my main machine.

[tweet]: https://twitter.com/___mna___/status/1247964015298478089
[review1]: https://www.engadget.com/2020-04-01-dell-xps-13-review-best-ultraportable.html
[review2]: https://www.techradar.com/reviews/dell-xps-13-2020
[review3]: https://www.theverge.com/2020/4/15/21221003/dell-xps-13-2020-review-core-i7-specs-features-price
[review4]: https://www.wired.com/review/dell-xps-13-2020/
[arch]: https://www.0value.com/using-arch-linux-on-a-macbook-pro
[gentoo]: https://www.gentoo.org/
[manjaro]: https://manjaro.org/
[void]: https://voidlinux.org/
[sway]: https://github.com/swaywm/sway
[fedora]: https://getfedora.org/
[gnome]: https://www.gnome.org/
[paper]: https://jvns.ca/blog/2020/01/05/paperwm/
[switcher]: https://github.com/daniellandau/switcher

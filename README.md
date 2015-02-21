# one-info
A go program, which gathers information from different web sites and rss feeds and shows it in one place.

Thread with more information about the project (in bulgarian) - http://fmi.golang.bg/topics/293

Dependencies:

* github.com/PuerkitoBio/goquery
* github.com/CzarekTomczak/cef2go
* github.com/gorilla/mux
* github.com/mattn/go-sqlite3
* github.com/peterbourgon/diskv
* github.com/gregjones/httpcache
* github.com/gregjones/httpcache/diskcache

Notes:

* cef2go might require some [additional steps](https://github.com/CzarekTomczak/cef2go#getting-started-on-windows) to work
* Although on theory the application may work on different operating systems it is being developed solely on windows. Support for other OS might be added later, but currently isn't a priority.
* The API is currently in developement and may go through some backwards incompatible changes 

Future developement ideas:

* Create advanced map and filter syntax
* Create advanced modules with client and/or server side logic
* Create web site for registering and loading modules
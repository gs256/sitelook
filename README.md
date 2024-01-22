# sitelook

**Simple frontend for Google search engine written in Go.**

It is meant to be fast and secure (leaks no tracking data to Google). It is also completely server-side rendered and works without JavaScript.

![thumbnail](dev/home.jpg)

## How Does It Work?

Application parses Google search page and renders a new page with the results. It can be deployed locally or to a remote server for better privacy.

## Stack

-   Go
-   Gin as a backend framework
-   Goquery for parsing html
-   Bootstrap for styling
-   Docker (WIP)

## Features

### Search Filters

Supported search filters are `All`, `Images` and `Videos`. Other filters most likely won't be implemented.

![](dev/search.jpg)
![](dev/image-search.jpg)
![](dev/video-search.jpg)

### Query parameters

Parameters' names are identical to Google's

-   `q` - search term
-   `start` - search offset
-   `tbm` - search type
    -   `tbm=isch` - image search
    -   `tbm=vid` - video
-   `lr` - search language (e.g. `lang_en`)
-   `hl` - interface language (e.g. `en`)

### Upcoming Features

You can find all upcoming and considered features in the project's [todo.md](dev/todo.md) file.

## Issues

### Captcha Issue

Google sometimes requires captcha to make a search request. Solving captchas on the frontend is not implemented and probably won't be implemented since I'm not sure if it is even possible. Instead you will be offered to proceed with the Google Search.

![captcha-error-example](dev/captcha-error-example.png)

## Licence

MIT

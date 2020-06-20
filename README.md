# titleparser

[![CircleCI](https://circleci.com/gh/lepinkainen/titleparser.svg?style=svg)](https://circleci.com/gh/lepinkainen/titleparser)

A golang implementation of an URL title parser running in AWS Lambda

- Fetches both the `<title>` element and Opengraph title (`<meta property="og:title" content="Title" />`), preference on the latter as it is usually less spammy.

Custom parsers for:
- Ylilauta
- Imgur
- IMDB (via OMDB)
- Twitter
- Hackernews (via API)

## TODO

Custom parsers for different sites, lifted from [Pyfibot's custom title parsers](https://github.com/lepinkainen/pyfibot/blob/master/pyfibot/modules/module_urltitle.py)

# titleparser

![test workflow](https://github.com/lepinkainen/titleparser/actions/workflows/workflow.yaml/badge.svg)

A Golang implementation of an URL title parser running in AWS Lambda

- Fetches both the `<title>` element and Opengraph title (`<meta property="og:title" content="Title" />`), preference on the latter as it is usually less spammy.

Custom parsers for:

- Imgur
- IMDB (via OMDB)
- Hackernews (via API)

## TODO

Custom parsers for different sites, lifted from [Pyfibot's custom title parsers](https://github.com/lepinkainen/pyfibot/blob/master/pyfibot/modules/module_urltitle.py)

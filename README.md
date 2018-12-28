[![Build Status](https://travis-ci.org/lepinkainen/titleparser.svg?branch=master)](https://travis-ci.org/lepinkainen/titleparser)

# titleparser

A golang implementation of an URL title parser running in AWS Lambda

- Fetches both the `<title>` element and Opengraph title (`<meta property="og:title" content="Title" />`), preference on the latter as it is usually less spammy.


## TODO:

Custom parsers for different sites, lifted from [Pyfibot's custom title parsers](https://github.com/lepinkainen/pyfibot/blob/master/pyfibot/modules/module_urltitle.py)

# Mind-Renewal Bible Study CLI

This command-line app was written to help me renew my mind to truths in scripture by studying the bible with a command-line-interface.

For most people, instead of this app, I would recommend the free https://www.stepbible.org with the ESV translation.  Like this, there are no commentaries, but it has many easily viewed cross-references and Greek and Hebrew helps.  Biblehub.com also has helps for the original languages

## Features

* Lookup single or multiple verses in the ESV translation.
* Show the ESV words next to the Strongs Greek or Hebrew translation numbers.
* Lookup definitions of Strongs translation numbers.
* Search for other verses that use a given Strongs number.
* Display declarations, which are verses that you have personalized to help you renew your mind to the truths inside.
* Print your declarations for offline review and study.

## Declarations

A declaration is a verse that you have personalized to help you renew your mind.  For example, Philippians 3:7 says 

    "But whatever gain I had, I counted as loss for the sake of Christ." 

and I personalized it as a declaration with: 

    "I count as a loss the things that I once thought mattered."

All Scripture quotations, unless otherwise indicated, are taken from The Holy Bible, English Standard Version. Copyright ©2001 by Crossway Bibles, a publishing ministry of Good News Publishers.

## Data comes from

* [English Standard Version API](https://api.esv.org/)
* ESV Translation File:
  * [Tyndale House, Cambridge](http://www.TyndaleHouse.com)
  * [github - raw data](https://github.com/tyndale/STEPBible-Data) 
* [Open Scriptures](https://github.com/openscriptures/strongs) - Strongs Greek and Hebrew definitions

## Examples

![Examples](https://github.com/joncrlsn/biblestudy/raw/master/images/biblestudy-cli.png)

## To Do

* Add pre-compiled executables for download
* Translate command:
  * add cross references from openbible.info
  * add translation next to each strongs number.

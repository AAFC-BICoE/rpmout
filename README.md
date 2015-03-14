rpmout
======

rpmout is a utility for extracting RPM packages and (optionally) R packages that are installed on a Linux system.
It is written in Go.

It extracts the RPM tag info (using the rpm command) and optionally, R packages, and produces an HTML list fragment (default), JSON, simple text, and LaTeX.
It can be restricted to the RPMS that are implicated in a particular set of directories.

The LaTeX version creates three indexes of packages:
1. Index by software package name
2. Index by RPM Group, package name
3. Index by License, package name

The LaTeX version creates table of counts by license, useful for auditing purposes.

My use is to produce a list of bioinformatics applications installed as RPMS on a Rocks cluster http://en.wikipedia.org/wiki/Rocks_Cluster_Distribution

For example: users want to know what is installed in the bioinformatics install dir /opt/bioinformatics, and 'rpmout' generates (by default) an HTML fragment made up of a list of rpms and their useful attributes.  
This fragment is meant to be embedded into a static HTML page that wraps it with the appropriate local style, titles it, etc.

###Usage###
     rpmout <args> <rootDir0>...<rootDirN>
     default <rootDir>=/
     Args:
     -R=false: Find R packages (only works with LaTeX right now)
     -outputFormat="html": Values: html|json|txt|latex
     -header="This is a header"; default: "Installed Software"
     Example:  rpmout -outputFormat=html /opt /usr/local

###LaTeX###

**NB**: The LaTeX output right now is the only one that includes R packages.

To generate a PDF from the LaTeX output:
Right now, for the LaTeX file produced by rpmout, rpmoutlatex2pdf.sh, and rpmout.R all need to be in the same directory (to be fixed):
```
    # generate LaTeX for all RPMs and all R packages system:
    $ ./rpmout -R -outputFormat=latex > sample.tex
    # generate PDF (takes a minute or so...)
    $ ./rpmoutlatex2pdf.sh sample.tex
    # display PDF file
    $ acroread all.pdf
```

###Dependencies###
* The 'rpm' program (http://www.rpm.org/max-rpm/rpm.8.html) needs to be in your PATH
* If you want 'R' packages, the 'R' command needs to be in your PATH
* `rpmoutlatex2pdf.sh` needs a reasonable modern instance of 'tex/LaTeX' installed
  * The LaTeX packages: 
color,
datetime,
fancyhdr,
hyperref,
longtable,
makeidx,
microtype,
savetrees,
seqsplit,
splitidx.

*NB*: The splitindex's  `splitnindex.pl` needs to be in the `PATH`.

rpmout is a 64bit compiled on Fedora 18 binary, go version go1.3 linux/amd64

###Sample output###

####HTML####
[sample.html.gz](https://github.com/gnewton/rpmout/blob/master/sample.html.gz) is an example, from running 'rpmout /' on my Fedora 18 laptop. As it is looking for all rpms, it is a big page (~1.4MB ungziped).

####LaTeX####
[sample.tex.gz](https://github.com/gnewton/rpmout/blob/master/sample.tex.gz) is an example, from running 'rpmout -R -outputFormat=latex /' on my Fedora 18 laptop. As it is looking for all rpms, it is a big document. The PDF is sample.pdf.gz, has 545 pages and is  ~1.8MB

###Idea###
The original single threaded Ruby version I prototyped takes about 4 1/2 minutes to run. This naively written Go implementation takes <22 seconds to do the same thing.
The reason for creating `rpmout` is to generate a list of packages for the [Rocks cluster](http://www.rocksclusters.org) I manage at Agriculture and Agri-Food Canada.

###TODO###

* Move to Go package html/template, and allow the user to supply an arbitrary template for HTML output
* testing
* more idiomatic Go
* user selection of rpm tags, beyond the defaults?
* list dependencies?
* list provides?
* show location of any executables?
* show location of any libraries? 

Copyright, License, Attribution& Acknowledgements
=====
* Copyright 2014 Government of Canada
* MIT License (See LICENSE file)
* Author: '''Glen Newton''' glen.newton@gmail.com
* Originally developed at: Microbial Biodiversity Bioinformatics Group @ Agriculture and Agri-Food Canada
* Ongoing developement by GNewton


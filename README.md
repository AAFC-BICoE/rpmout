rpmout
======

rpmout is a utility for creating user facing rpm information.

It extracts the rpm tag info (using the rpm command) and produces an HTML list fragment (default), JSON, simple text, and LaTeX.
It can be restricted to the rpms that are implicated in a particular set of directories.

My use is to produce a list of bioinformatics applications installed as RPMS on a Rocks cluster http://en.wikipedia.org/wiki/Rocks_Cluster_Distribution

For example: users want to know what is installed in the bioinformatics install dir /opt/bioinformatics, and 'rpmout' generates (by default) an HTML fragment made up of a list of rpms and their useful attributes. 
This fragment is meant to be embedded into a static HTML page that wrapes it with the appropriate local style, titles it, etc.

###Usage###
     rpmout <args> <rootDir0>...<rootDirN>
     default <rootDir>=/
     Args:
     -outputFormat="html": Values: html|json|txt|latex
     Example:  rpmout -outputFormat=html /opt /usr/local

###Misc###
Note that the 'rpm' program (http://www.rpm.org/max-rpm/rpm.8.html) needs to be in your PATH

rpmout is a 64bit compiled on Fedora 18 binary, go version devel +d58997478ec6 Mon Apr 08 00:09:35 2013 -0700 linux/amd64

###Sample output###

####HTML####
[sample.html.gz](https://github.com/gnewton/rpmout/blob/master/sample.html.gz) is an example, from running 'rpmout /' on my Fedora 18 laptop. As it is looking for all rpms, it is a big page (~1.4MB ungziped).

####LaTeX####
[sample.tex.gz](https://github.com/gnewton/rpmout/blob/master/sample.tex.gz) is an example, from running 'rpmout -outputFormat=latex /' on my Fedora 18 laptop. As it is looking for all rpms, it is a big document. The PDF is sample.pdf.gz, has 700 pages and is  ~1.3MB

###Idea###
The original single threaded Ruby version I prototyped takes about 4 1/2 minutes to run. This naively written Go implementation takes <22 seconds to do the same thing..

###TODO###

* Move to Go package html/template, and allow the user to supply an arbitrary template for HTML output
* testing
* more idiomatic Go
* user selection of rpm tags, beyond the defaults?
* list dependencies?
* list provides?
* show location of any executables?
* show location of any libraries? 
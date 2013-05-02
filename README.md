rpmout
======

rpmout produces a list of the rpm information of the installed packages found in one or a list of, directories. It is oriented to sites that offer a set of rpms to a set of clients.
In my case, it is a Rocks cluster http://en.wikipedia.org/wiki/Rocks_Cluster_Distribution offering various bioinformatics rpms.

Users want to know what is installed in the bioinformatics /opt, and 'rpmout' generates (by default) an HTML fragment made up of a list of rpms and their useful attributes.
This fragment is meant to be embedded into a static HTML page that puts it into the appropriate local style, titles it, etc.

rpmout can also generate json, plain text (mostly for debugging) and a LaTeX document with all the information in a (usually) large, page-spanning table.

Usage of rpmout:
	 rpmout <args> <rootDir0>...<rootDirN>
	 default <rootDir>: /
Args:
  -outputFormat="html": Values: html|json|txt|latex

Example:  rpmout -outputFormat=html /opt /usr/local

Note that the 'rpm' program (http://www.rpm.org/max-rpm/rpm.8.html) needs to be in your path

rpmout is a 64bit compiled on Fedora 18 binary, go version devel +d58997478ec6 Mon Apr 08 00:09:35 2013 -0700 linux/amd64

sample.html.gz is an example, from running 'rpmout /' on my Fedora 18 laptop. As it is looking for all rpms, it is a big page (~1.4MB ungziped).
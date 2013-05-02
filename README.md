rpmout
======

rpmout produces the lsit of rpms that install artifacts into a list of directories. It is oriented to sites that offer a set of rpms to a set of clients.
In my case, it is a Rocks cluster http://en.wikipedia.org/wiki/Rocks_Cluster_Distribution offering various bioinformatics rpms.
Users want to know what is installed in the bioinformatics /opt, and 'rpmout' generates (by default) an HTML fragment made up of a list of rpms and their useful attributes.

rpmout can also generate json, plain text (mostly for debugging) and a LaTeX document with all the information in a (uusally) large, page spanning table.

Usage of ./rpmout:
	 ./rpmout <args> <rootDir0>...<rootDirN>
	 default <rootDir>: /
Args:
  -outputFormat="html": Values: html|json|txt|latex

Example:  ./rpmout -outputFormat=html /opt /usr/local

Note that the 'rpm' program (http://www.rpm.org/max-rpm/rpm.8.html) needs to be in your path
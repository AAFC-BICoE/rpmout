#!/usr/bin/env Rscript

options(repos="http://probability.ca/cran/")
if("RJSONIO" %in% rownames(installed.packages(lib='/tmp')) == FALSE) {install.packages("RJSONIO",lib='/tmp')}
library("RJSONIO",lib='/tmp')
packages = installed.packages()
json = asJSVars(rdata=packages)
json = substring(json, 8, nchar(json)-2)
cat(json)
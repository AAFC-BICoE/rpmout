#!/usr/bin/Rscript

options(repos="http://probability.ca/cran/")
if("RJSONIO" %in% rownames(installed.packages()) == FALSE) {install.packages("RJSONIO")}
library("RJSONIO")
a = installed.packages(fields = c("URL","Title","Description"))
#print(asJSVars(myMatrix=a))
json = asJSVars(rdata=a)
json = substring(json, 8, nchar(json)-2)
cat(json)
#write(json, stdout())

#!/bin/env bash

texFile=$1
pdflatex $texFile


base=${texFile%.*}
echo "hello" $base

echo "headings_flag 1" > ${base}.ist
echo "heading_prefix \"\\\\textbf\{\\\\color\\{blue\\}\\{\"" >> ${base}.ist
echo "heading_suffix \"\\}\\}\"" >> ${base}.ist
splitindex.pl ${base}
makeindex -s ${base}.ist ${base}-s.idx
makeindex -s ${base}.ist ${base}-g.idx
makeindex -s ${base}.ist ${base}-l.idx
pdflatex $texFile
splitindex.pl ${base}
makeindex -s ${base}.ist ${base}-s.idx
makeindex -s ${base}.ist ${base}-g.idx
makeindex -s ${base}.ist ${base}-l.idx
pdflatex $texFile
pdflatex $texFile

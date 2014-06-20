#!/bin/env bash

texFile=$1
pdflatex $texFile
pdflatex $texFile

base=${texFile%.*}
echo "hello" $base

echo "headings_flag 1" > ${base}.ist
echo "heading_prefix \"\\\\textbf\{\\\\color\\{blue\\}\\{\"" >> ${base}.ist
echo "heading_suffix \"\\}\\}\"" >> ${base}.ist

makeindex -s ${base}.ist ${base}.idx
pdflatex $1

#!/usr/bin/env python
"""
Converts JSON output from rpmout into a CSV file.
"""

__version__ = "1.0"

from optparse import OptionParser
import json
import csv

def main():
    options = set_options()
    convert_json_to_csv(options)

def set_options():
  parser = OptionParser(usage="%prog [-i sample.json] [-o sample.csv]", version="%prog $version")

  parser.add_option("-i", "--input", dest="input_json", help="JSON file to parse.", default='sample.json')
  parser.add_option("-o", "--output", dest="output_csv", help="CSV file to write.", default='sample.csv')
  parser.add_option("-t", "--tags", dest="tags", help="CSV file to write.", default='name,group,license,version,buildtime,description')
  (options, args) = parser.parse_args()
  return options

def convert_json_to_csv(options):
  try:
    with open(options.input_json, 'r') as json_fp:
      json_data = json.loads(json_fp.read())

    with open(options.output_csv, "w") as csv_fp:
      csv_file = csv.writer(csv_fp)
      
      tags = options.tags.split(',')
      #
      csv_file.writerow(tags)

      for package in json_data:
        p = json_data[package]["Tags"]
        csv_row = []
        for tag in tags:
          if tag in p:
            csv_row.append(p[tag].encode('utf8').replace('\n',""))
          else:
            csv_row.append("")
        csv_file.writerow(csv_row)
  except IOError as e:
    print e
    sys.exit(1)

if __name__ == "__main__":
  main()

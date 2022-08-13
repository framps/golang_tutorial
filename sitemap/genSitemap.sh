#!/usr/bin/env bash

#######################################################################################################################
#
# Create a sitemap for a website
#
#######################################################################################################################
#
#    Copyright (c) 2022 framp at linux-tips-and-tricks dot de
#
#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.
#
#    This program is distributed in the hope that it will be useful,
#    but WITHOUT ANY WARRANTY; without even the implied warranty of
#    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#    GNU General Public License for more details.
#
#    You should have received a copy of the GNU General Public License
#    along with this program.  If not, see <http://www.gnu.org/licenses/>.
#
#######################################################################################################################

if [[ -z $1 ]]; then
  echo "Missing Website URL"
  exit 1
fi

go run sitemap.go "$1"

SITEMAP="sitemap.xml"
urlsFound=0

rm -f $SITEMAP > /dev/null

echo '<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd">' > $SITEMAP

while read line; do
  echo "   <url>" >> $SITEMAP
  echo "      <loc>$line</loc>" >> $SITEMAP
  echo "      <lastmod>2022-06-24</lastmod>" >> $SITEMAP
  echo "      <changefreq>weekly</changefreq>" >> $SITEMAP
  echo "      <priority>0.69</priority>" >> $SITEMAP
  echo "   </url>" >> $SITEMAP
  echo >> $SITEMAP
  ((urlsFound++))
done < <(sort sitemapGen.match | uniq)

echo "</urlset>" >> $SITEMAP

echo "URLs found: $urlsFound"

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

SITEMAP="sitemap.xml"
MYNAME="genSitemap"

if [[ -z $1 ]]; then
  echo "??? Missing Website URL"
  exit 1
fi

if [[ ! $1 =~ ^http[s]?:* ]]; then
	echo "??? Missing protocol (http or https)"
	exit
fi	

arch=$(uname -m)
arch_ext="arm"
[[ $arch =~ ^x86* ]] && arch_ext="x86"

if [[ ! -f ${MYNAME}_${arch_ext} ]] || [[ $MYNAME.go -nt ${MYNAME}_${arch_ext} ]]; then # check if source code was updated or does not exist
   if ! which go >/dev/null ; then 							# no go environment detected
     echo "--- Downloading executable ${MYNAME}_${arch_ext} from github ..."	# download code from github
	 curl -q -o ${MYNAME}_${arch_ext} https://raw.githubusercontent.com/framps/golang_tutorial/master/$MYNAME/${MYNAME}_${arch_ext}
	 rc=$?
	 if [[ $rc != 0 ]]; then
		echo "??? Download of executable ${MYNAME}_${arch_ext} from git failed with curl rc $rc"
		exit 1
	 fi
	 echo "--- Downloaded $MYNAME"
	 chmod +x ${MYNAME}_${arch_ext}
   else
	 echo "--- Compiling $MYNAME"
	 OOS=linux GOARCH=arm GOARM=7 go build -o ${MYNAME}_arm $MYNAME.go
     go build -o ${MYNAME}_x86 $MYNAME.go									# otherwise build new executable
   fi
fi

echo "--- Starting crawler"
./${MYNAME}_${arch_ext} "$@"												# start crawler

if (( ! $? )); then

  echo -e "\n--- Generating $SITEMAP ..."

  urlsFound=0

  rm -f $SITEMAP > /dev/null

  echo '<?xml version="1.0" encoding="UTF-8"?>' > $SITEMAP
  echo '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd">' >> $SITEMAP

  dte=$(date +"%Y-%m-%d %H:%M:%S")

  echo "<!-- ============================================================================================================= -->" >> $SITEMAP
  echo "<!-- Sitemap generated on $dte by genSitemap available on https://github.com/framps/golang_tutorial -->" >> $SITEMAP
  echo "<!-- ============================================================================================================= -->" >> $SITEMAP

  while read line; do
    echo "   <url>" >> $SITEMAP
    echo "      <loc>$line</loc>" >> $SITEMAP
    echo "      <lastmod>2022-09-24</lastmod>" >> $SITEMAP
    echo "      <changefreq>weekly</changefreq>" >> $SITEMAP
    echo "      <priority>0.69</priority>" >> $SITEMAP
    echo "   </url>" >> $SITEMAP
    echo >> $SITEMAP
    ((urlsFound++))
  done < <(sort $MYNAME.match | uniq)

  echo "<!-- Detected URLS: $urlsFound -->" >> $SITEMAP
  echo "</urlset>" >> $SITEMAP

  echo
  echo "--- URLs added in sitemap: $urlsFound"
else
  echo -e "\n--- Sitemap generation aborted"
fi

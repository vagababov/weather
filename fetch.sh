#!/bin/bash

# The following shell script helps with pulling data for the Washington state
# and parts of Idaho.

for station in SEA OLM SEW UIL BLI HWM;
do
  fetcher --start_date=20150228 --end_date=20150331 --station=${station} --noaa_station=sew --output_dir=./data;
done

for station in GEG LWS EAT OTX SFF OMK MWH PUW EPH DEW;
do
  fetcher --start_date=20150131 --end_date=20150331 --station=${station} --noaa_station=otx --output_dir=./data;
done

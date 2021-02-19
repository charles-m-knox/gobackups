# Copyright © 2020 Intel Corporation. All rights reserved.
# SPDX-License-Identifier: BSD-3-Clause

# This Dockerfile can be used to serve out the web UI live, or to build
# a static set of distributable files
FROM node:lts-alpine

# the angular CLI (ng command) is required in order to serve & build the site
RUN npm install -g @angular/cli

# creates the internal image's directory for the web UI
RUN mkdir -p /site

WORKDIR /site

# copy only the npm package descriptors to enable more efficient
# image building (source code gets copied later)
COPY package.json /site/package.json
COPY package-lock.json /site/package-lock.json

# installs the node packages into the image's node_modules folder
RUN npm install

# copy all files after doing npm install in order to allow efficient image
# layer caching with Docker
COPY . /site

# prepare the distributable directory
RUN mkdir -p /site/dist

# this directory is intended to be mounted at runtime - the files in
# this directory will contain the compiled, distributable files for the web UI
VOLUME [ "/site/dist" ]

# port 4200 is only used when running live and is not relevant when building
# the site
EXPOSE 4200

# to serve live, change the runtime command to the following:
# CMD ["ng", "serve","--host","0.0.0.0","--port","4200","--disable-host-check"]

# this command will produce a dist/ folder in the image's working directory,
# which should be mounted at runtime so that the results can be used later
CMD ["ng", "build", "--prod"]

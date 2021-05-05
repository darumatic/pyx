#!/usr/bin/env bash

cd dist

tar -xvf pyx_1.0.4_Darwin_x86_64.tar.gz
chmod +x pyx
tar -cvzf pyx_1.0.4_Darwin_x86_64.tar.gz pyx

tar -xvf pyx_1.0.4_Linux_arm64.tar.gz
chmod +x pyx
tar -cvzf pyx_1.0.4_Linux_arm64.tar.gz pyx

tar -xvf pyx_1.0.4_Linux_i386.tar.gz
chmod +x pyx
tar -cvzf pyx_1.0.4_Linux_i386.tar.gz pyx

tar -xvf pyx_1.0.4_Linux_x86_64.tar.gz
chmod +x pyx
tar -cvzf pyx_1.0.4_Linux_x86_64.tar.gz pyx

#!/usr/bin/env bash

cd dist

tar -xvf pyx_1.0.5_Darwin_x86_64.tar.gz
chmod +x pyx
tar -cvzf pyx_1.0.5_Darwin_x86_64.tar.gz pyx

tar -xvf pyx_1.0.5_Linux_arm64.tar.gz
chmod +x pyx
tar -cvzf pyx_1.0.5_Linux_arm64.tar.gz pyx

tar -xvf pyx_1.0.5_Linux_i386.tar.gz
chmod +x pyx
tar -cvzf pyx_1.0.5_Linux_i386.tar.gz pyx

tar -xvf pyx_1.0.5_Linux_x86_64.tar.gz
chmod +x pyx
tar -cvzf pyx_1.0.5_Linux_x86_64.tar.gz pyx

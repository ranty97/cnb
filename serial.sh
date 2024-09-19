#!/bin/bash

socat -u -u pty,raw,echo=0,link=/dev/ttyS0 pty,raw,echo=0,link=/dev/ttyS1

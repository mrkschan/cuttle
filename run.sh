#!/bin/bash
find cuttle/ -type f -not -path "*_test*" | xargs --verbose gom run

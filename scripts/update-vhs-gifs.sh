#! /bin/bash

for tape in $(ls -d docs/images/vhs/*.tape); do
  vhs $tape -o $(echo $tape | sed 's|/vhs/|/|;s/\.tape/.gif/');
done

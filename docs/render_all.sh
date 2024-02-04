#!/usr/bin/env bash

for diagram in ./*.wsd
do
  echo "Render $diagram"
  cat $diagram | plantuml -tsvg -pipe > $diagram.svg
done
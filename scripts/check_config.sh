#!/bin/bash

REPOPATH=$PWD
exporterFile=$REPOPATH'/plugins/execute/transformerExporter.go'
exportedTransformers=( $( sed -n '/vdb-mcd/p' $exporterFile | cut -d' ' -f 1 | sort ) )

configFile=$REPOPATH'/environments/docker.toml'
configTransformers=( $(sed -n "/exporter\..*/p" $configFile | awk -F "." '{print $2}' | tr -d "]" | sort ) )

if [ ${#exportedTransformers[@]} != ${#configTransformers[@]} ]; then
  echo execute/transformerExporter and config contain different transformer totals
  exit 1
fi

for ((i=0; i<${#exportedTransformers[@]}; i++)); do
  if [ ${exportedTransformers[$i]} != ${configTransformers[$i]} ]; then
    echo config contains: ${configTransformers[$i]}
    echo exporter contains: ${exportedTransformers[$i]}
    echo "execute/transformerExporter and config transformers don't match"
    exit 1
  fi
done

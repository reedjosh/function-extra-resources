#!/bin/bash
crossplane beta render \
  --extra-resources example/environmentConfigs.yaml \
  --include-context \
  example/xr.yaml example/composition.yaml example/functions.yaml

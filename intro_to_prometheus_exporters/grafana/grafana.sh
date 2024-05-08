#!/bin/bash

docker run -it --rm -p 3333:3000 \
    -e GF_AUTH_ANONYMOUS_ENABLED=true \
    -e GF_AUTH_ANONYMOUS_ORG_ROLE=Admin \
    -e GF_AUTH_DISABLE_LOGIN_FORM=true \
    -e GF_PATHS_PROVISIONING=/var/lib/grafana/provisioning \
    -v $(pwd)/provisioning:/var/lib/grafana/provisioning \
    grafana/grafana:11.0.0-preview

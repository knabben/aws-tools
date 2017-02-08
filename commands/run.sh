#!/bin/bash
REDIS_URL=${REDIS_URL:-localhost:6379}
REDIS_PASS=${REDIS_PASS:-pass}

./main serve -r ${REDIS_URL} -p ${REDIS_PASS}

#!/usr/bin/env bash
set -euo pipefail

/opt/kafka/bin/kafka-storage.sh format --ignore-formatted \
  --cluster-id "$FL_KAFKA_CLUSTER_ID" --config /etc/faultline/kafka.properties
exec /opt/kafka/bin/kafka-server-start.sh /etc/faultline/kafka.properties

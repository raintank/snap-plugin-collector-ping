# snap collector plugin - ping

## Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Description
----------|-----------------------
/raintank/ping/avg | average latency in milliseconds of all pings.
/raintank/ping/min | minimum latency in milliseconds of all pings.
/raintank/ping/max | maximum latency in milliseconds of all pings.
/raintank/ping/median | median latency in milliseconds of all pings.
/raintank/ping/mdev | standard deviation of latency in milliseconds.
/raintank/ping/loss | percentage of pings that were lost.
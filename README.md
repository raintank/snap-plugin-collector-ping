# snap collector plugin - ping
This plugin collects Ping latency metrics.  

It's used in the [snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
3. [License](#license-and-authors)
4. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements
* [golang 1.5+](https://golang.org/dl/)  - needed only for building

### Operating systems
All OSs currently supported by snap:
* Linux/amd64

### Installation
#### Download ping plugin binary:
TODO

#### To build the plugin binary:
Fork https://github.com/raintank/snap-plugin-collector-ping  
Clone repo into `$GOPATH/src/github.com/raintank/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-ping.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `/build/rootfs/`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Load the plugin and create a task, see example in [Examples](https://github.com/raintank/snap-plugin-collector-ping/blob/master/README.md#examples).

## Documentation

### Collected Metrics

List of collected metrics is described in [METRICS.md](https://github.com/raintank/snap-plugin-collector-ping/blob/master/METRICS.md).

### Example
Example running ping collector and writing data to a file.

Make sure that your `$SNAP_PATH` is set, if not:
```
$ export SNAP_PATH=<snapDirectoryPath>/build/linux/x86_64
```
Other paths to files should be set according to your configuration, using a file you should indicate where it is located.


In one terminal window, open the snap daemon (in this case with logging set to 1,  trust disabled and global configuration saved in config.json ):
```
$ $SNAP_PATH/snapteld -l 1 -t 0 --config config.json
```

In another terminal window:
Load ping plugin
```
$ $SNAP_PATH/snaptel plugin load snap-plugin-collector-ping
```
See available metrics for your system
```
$ $SNAP_PATH/snaptel metric list
```

Create a task manifest file  (exemplary files in [examples/tasks/] (https://github.com/raintank/snap-plugin-collector-ping/blob/master/examples/tasks/)):
```json
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/raintank/ping/*": {}
            },
            "config": {
            	"/raintank/ping": {
            		"hostname": "127.0.0.1"
            	}
            },
            "process": null,
            "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/published_ping"
                    }
                }
            ]
        }
    }
}
```
Load file plugin for publishing:
```
$ $SNAP_PATH/snaptel plugin load $SNAP_PATH/plugin/snap-plugin-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Fri, 20 Nov 2015 11:41:39 PST
```

Create a task:
```
$ $SNAP_PATH/snaptel task create -t examples/tasks/ping-file.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

Stop previously created task:
```
$ $SNAP_PATH/snaptel task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

## License
This plugin is Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@Anthony Woods](https://github.com/woodsaj/)

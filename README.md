[![Build Status](https://travis-ci.org/blockassets/bam_agent.svg?branch=master)](https://travis-ci.org/blockassets/bam_agent)

# Block Assets Manager (BAM) Agent

This is an agent that is intended to be installed on miners to help facilitate management of them via HTTP.

Thanks to [HyperBitShop.io](https://hyperbitshop.io) for sponsoring this project.

### Usage (defaults):

``
./bam_agent -port 1111 -no-update=false
``

By default, the BAM Agent will automatically attempt to self update from the Github [latest release](https://github.com/blockassets/bam_agent/releases) tab. It chooses a random hour of the day to update. This way, if you have a number of machines, they will not all DDOS Github and your network. You can override the update behavior to not perform any updates.

### Setup

Install [dep](https://github.com/golang/dep) and the dependencies...

`make dep`

### Build binary for arm

`make arm`

### Install onto miner

The [releases tab](https://github.com/blockassets/bam_agent/releases) has `master` binaries cross compiled for ARM suitable for running on the miner. These are built automatically on [Travis](https://travis-ci.org/blockassets/bam_agent).

Download the [latest release](https://github.com/blockassets/bam_agent/releases) and copy the gunzipped `bam_agent` binary to `/usr/bin`

```
gunzip bam_agent.gz
chmod ugo+x bam_agent
scp bam_agent root@MINER_IP:/usr/bin
```

Create `/etc/systemd/system/bam_agent.service`

```
ssh root@MINER_IP "echo '
[Unit]
Description=bam_agent
After=init.service

[Service]
Type=simple
ExecStart=/usr/bin/bam_agent
Restart=always
RestartSec=4s
StandardOutput=journal+console

[Install]
WantedBy=multi-user.target
' > /etc/systemd/system/bam_agent.service"
```

Enable the service:

```
ssh root@MINER_IP "systemctl enable bam_agent; systemctl start bam_agent"
```

[![Build Status](https://travis-ci.org/lookfirst/bam_agent.svg?branch=master)](https://travis-ci.org/lookfirst/bam_agent)

# Block Assets Manager (BAM) Agent

This is an agent that is intended to be installed on miners to help facilitate management of them via HTTP.

Thanks to [HyperBitShop.io](https://hyperbitshop.io) for sponsoring this project.

### Usage (defaults):

``
./bam_agent -port 1111
``

### Setup

Install [dep](https://github.com/golang/dep) and the dependencies...

`make dep`

### Build binary for arm

`make arm`

### Install onto miner

The [releases tab](https://github.com/lookfirst/bam_agent/releases) has `master` binaries cross compiled for ARM suitable for running on the miner. These are built automatically on [Travis](https://travis-ci.org/lookfirst/bam_agent).

Download the latest release and copy the `bam_agent` binary to `/usr/bin`

```
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

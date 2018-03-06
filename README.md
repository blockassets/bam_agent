[![Build Status](https://travis-ci.org/blockassets/bam_agent.svg?branch=master)](https://travis-ci.org/blockassets/bam_agent)

# Block Assets Manager (BAM) Agent

This is an agent that is intended to be installed on miners to help facilitate management of them via HTTP as well 
as monitoring for issues internal to the miner.

Thanks to [HyperBitShop.io](https://hyperbitshop.io) for sponsoring this project.

## Running (defaults):

``
./bam_agent-linux-arm -port 1111 -no-update=false
``

By default, the BAM Agent will automatically attempt to self update from the Github 
[latest release](https://github.com/blockassets/bam_agent/releases) tab. It chooses a random hour of the day to update. 
This way, if you have a number of machines, they will not all DDOS Github and your network. You can override the update 
behavior to not perform any updates by setting `-no-update=false` when starting the agent.

Automatic updates brings security into question. You don't want someone being able to install a binary on your
machine that has not been developed in a secure fashion. We employ a strict process around this that is completely
transparent on github.

The authors have 2FA setup on their Github accounts and the project requires 2FA. All commits are required to be 
signed. All merges to master require a PR and passing unit tests. Builds are tested to ensure the binary starts up 
and serves the /status page. Travis CI automates the updates to the release page. Updates have their downloaded 
binaries hash checked. Updates are zero downtime.

### Install onto miner

The [releases tab](https://github.com/blockassets/bam_agent/releases) has `master` binaries cross compiled for ARM 
suitable for running on the miner. These are built automatically on 
[Travis](https://travis-ci.org/blockassets/bam_agent).

Download the [latest release](https://github.com/blockassets/bam_agent/releases) and copy the gunzipped 
`bam_agent-linux-arm` binary to `/usr/bin`

```
gunzip bam_agent-linux-arm.gz
chmod ugo+x bam_agent-linux-arm
scp bam_agent-linux-arm root@MINER_IP:/usr/bin
```

Create `/etc/systemd/system/bam_agent.service`

```
ssh root@MINER_IP "echo '
[Unit]
Description=bam_agent
After=init.service

[Service]
Type=simple
ExecStart=/usr/bin/bam_agent-linux-arm
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

## Building from source

* Install golang
* Install [dep](https://github.com/golang/dep)
* `git clone https://github.com/blockassets/bam_agent.git`
* `make dep`
* `make arm` (Builds binary for arm)

## API

### `GET /config/pools`

```
{"pool1": "", "pool2": "", "pool3": ""}
```

### `PUT /config/pools`

Send PUT request with json body:

```
{"pool1": "", "pool2": "", "pool3": ""}
```

Restarts cgminer.

### `GET /status`

```
{
  "agent": "39892e1 2018-03-06 02:06:09",
  "miner": "value in /usr/app/version.txt",
  "uptime": "0s"
}
```

### `GET /reboot`

Reboots the miner. Obviously be careful with this one. :-) Do note that we send HTTP expiry headers so that
if you do run this from your browser (not really advised), it won't get cached.

## Monitors

Monitors allow us to execute code periodically.
Monitors are configured by editing the `/etc/bam_agent.conf` file. This file is created when the agent first starts.

### High Load

Enabled by default. If the 5m average load is above 5, `reboot -f` the miner. This works around a bug where the load 
spikes and the miner stops submitting shares to the pool.

### Periodic Quit cgminer

Disabled by default. Periodically quit the miner app to free up memory and start fresh.

### Periodic Reboot

Disabled by default. Periodically reboot the entire miner.

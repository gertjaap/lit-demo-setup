# LIT Demo / Test environment

This repository contains everything you need to set up a simple cluster of LIT nodes
to work with various regtest coins. You can use it to demo or test multi-hop payments
and swaps.

# Installation

## Docker

This setup uses docker to run the various components, so make sure you have Docker installed
on the environment where you want to run this. You can get Docker [here](https://docs.docker.com/install/)

## Coin Daemon Images

This demo uses three coin daemons: Bitcoin, Litecoin and Dummy USD. The images need to be tagged with
`bitcoind:latest`, `litecoind:latest` and `dummyusdd:latest`. You can build them yourself, or pull them
from Docker Hub.

```
docker pull gertjaapglasbergen/lit-demo-bitcoind
docker pull gertjaapglasbergen/lit-demo-litecoind
docker pull gertjaapglasbergen/lit-demo-dummyusdd
docker tag gertjaapglasbergen/lit-demo-bitcoind bitcoind:latest
docker tag gertjaapglasbergen/lit-demo-litecoind litecoind:latest
docker tag gertjaapglasbergen/lit-demo-dummyusdd dummyusdd:latest
```

## Lit tracker

This demo contains an embedded lit tracker to track the nodes internal on the Docker network. This image 
has to be tagged `littracker:lastest`. You can build this yourself or pull it from Docker Hub:

```
docker pull gertjaapglasbergen/lit-demo-tracker
docker tag gertjaapglasbergen/lit-demo-tracker littracker:latest
```

## Lit image

You need an image in your docker images repository tagged with `lit:latest`. You can do this
by cloning the lit repository and executing `docker build`:

```
git clone https://github.com/mit-dci/lit
cd lit
docker build . -t lit
```

# Build and run environment

Once you've done these prerequisite steps, you can proceed to run the demo environment. Clone this repository
and then execute the `buildandrun.sh` script:

```
git clone https://github.com/gertjaap/lit-demo-setup
cd lit-demo-setup
./buildandrun.sh
```

Once complete you will see something like this:

```
Successfully tagged adminpanel:latest
[...]
31c637163450277e049f9600e1aef97c494e81619dadc560f125f71d71d60f12
```

The environment is now starting up. You can track its progress by issuing:

```
docker logs litdemoadminpanel -f
```

This container is booting up other containers like the coin daemons, tracker, and the central node in the
network that allows nodes to pay each other easily. Then, it will fund this central node with the majority 
of the mined coins. Once that's done the log will show a line like this:

```
INFO: 2018/10/12 08:32:28 adminpanel.go:89: Listening on port %s 8000
```

Then open `http://localhost:8000/` to view the admin panel. This panel shows the current state of the network, block heights, channel graph
but also allows you to boot new nodes.

# Starting over

If for what ever reason you want to start over, just issue `./clean.sh` from the `lit-demo-setup` folder. The data folders
will be emptied, all containers stopped. You can then restart by calling `./run.sh`. This will start new blockchains as well.
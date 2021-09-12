# cacis
Container-based Adaptive Clustering IoT System

### Requirements

- Hardware
  - Edge devices (RaspberryPi/Jetson/...)
  - Wi-Fi Adapter
  - Bluetooth Adapter
- Software
  - Go (1.16.6)
  - Snap
  - Microk8s
  - hostapd
  - netplan

### Setup
```
## Install Go
$ rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.linux-amd64.tar.gz
$ export PATH=$PATH:/usr/local/go/bin
$ go version

## Install Snapd
$ sudo apt install snapd

## Install Microk8s
$ sudo snap install microk8s --classic

## Only Master
$ sudo apt install hostapd

## Only Slave
$ sudo apt install netplan.io
```

### Default Add-on
- registry
- dashboard
- dns

### How to run Docker image

```
docker tag [local image] [localhost:32000/hoge:fuga]
docker push [localhost:32000/hoge:fuga]
```

Then, create kuberentes yaml and apply it.

```
kubectl apply -f [piyo.yaml]
```


# proxmox-network-configuration

# Install dependence library
```
apt update
apt install libpve-network-perl ifupdown2
apt install openvswitch-switch
```

# Build binary
1. Clone code from git:
```
git clone https://gitlab.skylabteam.com/Skylab-STA/proxmox-network-configuration.git
```

2. How to build:
```
make
```

# How to deploy 
1. Prepare and config file environment
- Add env variable in node at: /etc/environment
```
   PROXMOX_NW_CONFIG_HTTPS_PORT="8443"
   PROXMOX_NW_CONFIG_HTTPS_CERT="/opt/c3/server.crt"
   PROXMOX_NW_CONFIG_HTTPS_KEY="/opt/c3/server.key"
   PROXMOX_NW_CONFIG_HTTP_PORT="10000"
   PROXMOX_NW_CONFIG_GRPC_PORT="10010"
```

2. Run script deploy
```
   ./deploy_proxmox_nw_config.sh <IP> <Port>
```

Example:
```
   ./deploy_proxmox_nw_config.sh 172.17.10.21 22
```


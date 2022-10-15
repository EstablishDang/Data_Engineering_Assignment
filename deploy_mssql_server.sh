#!/bin/bash
make
node_ip=$1
port=$2
echo "Deploying proxmox network config..."
if [ -z "${node_ip}" ]; then 
	echo "Missing IP address!"
	exit
fi
proxmox_node=root@${node_ip}

echo "Create dir for service if not exists."
read -sp "Enter password: " szPassword
sshpass -p $szPassword ssh -no StrictHostKeyChecking=no -p  $port $proxmox_node 'sudo mkdir -p /opt/c3/'
echo
echo "Copy resource to promox node."
sshpass -p $szPassword scp -o StrictHostKeyChecking=no -P $port certs/* $proxmox_node:/opt/c3/
sshpass -p $szPassword scp -o StrictHostKeyChecking=no -P $port bin/proxmox_network_config $proxmox_node:/home/

#check service exists or not
echo "Deploy and start service."
sshpass -p $szPassword ssh -no StrictHostKeyChecking=no -p  $port $proxmox_node 'sudo service proxmox_network_config status'

if [ $? = 0 ]; then
	#service exists
	sshpass -p $szPassword ssh -no StrictHostKeyChecking=no -p  $port $proxmox_node 'sudo systemctl stop proxmox_network_config && sudo yes|cp -vf /home/proxmox_network_config /opt/c3/ && sudo systemctl start proxmox_network_config'
else
	#service not exists
	sshpass -p $szPassword scp -o StrictHostKeyChecking=no -P $port proxmox_network_config.service $proxmox_node:/lib/systemd/system
	sshpass -p $szPassword ssh -no StrictHostKeyChecking=no -p  $port $proxmox_node 'sudo yes|cp -vf /home/proxmox_network_config /opt/c3/ && sudo systemctl daemon-reload && sudo systemctl start proxmox_network_config'
fi
#!/bin/bash

sysctl -w net.netfilter.nf_conntrack_acct=1

kubectl get pods --all-namespaces -o jsonpath='{range .items[?(@.status.phase=="Running")]}{.status.podIP}{" "}{.metadata.namespace}{" "}{.metadata.name}{" "}{.status.hostIP}{"\n"}{end}' > /tmp/ip2pod

while [ 1 ]
do
for i in `kubectl get nodes -o jsonpath='{.items[*].status.addresses[?(@.type=="InternalIP")].address}'`
do

#check local IP
IPLOCAL=`ip addr| grep inet | grep i | wc -l` 

echo $IPLOCAL

if [ "$IPLOCAL" -eq 0 ];
then
echo 1
ssh $i  "/opt/ican/monitor/bin/statsagent > /opt/ican/monitor/tmp/conn.txt " 
scp root@$i:/opt/ican/monitor/tmp/conn.txt /opt/ican/monitor/tmp/conn.txt  
else
echo 2
/opt/ican/monitor/bin/statsagent > /opt/ican/monitor/tmp/conn.txt  
fi

./ip2poddb.py $i
echo "OK $i"

done
sleep 10
done

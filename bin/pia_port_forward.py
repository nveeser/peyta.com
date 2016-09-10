#!/usr/bin/python

import urllib
import subprocess
import json

CONFIG="/var/lib/transmission/.config/transmission-daemon/settings.json"
START="/bin/systemctl start transmission-daemon.service"
STOP="/bin/systemctl stop transmission-daemon.service"


def GetTunnelIp():
    output = subprocess.check_output(['/sbin/ifconfig', 'tun0'])
    for line in output.split('\n'):
        if line.strip().startswith("inet"):
            parts = line.split()
            return parts[1]

    raise Exception("Could not find VPN IP")


def AssignPort():
    data = { 'local_ip' : GetTunnelIp() }
    with open('/etc/openvpn/keys/user-pass.txt', 'r') as f:
        data['user'] = f.readline().strip()
        data['pass'] = f.readline().strip()

    with open('/etc/openvpn/pia_client_id') as f:
        data['client_id'] = f.readline().strip()
        
    params = urllib.urlencode(data)
    output = urllib.urlopen('https://www.privateinternetaccess.com/vpninfo/port_forward_assignment', params).read()
    
    data = json.loads(output)
    if data.has_key('error'):
        raise Exception('Error with URL: %s' % data['error'])
        
    if not data.has_key('port'):
        raise Exception('No port data in JSON: %s' % output)
    
    return data['port']        



def main():
    port = AssignPort()

    with open(CONFIG, 'r') as f:
        settings = json.load(f)

    if settings['peer-port'] != port:
        subprocess.call(STOP.split())

        settings['peer-port'] = port    
        with open(CONFIG, 'w+') as f:
            json.dump(settings, f, indent=4, sort_keys=True)
            f.write('\n')

    	with open("/tmp/OUTPUT", 'a') as f:
	    f.write("DATA: %s\n" % settings)

        subprocess.call(START.split())
    

main()


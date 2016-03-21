import urllib
import subprocess
import string
import json

CONFIG="/var/lib/transmission/.config/transmission-daemon/settings.json"

def get_ip():
    output = subprocess.check_output(['/sbin/ifconfig', 'tun0'])
    for line in output.split('\n'):
        if line.strip().startswith("inet"):
            parts = line.split()
            return parts[1]

    raise Exception("Could not find VPN IP")


def assign_port():
    data = { 'local_ip' : get_ip() }
    with open('/etc/openvpn/keys/user-pass.txt', 'r') as f:
        data['user'] = f.readline().strip()
        data['pass'] = f.readline().strip()

    with open('/etc/openvpn/pia_client_id') as f:
        data['client_id'] = f.readline().strip()
        
    params = urllib.urlencode(data)
    output = urllib.urlopen('https://www.privateinternetaccess.com/vpninfo/port_forward_assignment', params).read()
    return json.loads(output)['port']

def update_port(port, config):
    with open(config, 'r+') as f:
        s = json.load(f)
        s['peer-port'] = port    
        f.truncate(0)
        json.dump(s, f, indent=4, sort_keys=True)
        f.write('\n')
    


port = assign_port()
update_port(port, CONFIG)



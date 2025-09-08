import yaml

services = {}
for i in range(1, 51):
    services[f'node{i}'] = {
        'build': '.',
        'networks': ['mynetwork']
    }

compose_data = {
    'networks': {'mynetwork': {'driver': 'bridge'}},
    'services': services
}

with open('docker-compose.yml', 'w') as f:
    yaml.dump(compose_data, f, sort_keys=False)
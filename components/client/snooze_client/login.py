'''Module to manage login methods'''

import requests

def login(server, method, payload):
    response = requests.post('{}/api/login/{}'.format(server, method), data=payload)
    if response.get('data'):
        token = response.get('data').get('token')
        return token
    else:
        raise Exception("Could not get token")

def ldap(server, username, password):
    payload = {
        'auth': {'username': username, 'password': password},
    }
    return login(server, 'ldap', payload)

def local(server, username, password):
    payload = {
        'auth': {'username': username, 'password': password},
    }
    return login(server, 'local', payload)

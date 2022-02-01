#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

import click
import requests
from bson.json_util import loads

from snooze.cli.login import get_token

@click.group()
def snooze():
    pass

@snooze.command()
@click.option('-d', '--data', help='JSON data of the record (inline)')
def post(data):
    token = get_token()
    headers = {'Authorization': 'JWT {}'.format(token)}
    json_data = loads(data)
    response = requests.post('http://localhost:5200/api/snooze', headers=headers, data = json_data)
    print(response)

@snooze.command()
def list():
    token = get_token()
    headers = {'Authorization': 'JWT {}'.format(token)}
    response = requests.get('http://localhost:5200/api/snooze', headers=headers)
    json = response.json()
    data = json['data']
    if data:
        for record in data:
            print(record)
    else:
        print('No data')

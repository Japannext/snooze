#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Webhook input plugin for Grafana notifications'''

from logging import getLogger

from bson.json_util import loads

from snooze.api.routes import WebhookRoute
from snooze.utils.functions import sanitize

log = getLogger('snooze.webhooks.grafana')

class GrafanaRoute(WebhookRoute):
    '''A falcon route to receive Grafana alerts as input'''
    auth = {
        'auth_disabled': True
    }

    def parse(self, match, media):
        '''Parse the data of the webhook to create an alert'''
        alert = {}
        tags = match.get('tags') or {}
        alert['metric'] = match.get('metric', '')
        alert['value'] = match.get('value', '')
        alert['image_url'] = media.get('imageUrl', '')
        alert['rule_id'] = media.get('ruleId', '')
        alert['rule_url'] = media.get('ruleUrl', '')
        alert['panel_id'] = media.get('panelId', '')
        alert['dashboard_id'] = media.get('dashboardId', '')
        alert['org_id'] = media.get('orgId', '')
        alert['rule_name'] = media.get('ruleName', '')

        alert['host'] = tags.pop('host', media.get('ruleName', ''))
        alert['process'] = tags.pop('process', match.get('metric', ''))
        alert['severity'] = tags.pop('severity', 'critical')
        alert['message'] = media.get('message', media.get('title', media.get('rule_name', '')))
        alert['source'] = 'grafana'
        alert['tags'] = {}
        alert['raw'] = sanitize(media)
        for tag_k, tag_v in tags.items():
            try:
                alert['tags'][tag_k] = loads(tag_v)
            except Exception:
                alert['tags'][tag_k] = tag_v
        alert['tags'].update(media.get('tags') or {})
        alert['tags'] = sanitize(alert['tags'])

        return alert

    def parse_webhook(self, req, media):
        alerts = []
        if media.get('state', '') == 'alerting':
            for match in media.get('evalMatches', []):
                alert = self.parse(match, media)
                alerts.append(alert)
        return alerts

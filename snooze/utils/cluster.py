#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''A module for managing the snooze cluster, which is used for sharing configuration
file updates accross the cluster.'''

import socket
from logging import getLogger
from queue import Queue
from threading import Event
from typing import List, Tuple

import falcon
import requests
import pkg_resources
from netifaces import interfaces, ifaddresses, AF_INET
from requests import Request
from requests.adapters import HTTPAdapter, Retry
from requests.exceptions import RequestException

from snooze.utils.threading import SurvivingThread
from snooze.utils.config import CoreConfig, ClusterConfig
from snooze.utils.typing import HostPort, PeerStatus

log = getLogger('snooze.cluster')
ADAPTER = HTTPAdapter(max_retries=Retry(total=20, backoff_factor=0.1))

class NonResolvableHost(RuntimeError):
    '''Thrown when one member of the cluster address cnanot be resolved
    by DNS.'''
    def __init__(self, host: str):
        self.host = host
        super().__init__(f"DNS cannot resolve {host}")

class SelfNotInCluster(RuntimeError):
    '''Thrown when the running application addresses are not defined in the cluster
    configuration'''

class SelfTooMuchInCluster(RuntimeError):
    '''Thrown when the current node has too many entries of his addresses in the cluster
    configuration'''


class Cluster(SurvivingThread):
    '''A class representing the cluster and used for interacting with it.'''

    def __init__(self, core_config: CoreConfig, config: ClusterConfig, exit_event: Event = None):
        if exit_event is None:
            exit_event = Event()
        self.config = config

        self.myself = HostPort(host=socket.gethostname(), port=core_config.port)
        self.others: List[HostPort] = []

        self.schema = 'https' if core_config.ssl.enabled else 'http'

        self.initialize_config(self.config)

        self.queue = Queue()
        SurvivingThread.__init__(self, exit_event)

    def initialize_config(self, config: ClusterConfig):
        '''Validate the config and set the attributes'''
        if self.config.enabled:
            log.debug('Init Cluster Manager')
            try:
                self.myself, self.others = who_am_i(self.config.members)
            except Exception as err:
                log.exception(err)
                log.error('Error while setting up the cluster. Disabling cluster...')
                self.enabled = False

    def handle_query(self, req: Request) -> dict:
        '''Handle a request to other members of the cluster. We will not catch exceptions here
        because we want to fail if the retry doesn't work.'''
        session = requests.Session()
        session.mount(f"{self.schema}", ADAPTER)
        resp = session.send(req.prepare(), timeout=10)
        return resp.json()

    def start_thread(self):
        while True:
            req: Request = self.queue.get()
            if req is None:
                break
            self.handle_query(req)

    def stop_thread(self):
        self.queue.put(None)
        self.queue.join()

    def status(self) -> PeerStatus:
        '''Return the status, health and info of the current node'''
        version = get_version()
        status = PeerStatus(self.myself.host, self.myself.port, version, True)
        log.debug("Self cluster status: %s", status)
        return status

    def members_status(self) -> List[PeerStatus]:
        '''Fetch the status of all members of the cluster'''
        statuses = []
        statuses.append(self.status())
        if self.enabled:
            for member in self.others:
                try:
                    url = f"{self.schema}://{member.host}:{member.port}/api/cluster"
                    params = {}
                    resp = requests.get(url, params=params, timeout=10)
                    statuses.append(PeerStatus(**resp.json()))
                except RequestException:
                    statuses.append(PeerStatus(member.host, member.port, 'unknown', False))
            log.debug("Cluster members: %s", statuses)
        return statuses

    """
    def reload_plugin(self, plugin_name):
        plugins_error = []
        plugins_success = []
        log.debug("Reloading plugins %s", plugins)
        for plugin_name in plugins:
            plugin = self.core.get_core_plugin(plugin_name)
            if plugin:
                plugin.reload_data()
                plugins_success.append(plugin)
            else:
                plugins_error.append(plugin)
        if plugins_error:
            return {'status': falcon.HTTP_404, 'text': f"The following plugins could not be found: {plugins_error}"}
        else:
            return {'status': falcon.HTTP_200, 'text': "Reloaded plugins: {plugin_success}"}

    def write_and_reload(self, filename: str, conf: dict, reload: dict):
        result_dict = {}
        log.debug("Will write to %s config %s and reload %s", filename, conf, reload)
        if filename and conf:
            path = write_config(filename, conf)
            result_dict = {'status': falcon.HTTP_200, 'text': f"Reloaded config file {path}"}
        if reload:
            auth_backends = reload.get('auth_backends', [])
            if auth_backends:
                result_dict = self.reload(filename, auth_backends)
            plugins = reload_conf.get('plugins', [])
            if plugins:
                result_dict = self.reload_plugins(plugins)
    """

    def reload_plugin_others(self, plugin_name):
        '''Async function to ask other members to reload the configuration of a plugin'''
        for member in self.others:
            self.queue.put(RequestReload(member, plugin_name))

    def write_and_reload_others(self, filename: str, conf: dict, reload: dict):
        '''Async function to ask other members to update their configuration and reload'''
        for member in self.others:
            self.queue.put(RequestWriteAndReload(member, filename, conf, reload))

class RequestReload(Request):
    '''Request another member to reload a given plugin'''
    def __init__(self, member: HostPort, plugin_name: str):
        url = f"{member.host}:{member.port}/api/reload"
        payload = {'reload': plugin_name}
        Request.__init__('POST', url, json=payload)

class RequestWriteAndReload(Request):
    '''Request another member to rewrite a config'''
    def __init__(self, member: HostPort, filename: str, conf: dict, reload):
        url = f"{member.host}:{member.port}/api/reload"
        payload = {'filename': filename, 'conf': conf, 'reload': reload}
        Request.__init__('POST', url, json=payload)

def who_am_i(members: List[HostPort]) -> Tuple[HostPort, List[HostPort]]:
    '''Return which member of the cluster the running program is.
    Raise exceptions in the following cases:
    * NonResolvableHost: if one member of the cluster has its DNS not resolvable
    * SelfNotInCluster: if the running node cannot be found in the cluster
    '''
    my_addresses = [
        link.get('addr') for interface in interfaces()
        for link in ifaddresses(interface).get(AF_INET, [])
    ]
    matches = []
    for member in members:
        try:
            host = socket.gethostbyname(member.host)
            if host in my_addresses:
                matches.append(member)
        except socket.gaierror as err:
            raise NonResolvableHost(member.host) from err
    if len(matches) == 1:
        myself = matches[0]
        return myself, [x for x in members if x != myself]
    elif len(matches) == 0:
        raise SelfNotInCluster()
    else:
        raise SelfTooMuchInCluster()

def get_version() -> str:
    '''Return the version of the installed snooze-server. Return 'unknown' if not found'''
    try:
        return pkg_resources.get_distribution('snooze-server').version
    except pkg_resources.DistributionNotFound as err:
        log.exception(err)
        return 'unknown'

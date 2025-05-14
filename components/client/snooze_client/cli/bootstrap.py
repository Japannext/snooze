'''Bootstrap subscripts for snooze'''

import click

from snooze_client.cli.utils import pass_server

@click.group()
@click.pass_context
def bootstrap(ctx, **kwargs):
    '''Group of commands related to bootstrapping things'''

@click.command()
@click.argument('collection')
@pass_server
def auditlogs(server, collection):
    '''Bootstrapping audit logs (from snooze 1.4 onwards)'''
    server.bootstrap_auditlogs(collection)

bootstrap.add_command(auditlogs)

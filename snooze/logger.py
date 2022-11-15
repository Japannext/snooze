'''Logging configuration'''

from logging import Filter

class ApiFilter(Filter):
    '''Inject the API context to the log record'''
    def filter(self, record):
        record.plugin = getattr(record, 'plugin', '')
        record.uid = getattr(record, 'uid', '')
        return True

class ProcessFilter(Filter):
    '''Inject the processing context to the log record'''
    def filter(self, record):
        record.plugin = getattr(record, 'plugin', '')
        record.rid = getattr(record, 'rid', '').split('-')[0]
        return True

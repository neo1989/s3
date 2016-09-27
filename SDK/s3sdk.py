# -*- coding: utf-8 -*-
"""
    python SDK for Wiz s3

    Usage:
        c = S3sdk(URL)
        c.upload("/tmp/test.jpg")
        print c.response_data
        print c.response_code
"""

import pycurl
from cStringIO import StringIO


class Pys3sdk:

    user_agent = 'Wiz_s3sdk/1.0'
    follow_location = True
    max_redirs = 5
    timeout = 15
    verbose = False
    encoding = 'gzip'
    cookie = ''

    def __init__(self, url, headers={}, verbose=False):

        assert url
        if type(url) == type(u''):
            url = url.encode('utf-8')
        self.url = url
        self.verbose = verbose
        self.response_code = None
        self.response_type = None
        self.response_data = None
        self.cookie = ''

        # init pycurl
        self.c = pycurl.Curl()
        self.c.setopt(pycurl.ENCODING, self.encoding)
        self.c.setopt(pycurl.URL, self.url)
        self.c.setopt(pycurl.VERBOSE, self.verbose)
        self.c.setopt(pycurl.USERAGENT, self.user_agent)
        self.c.setopt(pycurl.FOLLOWLOCATION, self.follow_location)
        self.c.setopt(pycurl.MAXREDIRS, self.max_redirs)
        self.c.setopt(pycurl.TIMEOUT, self.timeout)
        self.c.setopt(pycurl.COOKIEFILE, self.cookie)
        self.c.setopt(pycurl.COOKIEJAR, self.cookie)

        # header
        headers = map(lambda val: "%s: %s" % (val, headers[val]), headers)
        self.c.setopt(pycurl.HTTPHEADER, headers)

    def _perform(self):
        response_buffer = StringIO()
        self.c.setopt(pycurl.WRITEFUNCTION, response_buffer.write)
        self.c.perform()
        self.response_code = self.c.getinfo(pycurl.RESPONSE_CODE)
        self.response_type = self.c.getinfo(pycurl.CONTENT_TYPE)
        self.c.close()
        self.response_data = response_buffer.getvalue()
        response_buffer.close()
        return self.response_data

    def upload(self, onefile, folder=''):
        self.c.setopt(pycurl.HTTPPOST, [
            ('file', (pycurl.FORM_FILE, str(onefile))),
            ('type', folder),
            ])
        return self._perform()

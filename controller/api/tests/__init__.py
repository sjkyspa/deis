from __future__ import unicode_literals
import logging

from django.test.client import RequestFactory, Client
from django.test.simple import DjangoTestSuiteRunner
import requests
import json


# add patch support to built-in django test client

def construct_patch(self, path, data='',
                    content_type='application/octet-stream', **extra):
    """Construct a PATCH request."""
    return self.generic('PATCH', path, data, content_type, **extra)


def send_patch(self, path, data='', content_type='application/octet-stream',
               follow=False, **extra):
    """Send a resource to the server using PATCH."""
    # FIXME: figure out why we need to reimport Client (otherwise NoneType)
    from django.test.client import Client  # @Reimport
    response = super(Client, self).patch(
        path, data=data, content_type=content_type, **extra)
    if follow:
        response = self._handle_redirects(response, **extra)
    return response


RequestFactory.patch = construct_patch
Client.patch = send_patch


class SilentDjangoTestSuiteRunner(DjangoTestSuiteRunner):
    """Prevents api log messages from cluttering the console during tests."""

    def run_tests(self, test_labels, extra_tests=None, **kwargs):
        """Run tests with all but critical log messages disabled."""
        # hide any log messages less than critical
        logging.disable(logging.CRITICAL)
        return super(SilentDjangoTestSuiteRunner, self).run_tests(
            test_labels, extra_tests, **kwargs)


def mock_status_ok(*args, **kwargs):
    resp = requests.Response()
    resp.status_code = 200
    resp._content_consumed = True
    return resp


def mock_broker_stub(url):
    resp = requests.Response()
    resp.status_code = 200
    resp._content = json.dumps({
        "services": [{
            "id": "1211b57f-f1b3-4279-a4a9-bdc43503603a",
            "name": "mysql Cluster",
            "description": "A MySQL-compatible relational database",
            "bindable": "true",
            "plans": [{
                "id": "1211b57f-f1b3-4279-a4b1-bdc43503603a",
                "name": "small",
                "description": "A small shared database with 100mb storage quota "
                               "and 10 connections"
            }, {
                "id": "1211b57f-f1b3-4279-a4a9-aef43503602e",
                "name": "large",
                "description": "A large dedicated database with 10GB storage quota,"
                               " 512MB of RAM, and 100 connections",
                "free": "false"
            }],
            "dashboard_client": {
                "id": "client-id-1",
                "secret": "secret-1",
                "redirect_uri": "https://dashboard.service.com"
            }
        }]
    })
    return resp


def mock_provision(url, body):
    resp = requests.Response()
    resp.status_code = 201
    resp._content = json.dumps({
        "dashboard_url": "http://example-dashboard.com/9189kdfsk0vfnku"
    })
    return resp


def mock_binding(url, body):
    resp = requests.Response()
    resp.status_code = 201
    resp._content = json.dumps({
        "credentials": {
            "uri": "mysql://mysqluser:pass@mysqlhost:3306/dbname",
            "username": "mysqluser",
            "password": "pass",
            "host": "mysqlhost",
            "port": 3306,
            "database": "dbname"
        }
    })
    return resp


def mock_deprovision_200(url):
    resp = requests.Response()
    resp.status_code = 200
    return resp


def mock_deprovision_202(url):
    resp = requests.Response()
    resp.status_code = 202
    return resp


def mock_polling_last_operation_failed(url):
    # http://username:password@broker-url/v2/service_instances/:instance_id/last_operation
    resp = requests.Response()
    resp.status_code = 200
    resp._content = json.dumps({
        "state": "failed",
        "description": "Destroy service failed."
    })
    return resp


def mock_polling_last_operation_succeeded(url):
    resp = requests.Response()
    resp.status_code = 200
    resp._content = json.dumps({
        "state": "succeeded",
        "description": "Destroy service succeeded."
    })
    return resp


from .test_api_middleware import *  # noqa
from .test_app import *  # noqa
from .test_auth import *  # noqa
from .test_build import *  # noqa
from .test_certificate import *  # noqa
from .test_config import *  # noqa
from .test_container import *  # noqa
from .test_domain import *  # noqa
from .test_hooks import *  # noqa
from .test_key import *  # noqa
from .test_limits import *  # noqa
from .test_perm import *  # noqa
from .test_release import *  # noqa
from .test_scheduler import *  # noqa
from .test_users import *  # noqa
from .test_brokers import *  # noqa
from .test_service import *  # noqa
from .test_service_binding import *  # noqa
from .test_service_instance import *  # noqa

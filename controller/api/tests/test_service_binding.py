from __future__ import unicode_literals

import json
from django.conf import settings
from django.contrib.auth.models import User
from django.test import TestCase
from rest_framework.authtoken.models import Token
from api import broker_client
from . import mock_binding


class TestServiceBinding(TestCase):
    """ Tests brokers endpoint"""

    fixtures = ['tests.json', 'test_broker.json']

    def setUp(self):
        self.user = User.objects.get(username='autotest')
        self.token = Token.objects.get(user=self.user).key
        # provide mock authentication used for run commands
        settings.SSH_PRIVATE_KEY = '<some-ssh-private-key>'

    def tearDown(self):
        # reset global vars for other tests
        settings.SSH_PRIVATE_KEY = ''

    def test_service_binding(self):
        broker_client.binding = mock_binding
        url = '/v1/service_bindings'
        body = {
            "service_instance_id": "2909e1b9-1e70-42e6-a6e1-67d2fa81ee71",
            "app_id": "5a09a1e0-a27e-4839-928b-449310ed90e0",
            "parameters": {
                "the_service_broker": "wants this object"
            }
        }
        response = self.client.post(url, json.dumps(body),
                                    content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)

from __future__ import unicode_literals

import json
# import logging
# import mock
# import requests

from django.conf import settings
from django.contrib.auth.models import User
from django.test import TestCase
from rest_framework.authtoken.models import Token
from api import broker_client
from api.models import BrokerService, Broker
from . import mock_broker_stub


class TestBrokers(TestCase):
    """ Tests brokers endpoint"""

    fixtures = ['tests.json', 'test_broker.json']

    def setUp(self):
        self.user = User.objects.get(username='autotest')
        self.token = Token.objects.get(user=self.user).key
        # provide mock authentication used for run commands
        settings.SSH_PRIVATE_KEY = '<some-ssh-private-key>'
        # provide stub for get catalog from a broker service
        broker_client.catalog = mock_broker_stub

    def tearDown(self):
        # reset global vars for other tests
        settings.SSH_PRIVATE_KEY = ''

    def test_broker(self):
        """
        Test that a user can create, read and delete an broker
        """
        url = '/v1/brokers'
        body = {
            "name": "broker-auto-test",
            "url": "broker.example.com",
            "username": "admin",
            "password": "secretpassw0rd"
        }
        response = self.client.post(url, json.dumps(body),
                                    content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        uuid = response.data['uuid']  # noqa
        self.assertIn('uuid', response.data)
        """
        Test that when a valid broker is imported, the services and plans
        defined in the broker would be recorded
        """
        brokerIns = Broker.objects.filter(url="broker.example.com")
        brokerServices = BrokerService.objects.filter(broker=brokerIns)
        self.assertEqual(brokerServices.count(), 1)

        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)

        url = '/v1/brokers/{uuid}'.format(**locals())
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertIn('uuid', response.data)

        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 204)

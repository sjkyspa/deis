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
from . import mock_broker_stub, mock_provision, mock_binding
from api.models import BrokerService, Broker


class TestBrokers(TestCase):
    """ Tests brokers endpoint"""

    fixtures = ['tests.json']

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
            "url": "https://broker.example.com",
            "username": "admin",
            "password": "secretpassw0rd"
        }
        response = self.client.post(url, json.dumps(body),
                                    content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        uuid = response.data['uuid']  # noqa
        self.assertIn('uuid', response.data)

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

    def test_services(self):
        """
        Test that when a valid broker is imported, the services and plans
        defined in the broker would be recorded
        """
        url = '/v1/brokers'
        body = {
            "name": "broker-auto-test",
            "url": "https://broker.example.com",
            "username": "admin",
            "password": "secretpassw0rd"
        }

        response = self.client.post(url, json.dumps(body),
                                    content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        brokerIns = Broker.objects.filter(url="https://broker.example.com")
        brokerServices = BrokerService.objects.filter(broker=brokerIns)
        self.assertEqual(brokerServices.count(), 1)
        # url = '/v1/services'
        #
        # response = self.client.get(url,
        #                            HTTP_AUTHORIZATION='token {}'.format(self.token))
        # self.assertEqual(response.status_code, 200)
        # self.assertEqual(len(response.data['results']), 1)
        # self.assertIn('service_plans_url', response.data)
        # self.assertIn('name', response.data)
        # self.assertIn('id', response.data)
        #
        # url = '/v1/services/{uuid}'.format(**locals())
        # response = self.client.get(url,
        #                            HTTP_AUTHORIZATION='token {}'.format(self.token))
        # self.assertEqual(response.status_code, 200)
        # self.assertIn('uuid', response.data)
        #
        # response = self.client.delete(url,
        #                               HTTP_AUTHORIZATION='token {}'.format(self.token))
        # self.assertEqual(response.status_code, 204)

    def test_service_instance(self):
        broker_client.provision = mock_provision
        url = '/v1/service_instances'
        body = {
            "organization_guid": "org-guid-here",
            "plan_id":           "1211b57f-f1b3-4279-a4a9-bdc435936031",
            "service_id":        "1211b57f-f1b3-4279-a4a9-bdc43593603b",
            "space_guid":        "space-guid-here",
            "parameters":        {
                "parameter1": "value"
            }
        }
        response = self.client.post(url, json.dumps(body),
                                    content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)

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

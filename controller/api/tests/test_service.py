from __future__ import unicode_literals
import json
from django.conf import settings
from django.contrib.auth.models import User
from django.test import TestCase
from rest_framework.authtoken.models import Token
from api.models import BrokerService, Broker
from api import broker_client
from . import mock_broker_stub


class TestService(TestCase):
    """ Tests brokers endpoint"""

    fixtures = ['tests.json', "test_broker.json"]

    def setUp(self):
        self.user = User.objects.get(username='autotest')
        self.token = Token.objects.get(user=self.user).key
        # provide mock authentication used for run commands
        settings.SSH_PRIVATE_KEY = '<some-ssh-private-key>'
        broker_client.catalog = mock_broker_stub

    def tearDown(self):
        # reset global vars for other tests
        settings.SSH_PRIVATE_KEY = ''

    def test_services_list_retrieve_delete(self):
        """
        Test that when a valid broker is imported, the services and plans
        defined in the broker would be recorded
        """
        response = self.client.get('/v1/services',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        body = {'count': response.data['count'], 'results': response.data['results']}
        self.assertEqual(body['count'], 1)
        service = body['results'][0]
        self.assertEqual(service['name'], 'mysql')
        self.assertEqual(service['bindable'], True)

        plans = service['plans']
        self.assertEqual(len(plans), 2)
        self.assertEqual(plans[0]['name'], 'small')
        self.assertEqual(body['results'][0]['name'], 'mysql')
        self.assertEqual(body['results'][0]['bindable'], True)

        uuid = service['uuid']
        url = '/v1/services/{uuid}'.format(**locals())
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertIn('uuid', response.data)

        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 204)

    def test_services_creation_from_broker(self):
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

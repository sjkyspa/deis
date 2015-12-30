from __future__ import unicode_literals
from django.conf import settings
from django.contrib.auth.models import User
from django.test import TestCase
from rest_framework.authtoken.models import Token
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
        self.assertEqual(body['count'], 2)
        service = body['results'][0]
        self.assertEqual(service['name'], 'mysql')
        self.assertEqual(service['bindable'], True)

        plans = service['plans']
        self.assertEqual(len(plans), 2)
        self.assertEqual(plans[0]['name'], 'small')

        uuid = service['uuid']
        url = '/v1/services/{uuid}'.format(**locals())
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertIn('uuid', response.data)

        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 204)

    def test_services_list_with_query_parameters(self):
        response = self.client.get('/v1/services?name=mysql',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['count'], 1)

        response = self.client.get('/v1/services?name=redis',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['count'], 0)

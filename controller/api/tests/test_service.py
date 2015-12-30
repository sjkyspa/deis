from __future__ import unicode_literals
from django.conf import settings
from django.contrib.auth.models import User
from django.test import TestCase
from rest_framework.authtoken.models import Token


class TestService(TestCase):
    """ Tests brokers endpoint"""

    fixtures = ['tests.json', "test_broker.json"]

    def setUp(self):
        self.user = User.objects.get(username='autotest')
        self.token = Token.objects.get(user=self.user).key
        # provide mock authentication used for run commands
        settings.SSH_PRIVATE_KEY = '<some-ssh-private-key>'

    def tearDown(self):
        # reset global vars for other tests
        settings.SSH_PRIVATE_KEY = ''

    def test_services(self):
        """
        Test that when a valid broker is imported, the services and plans
        defined in the broker would be recorded
        """
        response = self.client.get('/v1/services',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)

        body = {'count': response.data['count'], 'results': response.data['results']}

        self.assertEqual(body['count'], 1)
        self.assertEqual(body['results'][0]['name'], 'mysql')
        self.assertEqual(body['results'][0]['bindable'], True)

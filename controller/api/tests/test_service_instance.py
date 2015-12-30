from __future__ import unicode_literals

import json

from django.conf import settings
from django.contrib.auth.models import User
from django.test import TestCase
from rest_framework.authtoken.models import Token
from api import broker_client
from api.models import ServiceInstance
from . import mock_provision, mock_deprovision_202, \
    mock_polling_last_operation_failed,\
    mock_polling_last_operation_succeeded


class TestServiceInstances(TestCase):
    """ Tests brokers endpoint"""

    fixtures = ['tests.json', 'test_broker.json']

    def setUp(self):
        self.user = User.objects.get(username='autotest')
        self.token = Token.objects.get(user=self.user).key
        # provide mock authentication used for run commands
        settings.SSH_PRIVATE_KEY = '<some-ssh-private-key>'
        # provide stub for get provision information from a broker service
        broker_client.provision = mock_provision

    def tearDown(self):
        # reset global vars for other tests
        settings.SSH_PRIVATE_KEY = ''

    def test_service_instance_creation(self):
        url = '/v1/service_instances'
        body = {
            "plan_id":           "1211b57f-f1b3-4279-a4a9-bdc435936031",
            "parameters":        {
                "parameter1": "value"
            }
        }
        response = self.client.post(url, json.dumps(body),
                                    content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)

    def test_service_instance_list_retrieve(self):
        url = '/v1/service_instances'
        response = self.client.get(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)

        uuid = response.data['results'][0]['uuid']

        url = '/v1/service_instances/{uuid}'.format(**locals())
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertIn('uuid', response.data)
        self.assertIn('service_id', response.data)
        self.assertIn('plan_id', response.data)
        self.assertIn('organization_guid', response.data)
        self.assertIn('dashboard_url', response.data)

    def test_service_instance_deletion(self):
        broker_client.deprovisioning = mock_deprovision_202
        broker_client.polling_last_operation = mock_polling_last_operation_failed

        uuid = '2909e1b9-1e70-42e6-a6e1-67d2fa81ee71'

        url = '/v1/service_instances/{uuid}'.format(**locals())
        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 500)
        self.assertEqual(ServiceInstance.objects.filter(uuid=uuid).count(), 1)

        # delete success
        broker_client.polling_last_operation = mock_polling_last_operation_succeeded

        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(ServiceInstance.objects.filter(uuid=uuid).count(), 0)

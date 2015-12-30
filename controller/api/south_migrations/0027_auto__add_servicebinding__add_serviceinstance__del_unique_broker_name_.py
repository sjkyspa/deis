# -*- coding: utf-8 -*-
from south.utils import datetime_utils as datetime
from south.db import db
from south.v2 import SchemaMigration
from django.db import models


class Migration(SchemaMigration):

    def forwards(self, orm):
        # Removing unique constraint on 'ServicePlan', fields ['name']
        db.delete_unique(u'api_serviceplan', ['name'])

        # Removing unique constraint on 'ServicePlan', fields ['id']
        db.delete_unique(u'api_serviceplan', ['id'])

        # Removing unique constraint on 'BrokerService', fields ['description']
        db.delete_unique(u'api_brokerservice', ['description'])

        # Removing unique constraint on 'BrokerService', fields ['name']
        db.delete_unique(u'api_brokerservice', ['name'])

        # Removing unique constraint on 'Broker', fields ['name']
        db.delete_unique(u'api_broker', ['name'])

        # Adding model 'ServiceBinding'
        db.create_table(u'api_servicebinding', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('service_instance_id', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.ServiceInstance'])),
            ('app_id', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.App'])),
            ('credentials', self.gf('json_field.fields.JSONField')(default={}, blank=True)),
        ))
        db.send_create_signal(u'api', ['ServiceBinding'])

        # Adding model 'ServiceInstance'
        db.create_table(u'api_serviceinstance', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('service_id', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.BrokerService'])),
            ('plan_id', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.ServicePlan'])),
            ('organization_guid', self.gf('django.db.models.fields.TextField')()),
            ('space_guid', self.gf('django.db.models.fields.TextField')()),
            ('dashboard_url', self.gf('django.db.models.fields.TextField')(null=True, blank=True)),
        ))
        db.send_create_signal(u'api', ['ServiceInstance'])

        # Deleting field 'BrokerService.id'
        db.delete_column(u'api_brokerservice', 'id')

        # Adding field 'BrokerService.bindable'
        db.add_column(u'api_brokerservice', 'bindable',
                      self.gf('django.db.models.fields.BooleanField')(default=True),
                      keep_default=False)


        # Changing field 'BrokerService.description'
        db.alter_column(u'api_brokerservice', 'description', self.gf('django.db.models.fields.TextField')(null=True))

    def backwards(self, orm):
        # Deleting model 'ServiceBinding'
        db.delete_table(u'api_servicebinding')

        # Deleting model 'ServiceInstance'
        db.delete_table(u'api_serviceinstance')

        # Adding unique constraint on 'Broker', fields ['name']
        db.create_unique(u'api_broker', ['name'])

        # Adding field 'BrokerService.id'
        db.add_column(u'api_brokerservice', 'id',
                      self.gf('django.db.models.fields.TextField')(default=1230, unique=True),
                      keep_default=False)

        # Deleting field 'BrokerService.bindable'
        db.delete_column(u'api_brokerservice', 'bindable')

        # Adding unique constraint on 'BrokerService', fields ['name']
        db.create_unique(u'api_brokerservice', ['name'])


        # Changing field 'BrokerService.description'
        db.alter_column(u'api_brokerservice', 'description', self.gf('django.db.models.fields.TextField')(default=1230, unique=True))
        # Adding unique constraint on 'BrokerService', fields ['description']
        db.create_unique(u'api_brokerservice', ['description'])

        # Adding unique constraint on 'ServicePlan', fields ['id']
        db.create_unique(u'api_serviceplan', ['id'])

        # Adding unique constraint on 'ServicePlan', fields ['name']
        db.create_unique(u'api_serviceplan', ['name'])


    models = {
        u'api.app': {
            'Meta': {'object_name': 'App'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'id': ('django.db.models.fields.SlugField', [], {'default': "'kingly-hacienda'", 'unique': 'True', 'max_length': '64'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'structure': ('json_field.fields.JSONField', [], {'default': '{}', 'blank': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.broker': {
            'Meta': {'object_name': 'Broker'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'name': ('django.db.models.fields.TextField', [], {}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'password': ('django.db.models.fields.TextField', [], {'null': 'True', 'blank': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'url': ('django.db.models.fields.TextField', [], {'unique': 'True'}),
            'username': ('django.db.models.fields.TextField', [], {'null': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.brokerservice': {
            'Meta': {'object_name': 'BrokerService'},
            'bindable': ('django.db.models.fields.BooleanField', [], {'default': 'True'}),
            'broker': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Broker']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'description': ('django.db.models.fields.TextField', [], {'null': 'True', 'blank': 'True'}),
            'name': ('django.db.models.fields.TextField', [], {}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.build': {
            'Meta': {'ordering': "[u'-created']", 'unique_together': "((u'app', u'uuid'),)", 'object_name': 'Build'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'dockerfile': ('django.db.models.fields.TextField', [], {'blank': 'True'}),
            'image': ('django.db.models.fields.CharField', [], {'max_length': '256'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'procfile': ('json_field.fields.JSONField', [], {'default': '{}', 'blank': 'True'}),
            'sha': ('django.db.models.fields.CharField', [], {'max_length': '40', 'blank': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.certificate': {
            'Meta': {'object_name': 'Certificate'},
            'certificate': ('django.db.models.fields.TextField', [], {}),
            'common_name': ('django.db.models.fields.TextField', [], {'unique': 'True'}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'expires': ('django.db.models.fields.DateTimeField', [], {}),
            u'id': ('django.db.models.fields.AutoField', [], {'primary_key': 'True'}),
            'key': ('django.db.models.fields.TextField', [], {}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'})
        },
        u'api.config': {
            'Meta': {'ordering': "[u'-created']", 'unique_together': "((u'app', u'uuid'),)", 'object_name': 'Config'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'cpu': ('json_field.fields.JSONField', [], {'default': '{}', 'blank': 'True'}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'memory': ('json_field.fields.JSONField', [], {'default': '{}', 'blank': 'True'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'tags': ('json_field.fields.JSONField', [], {'default': '{}', 'blank': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'}),
            'values': ('json_field.fields.JSONField', [], {'default': '{}', 'blank': 'True'})
        },
        u'api.container': {
            'Meta': {'ordering': "[u'created']", 'object_name': 'Container'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'num': ('django.db.models.fields.PositiveIntegerField', [], {}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'release': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Release']"}),
            'type': ('django.db.models.fields.CharField', [], {'max_length': '128'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.domain': {
            'Meta': {'object_name': 'Domain'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'domain': ('django.db.models.fields.TextField', [], {'unique': 'True'}),
            u'id': ('django.db.models.fields.AutoField', [], {'primary_key': 'True'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'})
        },
        u'api.key': {
            'Meta': {'unique_together': "((u'owner', u'fingerprint'),)", 'object_name': 'Key'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'fingerprint': ('django.db.models.fields.CharField', [], {'max_length': '128'}),
            'id': ('django.db.models.fields.CharField', [], {'max_length': '128'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'public': ('django.db.models.fields.TextField', [], {'unique': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.push': {
            'Meta': {'ordering': "[u'-created']", 'unique_together': "((u'app', u'uuid'),)", 'object_name': 'Push'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'fingerprint': ('django.db.models.fields.CharField', [], {'max_length': '255'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'receive_repo': ('django.db.models.fields.CharField', [], {'max_length': '255'}),
            'receive_user': ('django.db.models.fields.CharField', [], {'max_length': '255'}),
            'sha': ('django.db.models.fields.CharField', [], {'max_length': '40'}),
            'ssh_connection': ('django.db.models.fields.CharField', [], {'max_length': '255'}),
            'ssh_original_command': ('django.db.models.fields.CharField', [], {'max_length': '255'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.release': {
            'Meta': {'ordering': "[u'-created']", 'unique_together': "((u'app', u'version'),)", 'object_name': 'Release'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'build': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Build']", 'null': 'True'}),
            'config': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Config']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'summary': ('django.db.models.fields.TextField', [], {'null': 'True', 'blank': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'}),
            'version': ('django.db.models.fields.PositiveIntegerField', [], {})
        },
        u'api.servicebinding': {
            'Meta': {'object_name': 'ServiceBinding'},
            'app_id': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'credentials': ('json_field.fields.JSONField', [], {'default': '{}', 'blank': 'True'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'service_instance_id': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.ServiceInstance']"}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.serviceinstance': {
            'Meta': {'object_name': 'ServiceInstance'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'dashboard_url': ('django.db.models.fields.TextField', [], {'null': 'True', 'blank': 'True'}),
            'organization_guid': ('django.db.models.fields.TextField', [], {}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'plan_id': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.ServicePlan']"}),
            'service_id': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.BrokerService']"}),
            'space_guid': ('django.db.models.fields.TextField', [], {}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.serviceplan': {
            'Meta': {'object_name': 'ServicePlan'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'description': ('django.db.models.fields.TextField', [], {'null': 'True', 'blank': 'True'}),
            'free': ('django.db.models.fields.BooleanField', [], {'default': 'True'}),
            'id': ('django.db.models.fields.TextField', [], {}),
            'name': ('django.db.models.fields.TextField', [], {}),
            'service': ('django.db.models.fields.related.ForeignKey', [], {'related_name': "u'plans'", 'to': u"orm['api.BrokerService']"}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'auth.group': {
            'Meta': {'object_name': 'Group'},
            u'id': ('django.db.models.fields.AutoField', [], {'primary_key': 'True'}),
            'name': ('django.db.models.fields.CharField', [], {'unique': 'True', 'max_length': '80'}),
            'permissions': ('django.db.models.fields.related.ManyToManyField', [], {'to': u"orm['auth.Permission']", 'symmetrical': 'False', 'blank': 'True'})
        },
        u'auth.permission': {
            'Meta': {'ordering': "(u'content_type__app_label', u'content_type__model', u'codename')", 'unique_together': "((u'content_type', u'codename'),)", 'object_name': 'Permission'},
            'codename': ('django.db.models.fields.CharField', [], {'max_length': '100'}),
            'content_type': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['contenttypes.ContentType']"}),
            u'id': ('django.db.models.fields.AutoField', [], {'primary_key': 'True'}),
            'name': ('django.db.models.fields.CharField', [], {'max_length': '50'})
        },
        u'auth.user': {
            'Meta': {'object_name': 'User'},
            'date_joined': ('django.db.models.fields.DateTimeField', [], {'default': 'datetime.datetime.now'}),
            'email': ('django.db.models.fields.EmailField', [], {'max_length': '75', 'blank': 'True'}),
            'first_name': ('django.db.models.fields.CharField', [], {'max_length': '30', 'blank': 'True'}),
            'groups': ('django.db.models.fields.related.ManyToManyField', [], {'symmetrical': 'False', 'related_name': "u'user_set'", 'blank': 'True', 'to': u"orm['auth.Group']"}),
            u'id': ('django.db.models.fields.AutoField', [], {'primary_key': 'True'}),
            'is_active': ('django.db.models.fields.BooleanField', [], {'default': 'True'}),
            'is_staff': ('django.db.models.fields.BooleanField', [], {'default': 'False'}),
            'is_superuser': ('django.db.models.fields.BooleanField', [], {'default': 'False'}),
            'last_login': ('django.db.models.fields.DateTimeField', [], {'default': 'datetime.datetime.now'}),
            'last_name': ('django.db.models.fields.CharField', [], {'max_length': '30', 'blank': 'True'}),
            'password': ('django.db.models.fields.CharField', [], {'max_length': '128'}),
            'user_permissions': ('django.db.models.fields.related.ManyToManyField', [], {'symmetrical': 'False', 'related_name': "u'user_set'", 'blank': 'True', 'to': u"orm['auth.Permission']"}),
            'username': ('django.db.models.fields.CharField', [], {'unique': 'True', 'max_length': '30'})
        },
        u'contenttypes.contenttype': {
            'Meta': {'ordering': "('name',)", 'unique_together': "(('app_label', 'model'),)", 'object_name': 'ContentType', 'db_table': "'django_content_type'"},
            'app_label': ('django.db.models.fields.CharField', [], {'max_length': '100'}),
            u'id': ('django.db.models.fields.AutoField', [], {'primary_key': 'True'}),
            'model': ('django.db.models.fields.CharField', [], {'max_length': '100'}),
            'name': ('django.db.models.fields.CharField', [], {'max_length': '100'})
        }
    }

    complete_apps = ['api']
__author__ = 'hjliu'
import requests


def catalog(url):
    response = requests.get(url)
    return response


def provision(url, body):
    response = requests.post(url, body)
    return response

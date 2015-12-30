import requests


def catalog(url):
    response = requests.get(url)
    return response


def provision(url, body):
    response = requests.post(url, data=body)
    return response


def binding(url, body):
    response = requests.post(url, data=body)
    return response


def deprovisioning(url):
    response = requests.delete(url)
    return response


def polling_last_operation(url):
    response = requests.get(url)
    return response

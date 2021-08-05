import os

import requests.auth
from gql import Client, gql
from gql.transport.requests import RequestsHTTPTransport
from prometheus_client import Summary
from python_function_model import FunctionRequest, FunctionResponse

from .apig_sdk import signer

REQUEST_TIME = Summary('request_processing_seconds', 'Time spent processing request')
dataspace_url = os.environ.get('DATASPACE_URL')
app_key = os.environ.get('APP_KEY')
app_secret = os.environ.get('APP_SECRET')


# Метод handler. Этот метод будет вызываться при вызове функции
@REQUEST_TIME.time()
def handle(request: FunctionRequest):
    return FunctionResponse(
        "Hello from Python3.8-graphql function!\n" +
        "You said: " + request.payload + "\n" +
        "GraphQL status: " + graphql_status,
        200,
        {"Content-Type": "text/plain"}
    )


# Аутентификация запроса
class DataspaceAuth(requests.auth.AuthBase):
    def __call__(self, r):
        if app_key is None or app_secret is None:
            print("APP_SECRET or APP_KEY is undefined. Request will not be signed")
            return r

        sig = signer.Signer()
        sig.Key = app_key
        sig.Secret = app_secret
        request = signer.HttpRequest(r.method, r.url, r.headers, r.body.decode('utf-8'))
        sig.Sign(request)
        r.headers = request.headers
        return r


# Инициализация GraphQl клиента
if dataspace_url is not None:
    transport = RequestsHTTPTransport(url=dataspace_url, auth=DataspaceAuth(), verify=False)
    client = Client(transport=transport, fetch_schema_from_transport=False)
    graphql_status = "Dataspace URL: " + dataspace_url
else:
    client = None
    graphql_status = "DATASPACE_URL environment variable is not set. GraphQL client disabled"
print(graphql_status)


# Пример вызова DataSpace с подписью
def call_dataspace():
    # Запрос
    query = gql("query ($paramId: ID!) { some_operation( id: $paramId) { some_field } }")
    variable_values = {
        "paramId": "paramValue"
    }
    # Вызов Dataspace
    return client.execute(query, variable_values=variable_values)

# Раскомментируйте и измените этот метод, если вам нужен кастомный HealthCheck
# def health():
#     return True, "Custom healthcheck"

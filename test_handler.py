from unittest import TestCase, main

from python_function_model import FunctionRequest, FunctionResponse

from handler import handler


class HandlerTests(TestCase):

    def test_handle(self):
        # Given
        payload = "Request message"
        request = FunctionRequest(
            payload,
            headers=None
        )

        # When
        response: FunctionResponse = handler.handle(request)

        # Then
        assert response.status == 200


if __name__ == '__main__':
    main()

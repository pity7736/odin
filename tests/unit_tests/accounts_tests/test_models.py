import uuid

from odin.accounts.domain.models import User


def test_encrypt_password():
    user = User(
        email='me@raiseexception.com',
        password='test',
        first_name='julián',
        last_name='cortés',
        id=uuid.uuid4()
    )
    assert user.password != 'test'

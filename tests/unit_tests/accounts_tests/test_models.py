from odin.accounts.models import User


def test_encrypt_password():
    user = User(
        email='me@raiseexception.com',
        password='test',
        first_name='julián',
        last_name='cortés'
    )
    assert user.password != 'test'

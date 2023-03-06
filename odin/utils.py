import base64
import secrets
import string

import edgedb
from Crypto.Protocol.KDF import scrypt


def get_random_string(length=16, allowed_chars=None):
    allowed_chars = allowed_chars or string.ascii_letters + string.digits
    return ''.join(secrets.choice(allowed_chars) for _ in range(length))


def make_password(value, salt=None):
    salt = salt or get_random_string(32)
    password_hash = scrypt(
        password=value,
        salt=salt,
        key_len=32,
        N=16384,
        r=8,
        p=1
    )
    password_encoded = base64.b64encode(password_hash)
    return f'{salt}${password_encoded.decode()}'


class DBClient:

    _instance = None

    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._client = edgedb.create_client()
        return cls._instance

    def __getattribute__(self, item):
        if item != '_client':
            return getattr(self._client, item)
        return super().__getattribute__(item)

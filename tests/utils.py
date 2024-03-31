import re

UUID_PATTERN = r'^[a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12}$'


def is_uuid(text) -> bool:
    return bool(re.match(UUID_PATTERN, text))

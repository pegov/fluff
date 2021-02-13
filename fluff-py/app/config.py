import os
from string import ascii_letters, digits

DEBUG = os.environ.get("DEBUG") in ["1", "true", "True"]

REDIS_URL = os.environ.get("REDIS_URL") or "redis://localhost/"

LINK_MIN_LENGTH = 2
LINK_MAX_LENGTH = 32
LINK_NEW_LENGTH = 6
LINK_CHARS = "".join([ascii_letters, digits])

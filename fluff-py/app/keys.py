import random

from app.config import LINK_CHARS, LINK_LENGTH


def create_random_key():
    return "".join(random.choices(LINK_CHARS, k=LINK_LENGTH))
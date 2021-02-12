import random

from app.config import LINK_CHARS, LINK_NEW_LENGTH


def create_random_key() -> str:
    return "".join(random.choices(LINK_CHARS, k=LINK_NEW_LENGTH))

from typing import Optional

from pydantic import BaseModel, validator

from app.config import LINK_CHARS, LINK_LENGTH
from app.validators import URLValidator


class Link(BaseModel):
    url: str
    key: Optional[str] = None

    @validator("key")
    def validate_key(cls, v) -> str:
        if v is None:
            return v

        v = v.strip()
        if len(v) < LINK_LENGTH:
            raise ValueError("length")

        if any(c not in LINK_CHARS for c in v):
            raise ValueError("chars")

        return v

    @validator("url")
    def validate_url(cls, v) -> str:
        validator = URLValidator()
        return validator(v)

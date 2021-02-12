from typing import Optional
from pydantic import BaseModel, validator

from app.config import LINK_CHARS, LINK_LENGTH


class Link(BaseModel):
    url: str
    key: Optional[str]

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

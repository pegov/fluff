from typing import Optional

from fastapi import APIRouter, HTTPException, Request
from fastapi.responses import RedirectResponse
from pydantic import ValidationError

from app.db import db
from app.models import Link

api_router = APIRouter()


@api_router.post("/links", name="create")
async def create(*, request: Request):
    data = await request.json()
    try:
        link = Link(**data)
    except ValidationError as e:
        msg = e.errors()[0].get("msg")
        if "Error" in msg or "Exception" in msg:
            raise HTTPException(500)
        raise HTTPException(400, detail=msg)

    if link.key is None:
        link.key = await db.get_free_key()

    if not await db.set_link(link.key, link.url):
        raise HTTPException(400)

    return {"key": link.key}


router = APIRouter()


@router.get("/{key}", name="redirect")
async def redirect(key: Optional[str] = None):
    url = await db.get_link(key)
    if key is None or url is None:
        return RedirectResponse("/")

    return RedirectResponse(url)

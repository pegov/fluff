from typing import Optional
from fastapi import APIRouter, HTTPException
from fastapi.responses import RedirectResponse

from app.models import Link

from app.db import db

api_router = APIRouter()


@api_router.post("/links", name="create")
async def create(*, link: Link):
    # TODO: check url
    if link.key is None:
        link.key = await db.get_free_key()

    if not await db.set_link(link.key, link.url):
        raise HTTPException(400)

    return None


router = APIRouter()


@router.get("/{key}", name="redirect")
async def redirect(key: Optional[str] = None):
    url = await db.get_link(key)
    if key is None or url is None:
        return RedirectResponse("/")

    return RedirectResponse(url)

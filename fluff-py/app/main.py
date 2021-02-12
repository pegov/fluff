from fastapi import FastAPI

from app.config import DEBUG
from app.events import startup_event, shutdown_event

from app.routers import api_router, router

app = FastAPI(debug=DEBUG)

app.add_event_handler("startup", startup_event)
app.add_event_handler("shutdown", shutdown_event)

app.include_router(router)
app.include_router(api_router, prefix="/api")

from app.db import db


async def startup_event():
    await db.connect()
    await db.create_initial_keys()


async def shutdown_event():
    await db.disconnect()

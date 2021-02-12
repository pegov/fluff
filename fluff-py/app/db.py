import asyncio
from typing import List, Optional

import aioredis

from app.config import REDIS_URL
from app.keys import create_random_key

EXPIRE = 60 * 60 * 12
RESOURCE = 500

lock = asyncio.Lock()

QUEUE = "queue:free"


class DB:
    async def connect(self) -> None:
        self.redis = await aioredis.create_redis_pool(REDIS_URL, encoding="utf-8")
        await self.redis.ping()

    async def disconnect(self) -> None:
        self.redis.close()
        await self.redis.wait_closed()

    async def get_link(self, key: str) -> Optional[str]:
        async with lock:
            link = await self.redis.get(key)
            if link is not None:
                await self.redis.delete(key)
            return link

    async def set_link(self, key: str, url: str, expire: int = EXPIRE) -> bool:
        res = await self.redis.setnx(key, url)
        asyncio.create_task(self.redis.expire(key, expire))
        return res

    async def is_free(self, key: str) -> bool:
        link = await self.redis.get(key)
        return bool(link)

    async def get_queue(self) -> List[str]:
        return await self.redis.lrange(QUEUE, 0, -1)

    async def create_initial_keys(self):
        queue = await self.redis.lrange(QUEUE, 0, -1)
        if len(queue) < RESOURCE:
            for _ in range(RESOURCE - len(queue)):
                await self.create_free_key()

    async def create_free_key(self) -> str:
        queue = await self.redis.lrange(QUEUE, 0, -1)
        while True:
            key = create_random_key()
            if key not in queue:
                await self.redis.rpush(QUEUE, key)
                return key

    async def get_free_key(self) -> str:
        while True:
            key = await self.redis.rpop(QUEUE)
            if key is None:
                await asyncio.sleep(0.025)
                continue
            asyncio.create_task(self.create_free_key())
            return key


db = DB()

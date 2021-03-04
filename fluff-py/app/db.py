import asyncio
from collections import deque
from typing import Optional

import aioredis

from app.config import REDIS_URL
from app.keys import create_random_key

EXPIRE = 60 * 60 * 12
RESOURCE = 100

lock = asyncio.Lock()


class DB:
    queue: deque

    async def connect(self) -> None:
        self.redis = await aioredis.create_redis_pool(REDIS_URL, encoding="utf-8")
        await self.redis.ping()
        self.queue = deque()
        await self.create_initial_keys()

    async def disconnect(self) -> None:
        self.redis.close()
        await self.redis.wait_closed()

    async def get_link(self, key: str) -> Optional[str]:
        async with lock:
            link = await self.redis.get(key)
            if link is not None:
                await self.redis.delete(key)
            return link

    async def set_link(
        self, key: str, url: str, expire: int = EXPIRE, custom=False
    ) -> bool:
        res = await self.redis.setnx(key, url)
        if res:
            if custom and key in self.queue:
                self.queue.remove(key)
                asyncio.create_task(self.create_free_key())
            asyncio.create_task(self.redis.expire(key, expire))
        return res

    async def is_free(self, key: str) -> bool:
        link = await self.redis.get(key)
        return not bool(link)

    async def create_initial_keys(self):
        for _ in range(RESOURCE):
            await self.create_free_key()

    async def create_free_key(self) -> str:
        while True:
            key = create_random_key()
            if key not in self.queue and await self.is_free(key):
                self.queue.append(key)
                return key

    async def get_free_key(self) -> str:
        key = self.queue.popleft()
        asyncio.create_task(self.create_free_key())
        return key


db = DB()

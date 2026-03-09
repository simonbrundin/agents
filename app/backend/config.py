from pydantic_settings import BaseSettings
from functools import lru_cache


class Settings(BaseSettings):
    database_url: str = (
        "postgresql+asyncpg://postgres:postgres@localhost:5432/software_agent"
    )
    secret_key: str = "dev-secret-key-change-in-production"
    algorithm: str = "HS256"
    access_token_expire_minutes: int = 30
    github_webhook_secret: str = "github-webhook-secret"
    hasura_url: str = "http://localhost:8080/v1/graphql"
    hasura_admin_secret: str = "hasura-dev-secret"

    class Config:
        env_file = ".env"


@lru_cache
def get_settings():
    return Settings()

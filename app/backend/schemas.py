from pydantic import BaseModel, EmailStr
from typing import Optional
from datetime import datetime


class UserCreate(BaseModel):
    email: EmailStr
    password: str
    full_name: Optional[str] = None


class UserResponse(BaseModel):
    id: int
    email: str
    first_name: Optional[str]
    last_name: Optional[str]
    is_active: bool
    created_at: datetime

    class Config:
        from_attributes = True


class Token(BaseModel):
    access_token: str
    token_type: str


class TokenData(BaseModel):
    email: Optional[str] = None


class MessageCreate(BaseModel):
    body: str


class MessageResponse(BaseModel):
    id: int
    user_id: int
    body: str
    is_own: bool
    created_at: datetime

    class Config:
        from_attributes = True


class GitHubWebhookPayload(BaseModel):
    action: Optional[str] = None
    event_type: str
    event_id: Optional[str] = None
    payload: dict

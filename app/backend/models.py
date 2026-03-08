from sqlalchemy import (
    Column,
    Integer,
    String,
    Boolean,
    DateTime,
    Text,
    ForeignKey,
    Enum,
)
from sqlalchemy.orm import DeclarativeBase, relationship
from sqlalchemy.sql import func
import enum


class Base(DeclarativeBase):
    pass


class User(Base):
    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True)
    email = Column(String(255), unique=True, index=True, nullable=False)
    hashed_password = Column(String(255), nullable=False)
    first_name = Column(String(255), nullable=True)
    last_name = Column(String(255), nullable=True)
    is_active = Column(Boolean, default=True)
    created_at = Column(DateTime(timezone=True), server_default=func.now())

    messages = relationship("Message", back_populates="user")


class Message(Base):
    __tablename__ = "messages"

    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    body = Column(Text, nullable=False)
    is_own = Column(Boolean, default=False)
    created_at = Column(DateTime(timezone=True), server_default=func.now())

    user = relationship("User", back_populates="messages")


class GitHubEventType(str, enum.Enum):
    PUSH = "push"
    PULL_REQUEST = "pull_request"
    ISSUES = "issues"
    ISSUE_COMMENT = "issue_comment"
    CREATE = "create"
    DELETE = "delete"


class GitHubEvent(Base):
    __tablename__ = "github_events"

    id = Column(Integer, primary_key=True, index=True)
    event_type = Column(String(50), nullable=False)
    event_id = Column(String(255), index=True, nullable=True)
    payload = Column(Text, nullable=False)
    processed = Column(Boolean, default=False)
    created_at = Column(DateTime(timezone=True), server_default=func.now())

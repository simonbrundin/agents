import json
import hmac
import hashlib
from datetime import timedelta
from typing import List

import asyncpg
import httpx
from dateutil import parser
from dateutil.relativedelta import relativedelta
from fastapi import FastAPI, Depends, HTTPException, status, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.security import OAuth2PasswordRequestForm
from pydantic import BaseModel
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select
from typing import Optional
from datetime import datetime
import random

from config import get_settings
from database import get_db, init_db, async_session
from models import User, Message, GitHubEvent
from schemas import (
    UserCreate,
    UserResponse,
    Token,
    MessageCreate,
    MessageResponse,
    GitHubWebhookPayload,
    AgentMessageCreate,
    AgentMessageResponse,
)
from auth import (
    get_password_hash,
    verify_password,
    create_access_token,
    get_current_user,
)

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.on_event("startup")
async def startup():
    await init_db()


class Avatar(BaseModel):
    src: str


class From(BaseModel):
    name: str
    email: str
    avatar: Optional[Avatar] = None


class Mail(BaseModel):
    id: int
    unread: bool = False
    from_: From
    subject: str
    body: str
    date: str


class Participant(BaseModel):
    id: int
    name: str
    email: str
    avatar: Optional[Avatar] = None
    status: str = "subscribed"
    location: str = ""


class MessageModel(BaseModel):
    id: int
    body: str
    date: str
    is_own: bool


class Conversation(BaseModel):
    id: int
    participant: Participant
    messages: list[MessageModel]
    unread_count: int
    last_message_at: str


mails: list[Mail] = [
    Mail(
        id=1,
        from_=From(
            name="Alex Smith",
            email="alex.smith@example.com",
            avatar=Avatar(src="https://i.pravatar.cc/128?u=1"),
        ),
        subject="Meeting Schedule: Q1 Marketing Strategy Review",
        body="""Dear Team,

I hope this email finds you well. Just a quick reminder about our Q1 Marketing Strategy meeting scheduled for tomorrow at 10 AM EST in Conference Room A.

Agenda:
- Q4 Performance Review
- New Campaign Proposals
- Budget Allocation for Q2
- Team Resource Planning

Please come prepared with your department updates. I've attached the preliminary deck for your review.

Best regards,
Alex Smith
Senior Marketing Director
Tel: (555) 123-4567""",
        date=datetime.now().isoformat(),
        unread=False,
    ),
    Mail(
        id=2,
        unread=True,
        from_=From(
            name="Jordan Brown",
            email="jordan.brown@example.com",
            avatar=Avatar(src="https://i.pravatar.cc/128?u=2"),
        ),
        subject="RE: Project Phoenix - Sprint 3 Update",
        body="""Hi team,

Quick update on Sprint 3 deliverables:

✅ User authentication module completed
🏗️ Payment integration at 80%
⏳ API documentation pending review

Key metrics:
- Code coverage: 94%
- Sprint velocity: 45 points
- Bug resolution rate: 98%

Please review the attached report for detailed analysis. Let's discuss any blockers in tomorrow's stand-up.

Regards,
Jordan

--
Jordan Brown
Lead Developer | Tech Solutions
Mobile: +1 (555) 234-5678""",
        date=(datetime.now() - timedelta(minutes=7)).isoformat(),
    ),
    Mail(
        id=3,
        unread=True,
        from_=From(
            name="Taylor Green",
            email="taylor.green@example.com",
            avatar=Avatar(src="https://i.pravatar.cc/128?u=3"),
        ),
        subject="Lunch Plans",
        body="""Hi there!

I was wondering if you'd like to grab lunch this Friday? There's this amazing new Mexican restaurant downtown called "La Casa" that I've been wanting to try. They're known for their authentic tacos and house-made guacamole.

Would 12:30 PM work for you? It would be great to catch up and discuss the upcoming team building event while we're there.

Let me know what you think!

Best,
Taylor""",
        date=(datetime.now() - timedelta(hours=3)).isoformat(),
    ),
]


user_replies = [
    "Thanks for letting me know! I'll take a look.",
    "Sounds good, I'll review and get back to you.",
    "Perfect, looking forward to it!",
    "Got it, I'll be there.",
    "Thanks for the update!",
    "I'll check my schedule and confirm.",
    "Great, let's discuss this further.",
    "I agree, let's proceed with that plan.",
    "I'll send over the details shortly.",
    "Thanks for the heads up!",
]


def transform_mails_to_conversations(mails: list[Mail]) -> list[Conversation]:
    result = []
    for index, mail in enumerate(mails):
        user_reply = user_replies[index % len(user_replies)]
        mail_date = parser.parse(mail.date)

        result.append(
            Conversation(
                id=mail.id,
                participant=Participant(
                    id=mail.id,
                    name=mail.from_.name,
                    email=mail.from_.email,
                    avatar=mail.from_.avatar,
                    status="subscribed",
                    location="",
                ),
                messages=[
                    MessageModel(
                        id=1, body=mail.body, date=mail_date.isoformat(), is_own=False
                    ),
                    MessageModel(
                        id=2,
                        body=user_reply,
                        date=(
                            mail_date - timedelta(minutes=random.randint(5, 30))
                        ).isoformat(),
                        is_own=True,
                    ),
                    MessageModel(
                        id=3,
                        body=f"Re: {mail.subject}",
                        date=(mail_date - timedelta(hours=1)).isoformat(),
                        is_own=False,
                    ),
                ],
                unread_count=1 if mail.unread else 0,
                last_message_at=mail.date,
            )
        )
    return result


conversations = transform_mails_to_conversations(mails)


@app.get("/api/conversations", response_model=list[Conversation])
async def get_conversations():
    return conversations


@app.get("/health")
async def health():
    return {"status": "ok"}


@app.post(
    "/api/auth/register",
    response_model=UserResponse,
    status_code=status.HTTP_201_CREATED,
)
async def register(user: UserCreate, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(User).where(User.email == user.email))
    existing_user = result.scalar_one_or_none()
    if existing_user:
        raise HTTPException(status_code=400, detail="Email already registered")

    hashed_password = get_password_hash(user.password)

    first_name = None
    last_name = None
    if user.full_name:
        name_parts = user.full_name.split(" ", 1)
        first_name = name_parts[0]
        last_name = name_parts[1] if len(name_parts) > 1 else None

    db_user = User(
        email=user.email,
        hashed_password=hashed_password,
        first_name=first_name,
        last_name=last_name,
    )
    db.add(db_user)
    await db.commit()
    await db.refresh(db_user)
    return db_user


@app.post("/api/auth/login", response_model=Token)
async def login(
    form_data: OAuth2PasswordRequestForm = Depends(), db: AsyncSession = Depends(get_db)
):
    result = await db.execute(select(User).where(User.email == form_data.username))
    user = result.scalar_one_or_none()
    if not user or not verify_password(form_data.password, str(user.hashed_password)):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect email or password",
            headers={"WWW-Authenticate": "Bearer"},
        )

    access_token = create_access_token(
        data={"sub": user.email}, expires_delta=timedelta(minutes=30)
    )
    return {"access_token": access_token, "token_type": "bearer"}


@app.get("/api/auth/me", response_model=UserResponse)
async def get_user_info(current_user: User = Depends(get_current_user)):
    return current_user


@app.post(
    "/api/messages", response_model=MessageResponse, status_code=status.HTTP_201_CREATED
)
async def create_message(
    message: MessageCreate,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user),
):
    db_message = Message(user_id=current_user.id, body=message.body, is_own=True)
    db.add(db_message)
    await db.commit()
    await db.refresh(db_message)
    return db_message


@app.get("/api/messages", response_model=List[MessageResponse])
async def get_messages(
    db: AsyncSession = Depends(get_db), current_user: User = Depends(get_current_user)
):
    result = await db.execute(
        select(Message)
        .where(Message.user_id == current_user.id)
        .order_by(Message.created_at.desc())
    )
    messages = result.scalars().all()
    return messages


@app.post(
    "/api/agent/message",
    response_model=AgentMessageResponse,
    status_code=status.HTTP_201_CREATED,
)
async def agent_message(
    agent_message: AgentMessageCreate,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user),
):
    settings = get_settings()

    mutation = """
        mutation InsertMessage($conversationId: Int!, $userId: Int!, $message: String!) {
            insert_messages_one(object: {
                conversation_id: $conversationId,
                user_id: $userId,
                message: $message
            }) {
                id
            }
        }
    """

    async with httpx.AsyncClient() as client:
        response = await client.post(
            settings.hasura_url,
            json={
                "query": mutation,
                "variables": {
                    "conversationId": agent_message.conversation_id,
                    "userId": current_user.id,
                    "message": agent_message.message,
                },
            },
            headers={
                "x-hasura-admin-secret": settings.hasura_admin_secret,
                "Content-Type": "application/json",
            },
        )

        if response.status_code != 200:
            raise HTTPException(
                status_code=500, detail="Failed to insert message in Hasura"
            )

        result = response.json()

        if "errors" in result:
            raise HTTPException(status_code=500, detail=result["errors"])

        message_id = result["data"]["insert_messages_one"]["id"]

    return {"status": "mottaget", "message_id": message_id}


def verify_github_signature(payload: bytes, signature: str) -> bool:
    settings = get_settings()
    if not signature:
        return False
    key = settings.github_webhook_secret.encode()
    expected_signature = "sha256=" + hmac.new(key, payload, hashlib.sha256).hexdigest()
    return hmac.compare_digest(expected_signature, signature)


@app.post("/api/webhooks/github", status_code=status.HTTP_200_OK)
async def github_webhook(request: Request, db: AsyncSession = Depends(get_db)):
    signature = request.headers.get("X-Hub-Signature-256")
    event_type = request.headers.get("X-GitHub-Event", "unknown")
    event_id = request.headers.get("X-GitHub-Delivery")

    body = await request.body()

    if not verify_github_signature(body, signature or ""):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid signature"
        )

    try:
        payload = json.loads(body)
    except json.JSONDecodeError:
        raise HTTPException(status_code=400, detail="Invalid JSON")

    github_event = GitHubEvent(
        event_type=event_type,
        event_id=event_id,
        payload=json.dumps(payload),
        processed=False,
    )
    db.add(github_event)
    await db.commit()

    return {"status": "received", "event_type": event_type, "event_id": event_id}

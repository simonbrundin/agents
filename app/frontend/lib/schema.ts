import { pgTable, serial, varchar, text, timestamp, integer } from 'drizzle-orm/pg-core'

export const user = pgTable('users', {
  id: serial('id').primaryKey(),
  email: varchar('email', { length: 40 }).notNull().unique(),
  firstName: varchar('first_name', { length: 16 }),
  lastName: varchar('last_name', { length: 40 }),
  created: timestamp('created', { withTimezone: true, mode: 'timestamp' }).defaultNow().notNull()
})

export const conversation = pgTable('conversations', {
  id: serial('id').primaryKey(),
  userId: integer('user_id').references(() => user.id).notNull(),
  repoUrl: text('repo_url')
})

export const message = pgTable('messages', {
  id: serial('id').primaryKey(),
  conversationId: integer('conversation_id').references(() => conversation.id).notNull(),
  message: text('message').notNull(),
  time: timestamp('time', { withTimezone: true, mode: 'timestamp' }).defaultNow().notNull()
})

export type User = typeof user.$inferSelect
export type Conversation = typeof conversation.$inferSelect
export type Message = typeof message.$inferSelect

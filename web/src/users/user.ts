import z from "zod";

export const userSchema = z.object({
  id: z.string(),
});

export type User = z.infer<typeof userSchema>;

export const messageSchema = z.object({
  id: z.string(),
  role: z.enum(["user", "system"]),
  user_id: z.string(),
  message: z.string(),
  sent_at: z.string(),
}).transform((value) => ({
  id: value.id,
  role: value.role as "user" | "system",
  userId: value.user_id,
  message: value.message,
  sentAt: value.sent_at,
}));

export type Message = z.infer<typeof messageSchema>;

export const userWithMessagesSchema = z.object({
	user: userSchema,
	messages: z.array(messageSchema),
})

export type UserWithMessages = z.infer<typeof userWithMessagesSchema>;

import z from "zod";

export const userSchema = z.object({
  id: z.string(),
});

export type User = z.infer<typeof userSchema>;

export function userFromJson(input: string): User {
  const json = JSON.parse(input);

  return userSchema.parse({
    id: json.id,
    role: json.role,
    userId: json.user_id,
    message: json.message,
    sentAt: json.sent_at,
  });
}

export const messageSchema = z.object({
  id: z.string(),
  role: z.enum(["user", "system"]),
  userId: z.string(),
  message: z.string(),
  sentAt: z.string(),
});

export type Message = z.infer<typeof messageSchema>;

export function messageFromJson(input: string): Message {
  const json = JSON.parse(input);

  return messageSchema.parse({
    id: json.id,
    role: json.role,
    userId: json.user_id,
    message: json.message,
    sentAt: json.sent_at,
  });
}

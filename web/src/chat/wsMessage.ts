import z from "zod";
import type { Message } from "../users/user";

export const wsMessageSchema = z
  .object({
    user_id: z.string(),
    text: z.string(),
  })
  .transform((data) => ({
    userId: data.user_id,
    text: data.text,
  }));

// TODO: remove this and change the ws messages to the proper format
export const wsMessageToUserMessage = (
  message: WsMessage,
  role: "user" | "system",
): Message => {
  return {
    id: message.userId + "-" + Math.random(),
    role,
    userId: message.userId,
    message: message.text,
    sentAt: new Date().toISOString(),
  };
};

export type WsMessage = z.infer<typeof wsMessageSchema>;

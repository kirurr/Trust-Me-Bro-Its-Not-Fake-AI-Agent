import z from "zod";
import { userSchema, type User } from "./user";

const backendUrl = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";

export async function getAllUsers(): Promise<User[]> {
  const response = await fetch(new URL("/users", backendUrl));

  if (!response.ok) {
    throw new Error("Failed to fetch users");
  }

  const json = await response.json();
  return z.array(userSchema).parse(json);
}

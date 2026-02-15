import z from "zod";
import { loginSchema, signupSchema } from "./schema";

export type SignupInput = z.infer<typeof signupSchema>;
export type LoginInput = z.infer<typeof loginSchema>;
export type User = { name: string; email: string };

export type LoginResponse = {
  session_id: string;
  refresh_token: string;
  access_token_expires_at: string;
  refresh_token_expires_at: string;
  access_token: string;
  User: User;
};

export type AuthUserDetails = {
  User: User;
  access_token: string;
  access_token_expires_at: string;
  session_id: string;
};

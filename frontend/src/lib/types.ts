import z, { string } from "zod";
import { loginSchema, signupSchema } from "./schema";

export type SignupInput = z.infer<typeof signupSchema>;
export type LoginInput = z.infer<typeof loginSchema>;

export type User = {
  name: string;
  email: string;
};

export type AccessTokenDetails = {
  access_token: string;
  access_token_expires_at: string;
};

export type RefreshTokenDetails = {
  refresh_token: string;
  refresh_token_expires_at: string;
};

export type SessionDetails = {
  session_id: string;
};

export type TokenDetails = AccessTokenDetails & SessionDetails;

export type AuthUserDetails = {
  User: User;
} & AccessTokenDetails &
  SessionDetails;

export type LoginResponse = {
  User: User;
} & AccessTokenDetails &
  RefreshTokenDetails &
  SessionDetails;

export type Project = {
  Base: {
    id: string;
    created_at: string;
    updated_at: string;
    deleted_at: string;
  };
  name: string;
  git_url: string;
  sub_domain: string;
  custom_domain: string;
  user_id: string;
};

import z, { string } from "zod";
import { loginSchema, signupSchema } from "./schema";

export type SignupInput = z.infer<typeof signupSchema>;
export type LoginInput = z.infer<typeof loginSchema>;

export type User = {
  id: string;
  name: string;
  email: string;
  is_verified: boolean;
} | null;

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

export type TokenDetails = (AccessTokenDetails & SessionDetails) | null;

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
  ID: string;
  CreatedAt: string;
  Name: string;
  GitUrl: string;
  SubDomain: string;
};

// export type CreateProjectResponse = {
//   name: string;
//   id: string;
//   sub_domain: string;
// };

export type LogEvent = {
  log: string;
  created_at: string;
  metadata: any;
};

export type Deployment = {
  created_at: string;
  id: string;
  status: string;
  logs: LogEvent[];
};

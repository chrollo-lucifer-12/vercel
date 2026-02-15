import "server-only";

import z from "zod";

const envSchema = z.object({
  BACKEND_ENDPOINT: z.string().min(1),
  LOGIN_ENDPOINT: z.string().min(1),
  SIGNUP_ENDPOINT: z.string().min(1),
  LOGOUT_ENDPOINT: z.string().min(1),
});

export const serverEnv = envSchema.parse(process.env);

declare global {
  namespace NodeJS {
    interface ProcessEnv extends z.infer<typeof envSchema> {}
  }
}

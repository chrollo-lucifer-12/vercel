import z from "zod";
import { loginSchema } from "./schema";

export type LoginInput = z.infer<typeof loginSchema>;

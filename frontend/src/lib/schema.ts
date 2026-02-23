import z from "zod";

export const loginSchema = z.object({
  email: z.email({ error: "Enter a valid email." }),
  password: z
    .string()
    .min(8, { error: "Password must be atleast 8 characters long." }),
});

export const signupSchema = z
  .object({
    name: z.string().min(1, { error: "Name cannot be empty." }),
  })
  .extend(loginSchema.shape);

export const createProjectSchema = z.object({
  name: z
    .string()
    .min(3, { error: "Project name should be atleast 3 characters." }),

  git_url: z.url(),
});

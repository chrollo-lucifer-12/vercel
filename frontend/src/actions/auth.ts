"use server";

import { logout, signin, signup } from "@/lib/axios/auth";
import { loginSchema, signupSchema } from "@/lib/schema";
import { cookies } from "next/headers";
import z from "zod";

export const logoutAction = async (sessionId: string) => {
  try {
    await logout(sessionId);
    const cookieStore = await cookies();
    cookieStore.delete("refresh_token");
    return { success: true };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Failed to logout",
    };
  }
};

export const signupAction = async (formData: FormData) => {
  const parsedRes = signupSchema.safeParse({
    email: formData.get("email"),
    password: formData.get("password"),
    name: formData.get("name"),
  });
  if (!parsedRes.success) {
    return {
      success: false,
      error: z.treeifyError(parsedRes.error),
    };
  }

  try {
    const { email, name, password } = parsedRes.data;
    await signup({ email, name, password });
    return { success: true };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Failed to sign up",
    };
  }
};

export const signinAction = async (formData: FormData) => {
  const parsedRes = loginSchema.safeParse({
    email: formData.get("email"),
    password: formData.get("password"),
  });

  if (!parsedRes.success) {
    return {
      success: false,
      error: z.treeifyError(parsedRes.error),
    };
  }

  try {
    const { email, password } = parsedRes.data;
    const {
      User,
      access_token,
      access_token_expires_at,
      refresh_token,
      refresh_token_expires_at,
      session_id,
    } = await signin({ email, password });

    const refreshTokenStr = JSON.stringify({
      refresh_token,
      refresh_token_expires_at,
    });

    const cookieStore = await cookies();
    cookieStore.set("refresh_token", refreshTokenStr, {
      secure: true,
      httpOnly: true,
      expires: new Date(refresh_token_expires_at),
    });

    return {
      success: true,
      data: { User, access_token, access_token_expires_at, session_id },
    };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Failed to sign in",
    };
  }
};

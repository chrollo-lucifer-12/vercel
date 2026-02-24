"use server";

import { logout, refresh, signin, signup, verify } from "@/lib/axios/auth";
import { loginSchema, signupSchema } from "@/lib/schema";
import { RefreshTokenDetails, TokenDetails } from "@/lib/types";
import { cookies } from "next/headers";
import z from "zod";

export const verifyAction = async (formData: FormData) => {
  const parsedRes = z
    .object({ email: z.email() })
    .safeParse({ email: formData.get("email") });
  if (!parsedRes.success) {
    return {
      success: false,
      error: parsedRes.error.flatten().fieldErrors,
    };
  }
  try {
    await verify(parsedRes.data.email);
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Failed to send mail",
    };
  }
};

export const logoutAction = async (sessionId: string, accessToken: string) => {
  try {
    await logout(sessionId, accessToken);
    const cookieStore = await cookies();
    cookieStore.delete("refresh_token");
    return { success: true };
  } catch (error) {
    console.log(error);
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
      error: parsedRes.error.flatten().fieldErrors,
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

    console.log(refreshTokenStr);
    console.log(new Date(refresh_token_expires_at));

    const cookieStore = await cookies();
    cookieStore.set({
      name: "refresh_token",
      value: refreshTokenStr,
      path: "/",
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

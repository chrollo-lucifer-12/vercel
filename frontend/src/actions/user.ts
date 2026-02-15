"use server";

import { profile } from "@/lib/axios/user";

export const profileAction = async (accessToken: string) => {
  try {
    const res = await profile(accessToken);
    return {
      success: true,
      user: res,
    };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Failed to fetch profile",
    };
  }
};

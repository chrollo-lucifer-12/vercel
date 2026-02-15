import "server-only";

import { axiosInstance } from "./axios";
import { LoginInput, LoginResponse, SignupInput } from "./types";
import { serverEnv } from "./env/server";

export const signup = async (data: SignupInput) => {
  try {
    await axiosInstance.post(serverEnv.SIGNUP_ENDPOINT, data);
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to sign up");
  }
};

export const signin = async (data: LoginInput) => {
  try {
    const res = await axiosInstance.post(serverEnv.LOGIN_ENDPOINT, data);
    const {
      User,
      access_token,
      access_token_expires_at,
      refresh_token,
      refresh_token_expires_at,
      session_id,
    } = res.data as LoginResponse;

    return {
      User,
      access_token,
      access_token_expires_at,
      refresh_token,
      refresh_token_expires_at,
      session_id,
    };
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to sign in");
  }
};

export const logout = async (sessionId: string) => {
  try {
    await axiosInstance.delete(`${serverEnv.LOGOUT_ENDPOINT}/${sessionId}`);
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to sign in");
  }
};

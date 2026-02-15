import "server-only";
import { axiosInstance } from "./axios";
import {
  AccessTokenDetails,
  LoginInput,
  LoginResponse,
  RefreshTokenDetails,
  SignupInput,
  TokenDetails,
} from "../types";
import { serverEnv } from "../env/server";

export const refresh = async (refreshToken: string) => {
  try {
    const res = await axiosInstance.post(serverEnv.REFRESH_ENDPOINT, {
      refresh_token: refreshToken,
    });
    return res.data as TokenDetails;
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to refresh");
  }
};

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

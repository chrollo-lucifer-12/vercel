import { axiosInstance } from "./axios";
import {
  AccessTokenDetails,
  LoginInput,
  LoginResponse,
  RefreshTokenDetails,
  SignupInput,
  TokenDetails,
} from "../types";
import { clientEnv } from "../env/client";

export const refresh = async (refreshToken: string) => {
  try {
    const res = await axiosInstance.post(
      clientEnv.NEXT_PUBLIC_REFRESH_ENDPOINT,
      {
        refresh_token: refreshToken,
      },
    );
    return res.data as TokenDetails;
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to refresh");
  }
};

export const signup = async (data: SignupInput) => {
  try {
    await axiosInstance.post(clientEnv.NEXT_PUBLIC_SIGNUP_ENDPOINT, data);
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to sign up");
  }
};

export const signin = async (data: LoginInput) => {
  try {
    const res = await axiosInstance.post(
      clientEnv.NEXT_PUBLIC_LOGIN_ENDPOINT,
      data,
    );
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

export const logout = async (sessionId: string, accessToken: string) => {
  try {
    await axiosInstance.delete(
      `${clientEnv.NEXT_PUBLIC_LOGOUT_ENDPOINT}/${sessionId}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to sign in");
  }
};

export const verify = async (email: string) => {
  try {
    await axiosInstance.post(`${clientEnv.NEXT_PUBLIC_VERIFY_ENDPOINT}`, {
      email,
    });
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to verify");
  }
};

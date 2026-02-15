import "server-only";
import { axiosInstance } from "./axios";
import { serverEnv } from "../env/server";
import { User } from "../types";

export const profile = async (accessToken: string) => {
  try {
    const res = await axiosInstance.get(serverEnv.PROFILE_ENDPOINT, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
    return res.data as User;
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to fetch profile");
  }
};

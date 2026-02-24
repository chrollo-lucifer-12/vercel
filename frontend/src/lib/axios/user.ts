import { axiosInstance } from "./axios";
import { clientEnv } from "../env/client";
import { User } from "../types";

export const profile = async (accessToken: string) => {
  try {
    const res = await axiosInstance.get(
      clientEnv.NEXT_PUBLIC_PROFILE_ENDPOINT,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    return res.data as User;
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to fetch profile");
  }
};

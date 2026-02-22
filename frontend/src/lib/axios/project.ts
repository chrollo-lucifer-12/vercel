import { serverEnv } from "../env/server";
import { Project } from "../types";
import { axiosInstance } from "./axios";

export const getProjects = async (
  accessToken: string,
  limit: number,
  offset: number,
  name: string,
) => {
  try {
    const res = await axiosInstance.get(serverEnv.GET_PROJECT_ENDPOINT, {
      params: {
        limit,
        offset,
        name,
      },
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
    return res.data as Project[];
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to get projects");
  }
};

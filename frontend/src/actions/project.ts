"use server";

import { getProjects } from "@/lib/axios/project";

export const projectAction = async (
  accessToken: string,
  limit: number,
  offset: number,
  name: string,
) => {
  try {
    const data = await getProjects(accessToken, limit, offset, name);
    return {
      success: true,
      projects: data,
    };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Failed to fetch profile",
    };
  }
};

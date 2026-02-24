import axios from "axios";
import { clientEnv } from "../env/client";
import { CreateProjectResponse, Project } from "../types";
import { axiosInstance } from "./axios";

export const createProject = async (
  accessToken: string,
  name: string,
  gitUrl: string,
) => {
  try {
    const res = await axiosInstance.post(
      clientEnv.NEXT_PUBLIC_CREATE_PROJECT_ENDPOINT,
      {
        project_name: name,
        github_url: gitUrl,
      },
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    return res.data as CreateProjectResponse;
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to create project");
  }
};

export const getProjects = async (
  accessToken: string,
  limit: number,
  offset: number,
  name: string,
) => {
  try {
    const res = await axiosInstance.get(
      clientEnv.NEXT_PUBLIC_ALL_PROJECT_ENDPOINT,
      {
        params: {
          limit,
          offset,
          name,
        },
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );

    return res.data as Project[];
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to get projects");
  }
};

export const getProjectsTest = async (
  limit: number,
  offset: number,
  name?: string,
) => {
  try {
    const res = await axios.get("http://localhost:4000", {
      params: {
        _start: offset,
        _limit: limit,
        ...(name ? { q: name } : {}), // search
      },
    });

    return res.data as Project[];
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to get projects");
  }
};

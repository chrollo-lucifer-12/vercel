import axios from "axios";
import { clientEnv } from "../env/client";
import { Project, ProjectWithDeployment, WebsiteAnalytics } from "../types";
import { axiosInstance } from "./axios";

export const createProject = async (
  accessToken: string,
  name: string,
  gitUrl: string,
) => {
  try {
    const res = await axiosInstance.post<Project>(
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
    return res.data;
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
    const res = await axiosInstance.get<Project[]>(
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

    return res.data;
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to get projects");
  }
};

export const deleteProject = async (accessToken: string, projectId: string) => {
  try {
    await axiosInstance.delete(
      `${clientEnv.NEXT_PUBLIC_DELETE_PROJECT}/${projectId}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to delete project");
  }
};

export const getProject = async (accessToken: string, slug: string) => {
  try {
    const res = await axiosInstance.get<ProjectWithDeployment>(
      `${clientEnv.NEXT_PUBLIC_GET_PROJECT}/${slug}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    return res.data;
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to get project");
  }
};

export const getProjectAnalytics = async (
  accessToken: string,
  slug: string,
  from: Date | null,
  to: Date | null,
) => {
  try {
    const res = await axiosInstance.get<WebsiteAnalytics[]>(
      `${clientEnv.NEXT_PUBLIC_GET_ANALYTICS}/${slug}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    return res.data;
  } catch (err) {
    console.error(err);
    throw err instanceof Error
      ? err
      : new Error("Failed to get project analytics");
  }
};

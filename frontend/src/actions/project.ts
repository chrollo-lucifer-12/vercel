"use server";

import { createProject, getProjects } from "@/lib/axios/project";
import { createProjectSchema } from "@/lib/schema";

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

export const createProjectAction = async (
  accessToken: string,
  formData: FormData,
) => {
  const parsedSchema = createProjectSchema.safeParse({
    name: formData.get("name"),
    git_url: formData.get("git_url"),
  });

  if (!parsedSchema.success) {
    return {
      success: false,
      error: parsedSchema.error.flatten().fieldErrors,
    };
  }

  try {
    const data = await createProject(
      accessToken,
      parsedSchema.data.name,
      parsedSchema.data.git_url,
    );
    return { success: true, subDomain: data.sub_domain };
  } catch (error) {
    return {
      success: false,
      error:
        error instanceof Error ? error.message : "Failed to create project",
    };
  }
};

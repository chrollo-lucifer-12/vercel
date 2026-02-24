"use server";

import { createProject } from "@/lib/axios/project";
import { createProjectSchema } from "@/lib/schema";

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

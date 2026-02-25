"use client";

import { useProject } from "@/hooks/use-project";

const ProjectPage = ({ subdomain }: { subdomain: string }) => {
  const { data } = useProject(subdomain);

  console.log(data);

  return null;
};

export default ProjectPage;
